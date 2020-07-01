package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
)

func dataSourceTurbotResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurbotResourceRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"data": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"metadata": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": helpers.TagsSchema(),
			"akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"turbot": {
				Type:     schema.TypeMap,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func dataSourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resourceAka := d.Get("id").(string)
	resource, err := client.ReadSerializableResource(resourceAka)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// setting was not found - clear id
			d.SetId("")
		}
		return err
	}
	d.SetId(resource.Turbot["id"])
	d.Set("data", resource.Data)
	d.Set("metadata", resource.Metadata)
	d.Set("tags", helpers.TagsFromMap(resource.Tags))
	d.Set("akas", resource.Akas)
	d.Set("turbot", resource.Turbot)
	return nil
}
