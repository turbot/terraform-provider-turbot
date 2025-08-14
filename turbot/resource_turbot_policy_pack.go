package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

// properties which must be passed to a create/update call
var policyPackProperties = []interface{}{"title", "description", "parent", "filter", "tags", "akas"}

func getPolicyPackUpdateProperties() []interface{} {
	excludedProperties := []string{"parent"}
	return helpers.RemoveProperties(policyPackProperties, excludedProperties)
}

func resourceTurbotPolicyPack() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotPolicyPackCreate,
		Read:   resourceTurbotPolicyPackRead,
		Update: resourceTurbotPolicyPackUpdate,
		Delete: resourceTurbotPolicyPackDelete,
		Exists: resourceTurbotPolicyPackExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotPolicyPackImport,
		},
		Schema: map[string]*schema.Schema{
			//aka of the parent resource
			"parent": {
				Type: schema.TypeString,
				// when doing a diff, the state file will contain the id of the parent but the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
				Optional:         true,
				Default:          "tmod:@turbot/turbot#/",
			},
			//when doing a read, fetch the parent akas to use in suppressIfAkaMatches
			"parent_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"akas": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: suppressIfAkaRemoved(),
			},
		},
	}
}

func resourceTurbotPolicyPackExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotPolicyPackCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build map of folder properties
	input := mapFromResourceData(d, policyPackProperties)

	policyPack, err := client.CreateSmartFolder(input)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(policyPack.Turbot.Id)
	// TODO Remove Read call once schema changes are In.
	return resourceTurbotPolicyPackRead(d, meta)
}

func resourceTurbotPolicyPackUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build map of folder properties
	input := mapFromResourceData(d, getPolicyPackUpdateProperties())
	input["id"] = id

	_, err := client.UpdateSmartFolder(input)
	if err != nil {
		return err
	}
	// set 'Read' Properties
	// TODO Remove Read call once schema changes are In.
	return resourceTurbotPolicyPackRead(d, meta)
}

func resourceTurbotPolicyPackRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	policyPack, err := client.ReadSmartFolder(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(policyPack.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// NOTE currently turbot accepts array of filters but only uses the first
	if len(policyPack.Filters) > 0 {
		d.Set("filter", policyPack.Filters[0])
	}
	d.Set("parent", policyPack.Parent)
	d.Set("title", policyPack.Title)
	d.Set("description", policyPack.Description)
	d.Set("tags", policyPack.Turbot.Tags)
	d.Set("akas", policyPack.Turbot.Akas)

	return nil
}

func resourceTurbotPolicyPackDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	err := client.DeleteSmartFolder(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")

	return nil
}

func resourceTurbotPolicyPackImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotPolicyPackRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// Suppress the diff if trying to remove the akas completely
func suppressIfAkaRemoved() func(k, old, new string, d *schema.ResourceData) bool {
	return func(k, old, new string, d *schema.ResourceData) bool {
		oldVal, newVal := d.GetChange("akas")
		return slicesAreEquivalent(oldVal.([]interface{}), newVal.([]interface{}))
	}
}

func slicesAreEquivalent(old, new []interface{}) bool {
	if len(old) != len(new) {
		return false
	}
	freqMap := make(map[interface{}]int)
	for _, item := range old {
		freqMap[item]++
	}
	for _, item := range new {
		if count, exists := freqMap[item]; !exists || count == 0 {
			return false
		}
		freqMap[item]--
	}
	return true
}
