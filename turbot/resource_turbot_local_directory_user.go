package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
)

// properties which must be passed to a create/update call
var localDirectoryUserInputProperties = []interface{}{"parent", "tags"}
var localDirectoryUserDataProperties = []interface{}{"title", "email", "status", "display_name", "given_name", "middle_name", "family_name", "picture"}

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
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotLocalDirectoryUserCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build mutation input
	input := mapFromResourceData(d, localDirectoryUserInputProperties)
	data := mapFromResourceData(d, localDirectoryUserDataProperties)
	// set computed properties
	data["status"] = "Active"
	input["data"] = data

	// do create
	localDirectoryUser, err := client.CreateLocalDirectoryUser(input)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(localDirectoryUser.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(localDirectoryUser.Turbot.Id)

	d.Set("parent", localDirectoryUser.Parent)
	d.Set("title", localDirectoryUser.Title)
	d.Set("email", localDirectoryUser.Email)
	d.Set("display_name", localDirectoryUser.DisplayName)
	d.Set("given_name", localDirectoryUser.GivenName)
	d.Set("middle_name", localDirectoryUser.MiddleName)
	d.Set("family_name", localDirectoryUser.FamilyName)
	d.Set("picture", localDirectoryUser.Picture)
	// set the calculated properties
	d.Set("status", data["status"])
	return nil
}

func resourceTurbotLocalDirectoryUserUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build mutation payload
	input := mapFromResourceData(d, localDirectoryUserInputProperties)
	input["data"] = mapFromResourceData(d, localDirectoryUserDataProperties)
	input["id"] = d.Id()

	// do update
	localDirectoryUser, err := client.UpdateLocalDirectoryUserResource(input)
	if err != nil {
		return err
	}
	d.Set("parent", localDirectoryUser.Parent)
	d.Set("title", localDirectoryUser.Title)
	d.Set("email", localDirectoryUser.Email)
	d.Set("status", localDirectoryUser.Status)
	d.Set("display_name", localDirectoryUser.DisplayName)
	d.Set("given_name", localDirectoryUser.GivenName)
	d.Set("middle_name", localDirectoryUser.MiddleName)
	d.Set("family_name", localDirectoryUser.FamilyName)
	d.Set("picture", localDirectoryUser.Picture)
	// set parent_akas property by loading parent resource and fetching the akas
	return storeAkas(localDirectoryUser.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotLocalDirectoryUserRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	localDirectoryUser, err := client.ReadLocalDirectoryUser(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}
	// assign results back into ResourceData
	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(localDirectoryUser.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}

	d.Set("parent", localDirectoryUser.Parent)
	d.Set("title", localDirectoryUser.Title)
	d.Set("email", localDirectoryUser.Email)
	d.Set("stat	1us", localDirectoryUser.Status)
	d.Set("display_name", localDirectoryUser.DisplayName)
	d.Set("given_name", localDirectoryUser.GivenName)
	d.Set("middle_name", localDirectoryUser.MiddleName)
	d.Set("family_name", localDirectoryUser.FamilyName)
	d.Set("picture", localDirectoryUser.Picture)
	return nil
}

func resourceTurbotLocalDirectoryUserDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotLocalDirectoryUserImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotLocalDirectoryUserRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
