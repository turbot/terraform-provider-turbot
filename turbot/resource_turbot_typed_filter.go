package turbot

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
)

var typedFilterCreateProperties = []interface{}{"title", "description", "filter", "type", "tags", "akas"}
var typedFilterUpdateProperties = []interface{}{"title", "filter", "akas"}

func resourceTurbotTypedFilter() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotTypedFilterCreate,
		Read:   resourceTurbotTypedFilterRead,
		Update: resourceTurbotTypedFilterUpdate,
		Delete: resourceTurbotTypedFilterDelete,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotTypedFilterImport,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"type": {
				Type:         schema.TypeString,
				Required:     true,
				ForceNew:     true,
				ValidateFunc: validateTypedFilterType,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"akas": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: suppressIfAkaRemoved(),
			},
		},
	}
}

func resourceTurbotTypedFilterCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	input := mapFromResourceData(d, typedFilterCreateProperties)

	typedFilter, err := client.CreateTypedFilter(input)
	if err != nil {
		return err
	}

	d.SetId(typedFilter.Turbot.Id)
	return resourceTurbotTypedFilterRead(d, meta)
}

func resourceTurbotTypedFilterUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	input := mapFromResourceData(d, typedFilterUpdateProperties)
	input["id"] = d.Id()
	input["description"] = d.Get("description")
	if tags := mergedTagsUpdate(d); len(tags) > 0 {
		input["tags"] = tags
	}

	if err := client.UpdateTypedFilter(input); err != nil {
		return err
	}
	return resourceTurbotTypedFilterRead(d, meta)
}

func resourceTurbotTypedFilterRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	typedFilter, err := client.ReadTypedFilter(d.Id())
	if err != nil {
		if apiClient.IsTypedFilterNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}

	d.Set("title", typedFilter.Title)
	d.Set("description", typedFilter.Description)
	d.Set("filter", typedFilter.Filter)
	d.Set("type", typedFilter.Type)
	d.Set("tags", typedFilter.Turbot.Tags)
	d.Set("akas", typedFilter.Turbot.Akas)
	return nil
}

func resourceTurbotTypedFilterDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	if err := client.DeleteTypedFilter(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceTurbotTypedFilterImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotTypedFilterRead(d, meta); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("typed filter not found")
	}
	return []*schema.ResourceData{d}, nil
}

func validateTypedFilterType(v interface{}, k string) ([]string, []error) {
	if s := v.(string); s != "Account" && s != "Resource" {
		return nil, []error{fmt.Errorf("%q must be either \"Account\" or \"Resource\", got %q", k, s)}
	}
	return nil, nil
}

// Typed filter and guardrail updates merge tags, and a key sent as null is removed, so a tag dropped from the
// configuration has to be sent explicitly as null.
func mergedTagsUpdate(d *schema.ResourceData) map[string]interface{} {
	oldTags, newTags := d.GetChange("tags")
	tags := map[string]interface{}{}
	for key, value := range newTags.(map[string]interface{}) {
		tags[key] = value
	}
	for key := range oldTags.(map[string]interface{}) {
		if _, kept := tags[key]; !kept {
			tags[key] = nil
		}
	}
	return tags
}
