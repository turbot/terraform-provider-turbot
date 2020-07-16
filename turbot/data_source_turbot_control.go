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
		},
	}
}

func dataSourceTurbotControlRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	controlId, controlIdSet := d.GetOk("id")
	controlUri, controlUriSet := d.GetOk("uri")
	resourceId, resourceIdSet := d.GetOk("resource")
	var args string

	if controlIdSet {
		if controlUriSet || resourceIdSet {
			return fmt.Errorf("if controlId is set, controlUri and resourceId must not be set")
		}
		args = fmt.Sprintf(`id: "%s"`, controlId)
	} else {
		if !controlUriSet || !resourceIdSet {
			return fmt.Errorf("either controlId or controlUri AND resourceId must not be set")
		}
		args = fmt.Sprintf(`uri: %s, resourceId: %s`, controlUri, resourceId)
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
	return nil
}
