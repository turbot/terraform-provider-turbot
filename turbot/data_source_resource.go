package turbot

import (
	"encoding/json"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
)

func dataSourceTurbotResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurbotResourceRead,
		Schema: map[string]*schema.Schema{
			"aka": {
				Type:     schema.TypeString,
				Required: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"parent_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"json_data": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func dataSourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	aka := d.Get("aka").(string)

	resource, err := client.ReadFullResource(aka)
	if err != nil && !apiClient.NotFoundError(err) {
		return err
	}

	// assign results back into ResourceData
	if err != nil {
		return err
	}

	dataBytes, err := json.MarshalIndent(resource.Object, "", " ")
	if err != nil {
		return err
	}
	data := string(dataBytes)
	d.SetId(resource.Turbot.Id)
	d.Set("id", resource.Turbot.Id)
	d.Set("akas", resource.Turbot.Akas)
	d.Set("parent_id", resource.Turbot.ParentId)
	d.Set("json_data", data)
	return nil
}
