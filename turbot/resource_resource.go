package turbot

import (
	"encoding/json"
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
)

func resourceTurbotResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotResourceCreate,
		Read:   resourceTurbotResourceRead,
		Update: resourceTurbotResourceUpdate,
		Delete: resourceTurbotResourceDelete,
		Exists: resourceTurbotResourceExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotResourceImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent bu tthe config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: supressIfParentAkaMatches,
			},
			// when doing a read, fetch the parent akas to use in supressIfParentAkaMatches()
			"parent_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"payload": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressIfPayloadMatches,
			},
		},
	}
}

func resourceTurbotResourceExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parent := d.Get("parent").(string)
	resourceType := d.Get("type").(string)
	payload := d.Get("payload").(string)

	// create resource returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateResource(resourceType, parent, payload)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)
	// save formatted version of the payload for consistency
	d.Set("payload", formatPayload(payload))

	return nil
}

func resourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	// build required properties from payload
	properties, err := propertiesFromPayload(d.Get("payload").(string))
	if err != nil {
		return fmt.Errorf("error retrieving properties from resource payload: %s", err.Error())
	}

	resource, err := client.ReadResource(id, properties)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}

	// rebuild payload from the resource
	payload, err := payloadFromResource(resource.Data)
	if err != nil {
		return fmt.Errorf("error building resource payload: %s", err.Error())
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(resource.Turbot.ParentId, d, meta); err != nil {
		return err
	}
	d.Set("parent", resource.Turbot.ParentId)
	d.Set("payload", payload)

	return nil
}

func resourceTurbotResourceUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	payload := d.Get("payload").(string)
	parent := d.Get("parent").(string)
	resourceType := d.Get("type").(string)
	id := d.Id()
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateResource(id, parent, resourceType, payload)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}
	return nil
}

func resourceTurbotResourceDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()
	err := client.DeleteResource(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")

	return nil
}

func resourceTurbotResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotResourceRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// given a json string, unmarshal into a map and return a map of alias ->  propertyName
func propertiesFromPayload(payload string) (map[string]string, error) {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		return nil, err
	}
	var properties = map[string]string{}
	for k := range data {
		properties[k] = k
	}
	return properties, nil
}

// given a map of resource properties, marshal into a json string
func payloadFromResource(d map[string]interface{}) (string, error) {
	payload, err := json.MarshalIndent(d, "", " ")
	if err != nil {
		return "", err
	}
	return string(payload), nil
}

func setParentAkas(parentId string, d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)

	// load parent resource to get parent_akas
	parent, err := client.ReadResource(parentId, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to load parentAka resource; %s", err)
		return err
	}
	parentAkas := parent.Turbot.Akas
	// if this resource has no akas, just use the id
	if parentAkas == nil {
		parentAkas = []string{parentId}
	}

	// assign parent_akas
	d.Set("parent_akas", parent.Turbot.Akas)
	return nil
}

// the 'parent' in the config is an aka - however the state file will have an id.
// to perform a diff we also store parent_akas in state file, which is the list of akas for the parent
// if the new value of parent existts in parent_akas, then suppress diff
func supressIfParentAkaMatches(k, old, new string, d *schema.ResourceData) bool {
	parentAkasProperty, parentAkasSet := d.GetOk("parent_akas")
	// if parent_id has not been set yet, do not suppress the diff
	if !parentAkasSet {
		return false
	}

	parentAkas, ok := parentAkasProperty.([]interface{})
	if !ok {
		return false
	}
	// if parentAkas contains 'new', suppress diff
	for _, aka := range parentAkas {
		if aka.(string) == new {
			return true
		}
	}
	return false
}

// payload is a json string
// apply standard formatting to old and new payloads then compare
func suppressIfPayloadMatches(k, old, new string, d *schema.ResourceData) bool {
	if old == "" || new == "" {
		return false
	}
	return formatPayload(old) == formatPayload(new)
}

// apply standard formatting to the json payload to enable easy diffing
func formatPayload(payload string) string {
	data := map[string]interface{}{}
	if err := json.Unmarshal([]byte(payload), &data); err != nil {
		// ignore error and just return original payload
		return payload
	}
	payload, err := payloadFromResource(data)
	if err != nil {
		// ignore error and just return original payload
		return payload
	}
	return payload

}

func mapFromResourceData(d *schema.ResourceData, properties []string)map[string]interface{}{
	var propertyMap = map[string]interface{}{}
	for _, p := range properties {
		// get schema for property
		value := d.Get(p)
		if value != nil {
			propertyMap[p] = value;
		}
	}
	return propertyMap
}