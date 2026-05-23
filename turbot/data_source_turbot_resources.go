package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
)

func dataSourceTurbotResources() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurbotResourcesRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "The filter to apply to the list of resources.",
			},
			"ids": {
				Type:        schema.TypeList,
				Computed:    true,
				Description: "List of resource IDs matching the filter.",
				Elem:        &schema.Schema{Type: schema.TypeString},
			},
		},
	}
}

func dataSourceTurbotResourcesRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	filter := d.Get("filter").(string)
	resources, err := client.ReadResourceList(filter, nil)
	if err != nil {
		if errors.NotFoundError(err) {
			// setting was not found - clear id
			d.SetId("")
		}
		return err
	}

	ids := make([]string, 0, len(resources))
	for _, resource := range resources {
		ids = append(ids, resource.Turbot.Id)
	}

	d.SetId(filter)
	d.Set("ids", ids)
	d.Set("filter", d.Get("filter"))

	return nil
}
