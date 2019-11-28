package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
)

// these are the properties which must be passed to a create/update call
var samlDirectoryDataProperties = []interface{}{"description", "status", "entry_point", "issuer", "certificate", "profile_id_template", "group_id_template", "name_id_format", "sign_requests", "signature_private_key", "signature_algorithm", "pool_id"}
var samlDirectoryInputProperties = []interface{}{"parent", "tags"}

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
				// when doing a diff, the state file will contain the id of the parent but he config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
			},
			// when doing a read, fetch the parent akas to use in suppressIfAkaMatches("parent_akas")()
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
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
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
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotSamlDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build mutation payload
	input := mapFromResourceData(d, samlDirectoryInputProperties)
	data := mapFromResourceData(d, samlDirectoryDataProperties)
	// set computed properties
	data["status"] = "Active"
	data["directoryType"] = "saml"
	input["data"] = data

	samlDirectoryMetadata, err := client.CreateSamlDirectory(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(samlDirectoryMetadata.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(samlDirectoryMetadata.Turbot.Id)
	// assign Read query properties
	d.Set("parent", samlDirectoryMetadata.Parent)
	d.Set("title", samlDirectoryMetadata.Title)
	// assign computed properties
	d.Set("status", data["status"])
	d.Set("directoryType", data["directoryType"])
	return nil
}

func resourceTurbotSamlDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	samlDirectory, err := client.ReadSamlDirectory(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// saml directory was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(samlDirectory.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	d.Set("parent", samlDirectory.Parent)
	d.Set("title", samlDirectory.Title)
	return nil
}

func resourceTurbotSamlDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build mutation payload
	input := mapFromResourceData(d, samlDirectoryInputProperties)
	input["data"] = mapFromResourceData(d, samlDirectoryDataProperties)
	input["id"] = d.Id()

	// create folder returns turbot resource metadata containing the id
	samlDirectoryMetadata, err := client.UpdateSamlDirectory(input)
	if err != nil {
		return err
	}
	// assign Read query properties
	d.Set("parent", samlDirectoryMetadata.Parent)
	d.Set("title", samlDirectoryMetadata.Title)
	// set parent_akas property by loading parent resource and fetching the akas
	return storeAkas(samlDirectoryMetadata.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotSamlDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
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
