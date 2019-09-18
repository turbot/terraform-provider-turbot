package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

func dataSourceTurbotResource() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurbotResourceRead,
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
			"resources": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeMap,
				},
			},
		},
	}
}

func dataSourceTurbotResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	// build required properties from payload
	properties, err := propertiesFromPayload(d.Get("payload").(string))
	if err != nil {
		return fmt.Errorf("error retrieving properties from resource payload: %s", err.Error())
	}

	resource, err := client.ReadResource(id, properties)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// resource was not found - clear id
			d.SetId("")
		}
		return err
	}

	// rebuild payload from the resource
	payload, err := payloadFromResource(resource.Data)
	if err != nil {
		return fmt.Errorf("error building resource payload: %s", err.Error())
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(resource.Turbot.ParentId, d, meta); err != nil {
		return err
	}
	d.Set("parent", resource.Turbot.ParentId)
	d.Set("payload", payload)

	return nil
}
