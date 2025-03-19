package turbot

import (
	"reflect"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

// properties which must be passed to a create/update call
var campaignProperties = []interface{}{"title", "description", "status", "recipients", "preview", "check", "draft", "enforce", "guardrails", "accounts", "akas"}

func getCampaignUpdateProperties() []interface{} {
	excludedProperties := []string{"guardrails", "preview", "check", "draft", "enforce", "akas"}
	return helpers.RemoveProperties(campaignProperties, excludedProperties)
}

func resourceTurbotCampaign() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotCampaignCreate,
		Read:   resourceTurbotCampaignRead,
		Update: resourceTurbotCampaignUpdate,
		Delete: resourceTurbotCampaignDelete,
		Exists: resourceTurbotCampaignExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotCampaignImport,
		},
		Schema: map[string]*schema.Schema{
			"guardrails": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"accounts": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"recipients": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"title": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"parent": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"check": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: phaseSchema(),
				},
			},
			"detach": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: phaseSchema(),
				},
			},
			"draft": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: phaseSchema(),
				},
			},
			"enforce": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: phaseSchema(),
				},
			},
			"preview": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: phaseSchema(),
				},
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

func phaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"transition_at": {
			Type:     schema.TypeString,
			Required: true,
		},
		"transition_notice": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"transition_when": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"warn_at": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func resourceTurbotCampaignExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotCampaignCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build map of folder properties
	input := mapFromResourceData(d, campaignProperties)

	phases := map[string]interface{}{}
	if !isNil(input["preview"]) {
		phases["preview"] = setPhaseAttribute(input, "preview")
		delete(input, "preview")
	}

	if !isNil(input["check"]) {
		phases["check"] = setPhaseAttribute(input, "check")
		delete(input, "check")
	}
	input["phases"] = phases

	// panic(fmt.Sprintf("HERE >>> %+v", input))

	campaign, err := client.CreateCampaign(input)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(campaign.Turbot.Id)

	// TODO Remove Read call once schema changes are In.
	return resourceTurbotCampaignRead(d, meta)
}

func resourceTurbotCampaignUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build map of folder properties
	input := mapFromResourceData(d, getCampaignUpdateProperties())
	input["id"] = id

	_, err := client.UpdateCampaign(input)
	if err != nil {
		return err
	}
	// set 'Read' Properties
	return resourceTurbotCampaignRead(d, meta)
}

func resourceTurbotCampaignRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	campaign, err := client.ReadCampaign(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	d.Set("description", campaign.Description)
	d.Set("status", campaign.Status)
	d.Set("title", campaign.Turbot.Title)
	d.Set("recipients", campaign.Recipients)
	d.Set("akas", campaign.Turbot.Akas)

	if len(campaign.Accounts.Items) > 0 {
		accounts := []string{}
		for _, item := range campaign.Accounts.Items {
			accounts = append(accounts, item.Turbot.Id)
		}
		d.Set("accounts", accounts)
	}
	if len(campaign.Guardrails.Items) > 0 {
		guardrails := []string{}
		for _, item := range campaign.Guardrails.Items {
			guardrails = append(guardrails, item.Turbot.Id)
		}
		d.Set("guardrails", guardrails)
	}
	if !isNil(campaign.Turbot.ParentId) {
		d.Set("parent", campaign.Turbot.ParentId)
	}

	d.Set("preview", []interface{}{campaign.Phases.Preview})
	d.Set("check", []interface{}{campaign.Phases.Check})

	return nil
}

func resourceTurbotCampaignDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotCampaignImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotCampaignRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func isNil(i interface{}) bool {
	if i == nil {
		return true
	}
	switch reflect.TypeOf(i).Kind() {
	case reflect.Ptr, reflect.Map, reflect.Chan, reflect.Slice:
		return reflect.ValueOf(i).IsNil()
	}
	return false
}

func setPhaseAttribute(input map[string]interface{}, attributeName string) map[string]interface{} {
	phaseList := input[attributeName].([]interface{})
	if len(phaseList) > 0 {
		phase := phaseList[0].(map[string]interface{})

		return map[string]interface{}{
			"transitionAt":     phase["transition_at"].(string),
			"transitionNotice": phase["transition_notice"].(string),
			"transitionWhen":   phase["transition_when"].(string),
			"warnAt":           phase["warn_at"].([]interface{}),
		}
	}
	return nil
}
