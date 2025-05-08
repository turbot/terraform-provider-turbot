package turbot

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

var fileProperties = []interface{}{"parent", "tags", "akas"}

func resourceTurbotFile() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotFileCreate,
		Read:   resourceTurbotFileRead,
		Update: resourceTurbotFileUpdate,
		Delete: resourceTurbotFileDelete,
		Exists: resourceTurbotFileExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotFileImport,
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"content": {
				Type:             schema.TypeString,
				Optional:         true,
				DiffSuppressFunc: suppressIfDataMatches,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"akas": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTurbotFileExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotFileCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	title := d.Get("title")
	description := d.Get("description")
	var err error
	// build input map to pass to mutation
	input, err := buildFileInput(d, fileProperties)
	if err != nil {
		return err
	}
	// set type property
	input["type"] = "tmod:@turbot/turbot#/resource/types/file"

	turbotMetadata, err := client.CreateResource(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)
	// save the formatted data: this is to ensure the acceptance tests behave in a consistent way regardless of the ordering of the json data
	d.Set("content", helpers.FormatJson(d.Get("content").(string)))
	d.Set("title", title)
	d.Set("description", description)
	return nil
}

func resourceTurbotFileRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	resource, err := client.ReadFullResource(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}
	// rebuild content from the resource
	if len(resource.Data) != 0 {
		content, err := helpers.MapToJsonString(resource.Data)
		if err != nil {
			return fmt.Errorf("error building resource content: %s", err.Error())
		}
		d.Set("content", content)
	}

	customMetadata := resource.Turbot.Custom
	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(resource.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}

	if v, ok := customMetadata["title"]; ok {
		d.Set("title", v)
	}
	if v, ok := customMetadata["description"]; ok {
		d.Set("description", v)
	}
	// assign results back into ResourceData
	d.Set("parent", resource.Turbot.ParentId)
	return nil
}

func resourceTurbotFileUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build input map to pass to mutation
	id := d.Id()

	input, err := buildFileInput(d, fileProperties)
	if err != nil {
		return err
	}

	input["data"], err = buildInputDataMap(d)
	if err != nil {
		return err
	}
	input["id"] = id

	// Check if parent actually changed
	oldParent, newParent := d.GetChange("parent")
	if oldParent != newParent {
		// Only update parent if it actually changed
		if err := updateParent(client, id, newParent.(string)); err != nil {
			return err
		}
	}

	// Delete `parent` from input because the putResource mutation does not expect `parent` in the input
	delete(input, "parent")

	turbotMetadata, err := client.PutResource(input)
	if err != nil {
		return err
	}
	// save the formatted data: this is to ensure the acceptance tests behave in a consistent way regardless of the ordering of the json data
	d.Set("content", helpers.FormatJson(d.Get("content").(string)))

	metadataMap := turbotMetadata.Custom
	if v, ok := metadataMap["description"]; ok {
		d.Set("description", v)
	}
	d.Set("title", metadataMap["title"])
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta)
}

func resourceTurbotFileDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotFileImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotResourceRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// Helper function for parent updates
func updateParent(client *apiClient.Client, id, parent string) error {
	_, err := client.UpdateResource(map[string]interface{}{
		"id":     id,
		"parent": parent,
		"data":   map[string]interface{}{},
	})
	if err != nil {
		return fmt.Errorf("failed to update parent: %w", err)
	}
	return nil
}

func buildFileInput(d *schema.ResourceData, properties []interface{}) (map[string]interface{}, error) {
	// initialize empty input map
	var err error
	var input = make(map[string]interface{})

	input = mapFromResourceData(d, properties)
	// convert data from json string to map
	// empty `data` object to handle no `content` given in config
	input["data"] = make(map[string]string)
	if contentString, ok := d.GetOk("content"); ok {
		if input["data"], err = helpers.JsonStringToMap(contentString.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal content: \n%s\nerror: %s", contentString, err.Error())
		}
	}
	input["metadata"] = buildInputMetadataMap(d)
	return input, nil
}

func buildInputMetadataMap(d *schema.ResourceData) map[string]interface{} {
	var metadataMap = make(map[string]interface{})
	title := d.Get("title")
	metadataMap["title"] = title

	old, new := d.GetChange("description")

	if old != "" && new == "" {
		metadataMap["description"] = nil
	}
	if new != "" {
		metadataMap["description"] = new
	}
	return metadataMap
}

// build the data field of the mutation input from content attribute
// any property which should be removed must be explicitly set to null in the mutation input
func buildInputDataMap(d *schema.ResourceData) (map[string]interface{}, error) {
	var oldContent, newContent map[string]interface{}
	var err error
	// fetch old(state-file) and new(config) content
	if old, new := d.GetChange("content"); old != nil {
		if oldContent, err = helpers.JsonStringToMap(old.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal content: \n%s\nerror: %s", old.(string), err.Error())
		}
		if newContent, err = helpers.JsonStringToMap(new.(string)); err != nil {
			return nil, fmt.Errorf("error build resource mutation input, failed to unmarshal content: \n%s\nerror: %s", new.(string), err.Error())
		}
		// extract keys from old content not in new
		excludeContentProperties := helpers.GetOldMapProperties(oldContent, newContent)
		for _, key := range excludeContentProperties {
			// set keys of old content to `nil` in new content
			// any property which doesn't exist in config is set to nil
			if _, ok := oldContent[key.(string)]; ok {
				newContent[key.(string)] = nil
			}
		}
	}
	return newContent, nil
}
