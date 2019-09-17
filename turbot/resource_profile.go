package turbot

import (
	// "fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
	// "strings"
)

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
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"external_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile_id": {
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
		},
	}
}

func resourceTurbotProfileExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	var payload = apiclient.ProfilePayload{
		Parent:          d.Get("parent").(string),
		Title:           d.Get("title").(string),
		Status:          d.Get("status").(string),
		DisplayName:     d.Get("display_name").(string),
		GivenName:       d.Get("given_name").(string),
		FamilyName:      d.Get("family_name").(string),
		Email:           d.Get("email").(string),
		DirectoryPoolId: d.Get("directory_pool_id").(string),
	}

	// create profile returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateProfile(&payload)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)

	return nil
}

func resourceTurbotProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	var payload = apiclient.ProfilePayload{
		Parent:          d.Get("parent").(string),
		Title:           d.Get("title").(string),
		Status:          d.Get("status").(string),
		DisplayName:     d.Get("display_name").(string),
		GivenName:       d.Get("given_name").(string),
		FamilyName:      d.Get("family_name").(string),
		Email:           d.Get("email").(string),
		DirectoryPoolId: d.Get("directory_pool_id").(string),
	}
	id := d.Id()

	log.Println("[INFO} resourceTurbotProfileUpdate", id, payload)
	// create profile returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateProfile(id, &payload)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}
	return nil
}

func resourceTurbotProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	profile, err := client.ReadProfile(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// profile was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(profile.Turbot.ParentId, d, meta); err != nil {
		return err
	}
	d.Set("parent", profile.Parent)
	d.Set("title", profile.Title)
	d.Set("status", profile.Status)
	d.Set("email", profile.Email)
	d.Set("given_name", profile.GivenName)
	d.Set("family_name", profile.FamilyName)
	d.Set("directory_pool_id", profile.DirectoryPoolId)

	return nil
}

func resourceTurbotProfileDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotProfileImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotProfileRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
