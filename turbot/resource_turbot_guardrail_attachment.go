package turbot

import (
	"fmt"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
)

func resourceTurbotGuardrailAttachment() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotGuardrailAttachmentCreate,
		Read:   resourceTurbotGuardrailAttachmentRead,
		Delete: resourceTurbotGuardrailAttachmentDelete,
		Exists: resourceTurbotGuardrailAttachmentExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotGuardrailAttachmentImport,
		},
		Schema: map[string]*schema.Schema{
			"resource": {
				Type:             schema.TypeString,
				Required:         true,
				ForceNew:         true,
				DiffSuppressFunc: suppressIfAkaMatches("resource_akas"),
			},
			"guardrail": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"guardrail_phase": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"resource_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
		},
	}
}

func resourceTurbotGuardrailAttachmentExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	guardrailId, resource := parseSmartFolderId(d.Id())

	// execute api call
	guardrail, err := client.ReadGuardrail(guardrailId)
	if err != nil {
		return false, fmt.Errorf("error reading guardrail: %s", err.Error())
	}

	//find resource aka in list of attached resources
	for _, attachedResource := range guardrail.Accounts.Items {
		if resource == attachedResource.Turbot.Id {
			return true, nil
		}

		for _, aka := range attachedResource.Turbot.Akas {
			if aka == resource {
				return true, nil
			}
		}
	}
	return false, nil
}

func resourceTurbotGuardrailAttachmentCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resource := d.Get("resource").(string)
	guardrail := d.Get("guardrail").(string)
	guardrailPhase := d.Get("guardrail_phase").(string)

	input := map[string]interface{}{
		"resource": resource,
		"guardrailsWithPhase": []map[string]interface{}{
			{
				"guardrail": guardrail,
				"phase":     guardrailPhase,
			},
		},
	}

	_, err := client.AttachGuardrail(input)
	if err != nil {
		return err
	}

	// set resource_akas property by loading resource and fetching the akas
	if err := storeAkas(resource, "resource_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	var stateId = buildId(guardrail, resource)
	d.SetId(stateId)
	d.Set("resource", resource)
	d.Set("guardrail", guardrail)
	d.Set("guardrail_phase", guardrailPhase)
	return nil
}

func resourceTurbotGuardrailAttachmentRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// NOTE: This will not be called if the attachment does not exist
	guardrail, resource := parseSmartFolderId(d.Id())
	guardrailPhase := d.Get("guardrail_phase").(string)

	turbotResource, err := client.ReadResource(resource, nil)
	if err != nil {
		return err
	}
	// set resource_akas property by loading resource and fetching the akas
	if err := storeAkas(turbotResource.Turbot.Id, "resource_akas", d, meta); err != nil {
		return err
	}
	// assign results directly back into ResourceData
	d.Set("resource", resource)
	d.Set("guardrail", guardrail)
	d.Set("guardrail_phase", guardrailPhase)
	return nil
}

func resourceTurbotGuardrailAttachmentDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	resource := d.Get("resource").(string)
	guardrail := d.Get("guardrail").(string)
	input := map[string]interface{}{
		"resource":   resource,
		"guardrails": []string{guardrail},
	}
	err := client.DetachGuardrail(input)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")
	return nil
}

func resourceTurbotGuardrailAttachmentImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotGuardrailAttachmentRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
