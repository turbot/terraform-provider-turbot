package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
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
			"resource_id": {
				Type:     schema.TypeString,
				Required: true,
				Forces:   New,
			},
			"resource_group_id": {
				Type:     schema.TypeString,
				Required: true,
				Forces:   New,
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
	resourceId := d.Get("resource_id").(string)
	resourceGroupId := d.Get("resource_group_id").(string)
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.CreateSmartFolderAttachment(resourceId, resourceGroupId)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)

	return nil
}

func resourceSmartFolderAttachemntRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	folder, err := client.ReadSmartFolderAttachment(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	d.Set("resource", SmartFolderAttachment.Resource)
	d.Set("smart_folder_attachment", SmartFolderAttachment.SmartFolderAttachment)

	return nil
}

func resourceTurbotSmartFolderAttachemntDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotSmartFolderAttachemntImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSmartFolderAttachemntRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
