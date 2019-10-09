package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
)

// properties which must be passed to a create/update call
var profileProperties = []interface{}{"title", "status", "display_name", "given_name", "family_name", "email", "directory_pool_id", "profile_id"}

func resourceTurbotProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotProfileCreate,
		Read:   resourceTurbotProfileRead,
		Update: resourceTurbotProfileUpdate,
		Delete: resourceTurbotProfileDelete,
		Exists: resourceTurbotProfileExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotProfileImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent but the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
			},
			// when doing a read, fetch the parent akas to use in suppressIfAkaMatches
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
			"profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"given_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"family_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"directory_pool_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"last_login_timestamp": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"middle_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"picture": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "Active",
			},
		},
	}
}

func resourceTurbotProfileExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	parentAka := d.Get("parent").(string)

	// build map of profile properties
	data := helpers.MapFromResourceData(d, profileProperties)
	// create profile returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateProfile(parentAka, data)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := helpers.StoreAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)

	return nil
}

func resourceTurbotProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build map of profile properties
	data := helpers.MapFromResourceData(d, folderDataProperties)

	// create profile returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateProfile(id, parentAka, data)
	if err != nil {
		return err
	}
	// set parent_akas property by loading resource and fetching the akas
	return helpers.StoreAkas(turbotMetadata.ParentId, "parent_akas", d, meta)
}

func resourceTurbotProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	profile, err := client.ReadProfile(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// profile was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("parent", profile.Parent)
	d.Set("title", profile.Title)
	d.Set("status", profile.Status)
	d.Set("email", profile.Email)
	d.Set("given_name", profile.GivenName)
	d.Set("family_name", profile.FamilyName)
	d.Set("directory_pool_id", profile.DirectoryPoolId)
	/// set parent_akas property by loading resource and fetching the akas
	return helpers.StoreAkas(profile.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotProfileDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotProfileImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotProfileRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
