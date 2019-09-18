package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

func dataSourceTurbotPolicy() *schema.Resource {
	return &schema.Resource{
		Read: dataSourceTurbotPolicyValueRead,

		Schema: map[string]*schema.Schema{
			"resource": {
				Type:     schema.TypeString,
				Required: true,
			},
			"policy_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"value": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"value_source": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"precedence": {
				Type:     schema.TypeString,
				Computed: true,
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
			"setting_id": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}
func dataSourceTurbotPolicyValueRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	policyTypeUri := d.Get("policy_type").(string)
	resourceAka := d.Get("resource").(string)

	policyValue, err := client.ReadPolicyValue(policyTypeUri, resourceAka)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// setting was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.SetId(policyValue.Turbot.Id)

	d.Set("value", fmt.Sprintf("%v", policyValue.Value))
	d.Set("value_source", policyValue.Setting.ValueSource)
	d.Set("precedence", policyValue.Precedence)
	d.Set("state", policyValue.State)
	d.Set("reason", policyValue.Reason)
	d.Set("details", policyValue.Details)
	d.Set("setting_id", policyValue.Setting.Turbot.Id)
	return nil
}
