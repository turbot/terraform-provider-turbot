package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

// properties which must be passed to a create/update call
var groupProfileInputProperties = []interface{}{"directory", "title", "group_profile_id", "status"}

func getGroupProfileUpdateProperties() []interface{} {
	excludedProperties := []string{"directory", "group_profile_id"}
	return helpers.RemoveProperties(groupProfileInputProperties, excludedProperties)
}
func resourceTurbotGroupProfile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotGroupProfileCreate,
		Read:   resourceTurbotGroupProfileRead,
		Update: resourceTurbotGroupProfileUpdate,
		Delete: resourceTurbotGroupProfileDelete,
		Exists: resourceTurbotGroupProfileExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotGroupProfileImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent directory
			"directory": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"group_profile_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTurbotGroupProfileExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotGroupProfileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	input := mapFromResourceData(d, groupProfileInputProperties)
	// do create
	groupProfile, err := client.CreateGroupProfile(input)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(groupProfile.Turbot.Id)
	// assign results back into ResourceData
	d.Set("directory", groupProfile.Directory)
	d.Set("title", groupProfile.Title)
	d.Set("status", groupProfile.Status)
	d.Set("group_profile_id", groupProfile.GroupProfileId)
	return nil
}

func resourceTurbotGroupProfileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	groupProfile, err := client.ReadGroupProfile(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// profile was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("directory", groupProfile.Directory)
	d.Set("title", groupProfile.Title)
	d.Set("status", groupProfile.Status)
	d.Set("group_profile_id", groupProfile.GroupProfileId)
	return nil
}

func resourceTurbotGroupProfileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build mutation data
	input := mapFromResourceData(d, groupProfileInputProperties)
	input["id"] = d.Id()

	// do create
	groupProfile, err := client.UpdateGroupProfile(input)
	if err != nil {
		return err
	}

	// assign results back into ResourceData
	d.Set("title", groupProfile.Title)
	d.Set("status", groupProfile.Status)
	return nil
}

func resourceTurbotGroupProfileDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	err := client.DeleteGroupProfile(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")
	return nil
}

func resourceTurbotGroupProfileImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotGroupProfileRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
