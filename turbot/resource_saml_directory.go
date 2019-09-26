package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// these are the properties which must be passed to a create/update call
var samlDirectoryProperties = []string{"description", "directory_type", "status", "entry_point", "issuer", "certificate", "profile_id_template", "group_id_template", "name_id_format", "sign_requests", "signature_private_key", "signature_algorithm", "pool_id"}
var samlDirectoryMetadataProperties = []string{"tags"}

func getSamlDirectoryUpdateProperties() []string {
	excludedProperties := []string{"directory_type"}
	return removeProperties(samlDirectoryProperties, excludedProperties)
}

func resourceTurbotSamlDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotSamlDirectoryCreate,
		Read:   resourceTurbotSamlDirectoryRead,
		Update: resourceTurbotSamlDirectoryUpdate,
		Delete: resourceTurbotSamlDirectoryDelete,
		Exists: resourceTurbotSamlDirectoryExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotSamlDirectoryImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resourcesamlDirectoryProperties
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent bu the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfParentAkaMatches,
			},
			// when doing a read, fetch the parent akas to use in suppressIfParentAkaMatches()
			"parent_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"description": {
				Type:     schema.TypeString,
				Required: true,
			},
			"directory_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entry_point": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile_id_template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_id_format": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"sign_requests": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signature_private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTurbotSamlDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotSamlDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	// build mutation payload
	payload := map[string]map[string]interface{}{
		"data":       mapFromResourceData(d, samlDirectoryProperties),
		"turbotData": mapFromResourceData(d, samlDirectoryMetadataProperties),
	}
	// set computed properties
	payload["data"]["status"] = "New"
	payload["data"]["directoryType"] = "saml"

	turbotMetadata, err := client.CreateSamlDirectory(parentAka, payload)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(turbotMetadata.ParentId)
	if err != nil {
		return err
	}
	d.Set("parent_akas", parentAkas)
	// assign the id
	d.SetId(turbotMetadata.Id)
	d.Set("status", payload["data"]["status"])
	d.Set("directoryType", payload["data"]["directoryType"])
	return nil
}

func resourceTurbotSamlDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	payload := map[string]map[string]interface{}{
		"data":       mapFromResourceData(d, getSamlDirectoryUpdateProperties()),
		"turbotData": mapFromResourceData(d, samlDirectoryMetadataProperties),
	}

	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateSamlDirectory(id, parentAka, payload)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(turbotMetadata.ParentId)
	if err != nil {
		return err
	}
	d.Set("parent_akas", parentAkas)
	return nil
}

func resourceTurbotSamlDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	samlDirectory, err := client.ReadSamlDirectory(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// saml directory was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(samlDirectory.Turbot.ParentId)
	if err != nil {
		return err
	}
	d.Set("parent_akas", parentAkas)
	d.Set("parent", samlDirectory.Parent)
	d.Set("title", samlDirectory.Title)
	return nil
}

func resourceTurbotSamlDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotSamlDirectoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSamlDirectoryRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
