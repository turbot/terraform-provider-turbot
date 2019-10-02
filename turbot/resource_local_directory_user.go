package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// properties which must be passed to a create/update call
var localDirectoryUserProperties = []string{"title", "email", "status", "display_name", "given_name", "middle_name", "family_name", "picture"}
var localDirectoryUserMetadataProperties = []string{"tags"}

func resourceTurbotLocalDirectoryUser() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotLocalDirectoryUserCreate,
		Read:   resourceTurbotLocalDirectoryUserRead,
		Update: resourceTurbotLocalDirectoryUserUpdate,
		Delete: resourceTurbotLocalDirectoryUserDelete,
		Exists: resourceTurbotLocalDirectoryUserExists,
		Importer: &schema.ResourceImporter{ //need to understand
			State: resourceTurbotLocalDirectoryUserImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent bu tthe config contains the aka,
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
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"email": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"display_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"given_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"middle_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"family_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"picture": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"password_timestamp": {
				Type:     schema.TypeString,
				Computed: true,
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

func resourceTurbotLocalDirectoryUserExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotLocalDirectoryUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	// build mutation payload
	payload := map[string]map[string]interface{}{
		"data":       mapFromResourceData(d, localDirectoryUserProperties),
		"turbotData": mapFromResourceData(d, localDirectoryUserMetadataProperties),
	}
	// CreateLocalDirectoryUser returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateLocalDirectoryUser(parentAka, payload)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(turbotMetadata.ParentId)
	if err != nil {
		return err
	}
	// assign parent_akas
	d.Set("parent_akas", parentAkas)
	// assign the id
	d.SetId(turbotMetadata.Id)
	return nil
}

func resourceTurbotLocalDirectoryUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	// build mutation payload
	payload := map[string]map[string]interface{}{
		"data":       mapFromResourceData(d, localDirectoryUserProperties),
		"turbotData": mapFromResourceData(d, localDirectoryUserMetadataProperties),
	}
	id := d.Id()
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateLocalDirectoryUserResource(id, parentAka, payload)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(turbotMetadata.ParentId)
	if err != nil {
		return err
	}
	// assign parent_akas
	d.Set("parent_akas", parentAkas)
	return nil
}

func resourceTurbotLocalDirectoryUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()
	localDirectoryUser, err := client.ReadLocalDirectoryUser(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}
	// assign results back into ResourceData
	// set parent_akas property by loading parent resource and fetching the akas
	parentAkas, err := client.GetResourceAkas(localDirectoryUser.Turbot.ParentId)
	if err != nil {
		return err
	}
	d.Set("parent_akas", parentAkas)

	d.Set("parent", localDirectoryUser.Parent)
	d.Set("title", localDirectoryUser.Title)
	d.Set("email", localDirectoryUser.Email)
	d.Set("status", localDirectoryUser.Status)
	d.Set("display_name", localDirectoryUser.DisplayName)
	d.Set("given_name", localDirectoryUser.GivenName)
	d.Set("middle_name", localDirectoryUser.MiddleName)
	d.Set("family_name", localDirectoryUser.FamilyName)
	d.Set("picture", localDirectoryUser.Picture)
	return nil
}

func resourceTurbotLocalDirectoryUserDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotLocalDirectoryUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotLocalDirectoryUserRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
