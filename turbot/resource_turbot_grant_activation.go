package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
)

var grantActivationInputProperties = []interface{}{"grant", "resource"}

func resourceTurbotGrantActivation() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotGrantActivateCreate,
		Read:   resourceTurbotGrantActivateRead,
		Delete: resourceTurbotGrantActivateDelete,
		Exists: resourceTurbotGrantActivateExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotGrantActivateImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the resource resource
			"resource": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the resource but the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("resource_akas"),
				ForceNew:         true,
			},
			// when doing a read, fetch the resource akas to use in suppressIfAkaMatches
			"resource_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				ForceNew: true,
			},
			// the grant id (grants do not (currently) have akas
			"grant": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTurbotGrantActivateExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.GrantActivationExists(id)
}

func resourceTurbotGrantActivateCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resourceAka := d.Get("resource").(string)
	input := mapFromResourceData(d, grantActivationInputProperties)
	TurbotGrantMetadata, err := client.CreateGrantActivation(input)
	if err != nil {
		return err
	}

	// set resource_akas property by loading resource and fetching the akas
	if err := storeAkas(resourceAka, "resource_akas", d, meta); err != nil {
		return err
	}
	// assign results back into ResourceData
	d.Set("grant", TurbotGrantMetadata.GrantId)
	d.Set("resource", TurbotGrantMetadata.ResourceId)
	// assign the id
	d.SetId(TurbotGrantMetadata.Id)
	return nil
}

func resourceTurbotGrantActivateRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	activeGrant, err := client.ReadGrantActivation(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// Grant was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("grant", activeGrant.Turbot.GrantId)
	d.Set("resource", activeGrant.Turbot.ResourceId)
	// set resource_akas property by loading resource and fetching the akas
	return storeAkas(activeGrant.Turbot.ResourceId, "resource_akas", d, meta)
}

func resourceTurbotGrantActivateDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	err := client.DeleteGrantActivation(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")
	return nil
}

func resourceTurbotGrantActivateImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotGrantActivateRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
