package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"log"
	"strings"
)

var smartFolderAttachProperties = map[string]string{
	"resource":     "resource",
	"smart_folder": "smartFolders",
}

func resourceTurbotSmartFolderAttachemnt() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotSmartFolderAttachmentCreate,
		Read:   resourceTurbotSmartFolderAttachmentRead,
		Delete: resourceTurbotSmartFolderAttachmentDelete,
		Exists: resourceTurbotSmartFolderAttachmentExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotSmartFolderAttachmentImport,
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

func resourceTurbotSmartFolderAttachmentExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	smartFolderId, resource := parseId(d.Id())
	// execute api call
	smartFolder, err := client.ReadSmartFolder(smartFolderId)
	if err != nil {
		return false, fmt.Errorf("error reading smart folder: %s", err.Error())
	}

	//find resource aka in list of attached resources
	for _, attachedResource := range smartFolder.AttachedResources.Items {
		log.Println("attachedResource", attachedResource)
		if resource == attachedResource.Turbot.Id {
			return true, nil
		}
		log.Println("Target resource", resource)
		log.Println("Resource id", attachedResource.Turbot.Id)
		for _, aka := range attachedResource.Turbot.Akas {
			if aka == resource {
				log.Println("Exists", aka)
				return true, nil
			}
		}
	}
	return false, nil
}

func resourceTurbotSmartFolderAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resource := d.Get("resource").(string)
	smartFolder := d.Get("smart_folder").(string)
	input := mapFromResourceDataWithPropertyMap(d, smartFolderAttachProperties)
	// create folder returns turbot resource metadata containing the id
	_, err := client.CreateSmartFolderAttachment(input)
	if err != nil {
		return err
	}

	// assign the id
	var stateId = buildId(smartFolder, resource)
	d.SetId(stateId)
	return nil
}

func resourceTurbotSmartFolderAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	smartFolder, resource := parseId(d.Id())
	// assign results back into ResourceData

	d.Set("resource", resource)
	d.Set("smart_folder", smartFolder)

	return nil
}

func resourceTurbotSmartFolderAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	input := mapFromResourceDataWithPropertyMap(d, smartFolderAttachProperties)
	err := client.DeleteSmartFolderAttachment(input)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")

	return nil
}

func resourceTurbotSmartFolderAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSmartFolderAttachmentRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func buildId(smartFolder, resource string) string {
	return smartFolder + "_" + resource
}

func parseId(id string) (smartFolder, resource string) {
	segments := strings.Split(id, "_")
	smartFolder = segments[0]
	resource = segments[1]
	return
}
