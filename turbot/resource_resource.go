package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
)

func resourceTurbotResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotResourceCreate,
		Read:   resourceTurbotResourceRead,
		Update: resourceTurbotResourceCreate,
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
			"resourceType": {
				Type:     schema.TypeString,
				Required: true,
			}, "payload": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTurbotResourceExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	log.Println("resourceTurbotResourceExists")
	client := meta.(*apiclient.Client)
	id := d.Id()

	_, err := client.ReadResource(id, nil)
	if err != nil {
		if apiclient.NotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceTurbotResourceCreate(d *schema.ResourceData, meta interface{}) error {
	log.Println("resourceTurbotResourceCreate")
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
	if err = setParentAkas(d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)

	return nil
}

func resourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("resourceTurbotResourceRead")
	client := meta.(*apiclient.Client)
	id := d.Id()

	// todo build required properties from payload

	resource, err := client.ReadResource(id, nil)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}

	// todo rebuild payload from properties
	payload := ""

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(d, meta); err != nil {
		return err
	}

	d.Set("parent", resource.Turbot.ParentId)
	d.Set("payload", payload)

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

func setParentAkas(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)

	// load parent resource to get parent_akas
	parent, err := client.ReadResource(parentAka, nil)
	if err != nil {
		log.Printf("[ERROR] Failed to load parentAka resource; %s", err)
		return err
	}
	// assign parent_akas
	d.Set("parent_akas", parent.Turbot.Akas)
	return nil
}

// the 'parent' in the config is an aka - however the state file will have an id.
// to perform a diff we also store parent_akas in state file, which is the list of akas for the parent
// if the new value of parent existists in parent_akas, then suppress diff
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
