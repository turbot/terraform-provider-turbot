package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"strings"
)

func resourceTurbotSmartFolderAttachemnt() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotSmartFolderAttachemntCreate,
		Read:   resourceTurbotSmartFolderAttachemntRead,
		Delete: resourceTurbotSmartFolderAttachemntDelete,
		Exists: resourceTurbotSmartFolderAttachemntExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotSmartFolderAttachemntImport,
		},
		Schema: map[string]*schema.Schema{
			"resource": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"smart_folder": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTurbotSmartFolderAttachemntExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotSmartFolderAttachemntCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	resource := d.Get("resource").(string)
	smartFolder := d.Get("smart_folder").(string)
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateSmartFolderAttachment(resource, smartFolder)
	if err != nil {
		return err
	}

	// assign the id
	var stateId = buildId(turbotMetadata.Id, resource)
	d.SetId(stateId)
	return nil
}

func resourceTurbotSmartFolderAttachemntRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	var id = strings.Split(d.Id(), "_")[0]
	_, err := client.ReadSmartFolderAttachment(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	//d.Set("resource", smartFolderAttachment.Resource)
	//d.Set("smart_folder", smartFolderAttachment.SmartFolderAttachment)

	return nil
}

func resourceTurbotSmartFolderAttachemntDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	resourceId := d.Get("resource_id").(string)
	resourceGroupId := d.Get("resource_group_id").(string)
	err := client.DeleteSmartFolderAttachment(resourceId, resourceGroupId)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")

	return nil
}

func resourceTurbotSmartFolderAttachemntImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSmartFolderAttachemntRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func buildId(smartFolder, resource string) string {
	return smartFolder + "_" + resource
}
