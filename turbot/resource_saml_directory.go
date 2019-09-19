package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// these are the properties which must be passed to a create/update call
var samlDirectoryProperties = []string{"description", "directoryType", "status", "entryPoint", "issuer", "certificate", "profileIdTemplate", "groupIdTemplate", "nameIdFormat", "signRequests", "signaturePrivateKey", "signatureAlgorithm", "poolId"}

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
			// aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent bu the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: supressIfParentAkaMatches,
			},
			// when doing a read, fetch the parent akas to use in supresIfParentAkaMatches()
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
			"directoryType": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"entryPoint": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profileIdTemplate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"groupIdTemplate": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"nameIdFormat": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signRequests": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signaturePrivateKey": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"signatureAlgorithm": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"poolId": {
				Type:     schema.TypeString,
				Computed: true,
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
	// build map of Saml directory properties
	data := mapFromResourceData(d, samlDirectoryProperties)
	data["status"] = "New"
	data["directoryType"] = "saml"
	turbotMetadata, err := client.CreateSamlDirectory(parentAka, data)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)
	d.Set("status", data["status"])
	d.Set("directoryType", data["directoryType"])
	return nil
}

func resourceTurbotSamlDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build map of Saml directory properties
	data := mapFromResourceData(d, folderProperties)
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateDirectory(id, parentAka, data)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}
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
	if err = setParentAkas(samlDirectory.Turbot.ParentId, d, meta); err != nil {
		return err
	}
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
