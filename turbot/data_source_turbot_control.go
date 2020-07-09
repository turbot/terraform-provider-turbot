package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
)

func dataSourceTurbotControl() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurbotControlRead,
		Schema: map[string]*schema.Schema{
			"id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"uri": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"state": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"reason": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"details": {
				Type:     schema.TypeString,
				Computed: true,
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

func dataSourceTurbotControlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	controlId := d.Get("id").(string)
	controlUri, controlOk := d.GetOk("uri")
	resourceId, resourceOk := d.GetOk("resource")
	var args string

	if controlOk || resourceOk {
		if controlOk && resourceOk {
			args = fmt.Sprintf(`uri: %s, resourceId: %s`, controlUri, resourceId)
		}
		if !resourceOk {
			return fmt.Errorf("if controlUri is given, resourceId is required")
		}
		if !controlOk {
			return fmt.Errorf("if resourceId is given, controlUri is required")
		}
	}

	if args == "" {
		args = fmt.Sprintf(`id: "%s"`, controlId)
	}

	control, err := client.ReadControl(args)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// setting was not found - clear id
			d.SetId("")
		}
		return err
	}

	d.SetId(control.Turbot["id"])
	d.Set("uri", control.Type.Uri)
	d.Set("resource", control.Turbot["resourceId"])
	d.Set("state", control.State)
	d.Set("reason", control.Reason)
	d.Set("details", control.Details)
	d.Set("turbot", control.Turbot)
	return nil
}
