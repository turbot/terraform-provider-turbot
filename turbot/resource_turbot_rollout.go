package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

// properties which must be passed to a create/update call
var rolloutProperties = []interface{}{"title", "description", "status", "recipients", "preview", "check", "draft", "enforce", "detach", "guardrails", "accounts", "akas"}

func getRolloutUpdateProperties() []interface{} {
	excludedProperties := []string{"guardrails", "preview", "check", "draft", "enforce", "detach", "akas"}
	return helpers.RemoveProperties(rolloutProperties, excludedProperties)
}

func resourceTurbotRollout() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotRolloutCreate,
		Read:   resourceTurbotRolloutRead,
		Update: resourceTurbotRolloutUpdate,
		Delete: resourceTurbotRolloutDelete,
		Exists: resourceTurbotRolloutExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotRolloutImport,
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
					Schema: rolloutPhaseSchema(),
				},
			},
			"detach": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: rolloutPhaseSchema(),
				},
			},
			"draft": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: rolloutDraftSchema(),
				},
			},
			"enforce": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: rolloutPhaseSchema(),
				},
			},
			"preview": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: rolloutPreviewSchema(),
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

func rolloutPhaseSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_at": {
			Type:     schema.TypeString,
			Required: true,
		},
		"start_notice": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"start_early_if": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"recipients": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
		"warn_at": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func rolloutDraftSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_at": {
			Type:     schema.TypeString,
			Required: true,
		},
	}
}

func rolloutPreviewSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"start_at": {
			Type:     schema.TypeString,
			Required: true,
		},
		"start_notice": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"start_early_if": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"recipients": {
			Type:     schema.TypeList,
			Optional: true,
			Elem:     &schema.Schema{Type: schema.TypeString},
		},
	}
}

func resourceTurbotRolloutExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotRolloutCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build map of rollout properties
	input := mapFromResourceData(d, rolloutProperties)

	// extract and set phase-related inputs
	phases := make(map[string]interface{})

	phaseKeys := []string{"check", "enforce", "detach"}
	for _, key := range phaseKeys {
		if !isNil(input[key]) {
			phases[key] = setPhaseAttribute(input, key)
			delete(input, key)
		}
	}

	// special handling for "draft"
	if !isNil(input["draft"]) {
		phases["draft"] = setDraftInputAttribute(input, "draft")
		delete(input, "draft")
	}

	if !isNil(input["preview"]) {
		phases["preview"] = setPreviewInputAttribute(input, "preview")
		delete(input, "preview")
	}

	input["phases"] = phases

	// create the rollout
	rollout, err := client.CreateRollout(input)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(rollout.Turbot.Id)

	// TODO Remove Read call once schema changes are In.
	return resourceTurbotRolloutRead(d, meta)
}

func resourceTurbotRolloutUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build map of folder properties
	input := mapFromResourceData(d, getRolloutUpdateProperties())
	input["id"] = id

	_, err := client.UpdateRollout(input)
	if err != nil {
		return err
	}
	// set 'Read' Properties
	return resourceTurbotRolloutRead(d, meta)
}

func resourceTurbotRolloutRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	rollout, err := client.ReadRollout(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// Set basic attributes
	d.Set("description", rollout.Description)
	d.Set("status", rollout.Status)
	d.Set("title", rollout.Turbot.Title)
	d.Set("recipients", rollout.Recipients)
	d.Set("akas", rollout.Turbot.Akas)

	// helper to extract Turbot.Id from a slice of items
	extractIds := func(items []struct {
		Turbot apiClient.TurbotResourceMetadata
	}) []string {
		ids := make([]string, len(items))
		for i, item := range items {
			ids[i] = item.Turbot.Id
		}
		return ids
	}

	if len(rollout.Accounts.Items) > 0 {
		d.Set("accounts", extractIds(rollout.Accounts.Items))
	}
	if len(rollout.Guardrails.Items) > 0 {
		d.Set("guardrails", extractIds(rollout.Guardrails.Items))
	}
	if !isNil(rollout.Turbot.ParentId) {
		d.Set("parent", rollout.Turbot.ParentId)
	}

	// helper to set phase if it's not nil
	setPhase := func(key string, phase interface{}) {
		if !isNil(phase) {
			d.Set(key, []interface{}{phase})
		}
	}

	setPhase("preview", rollout.Phases.Preview)
	setPhase("check", rollout.Phases.Check)
	setPhase("enforce", rollout.Phases.Enforce)
	setPhase("detach", rollout.Phases.Detach)
	setPhase("draft", rollout.Phases.Draft)

	return nil
}

func resourceTurbotRolloutDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotRolloutImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotRolloutRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

func setPhaseAttribute(input map[string]interface{}, attributeName string) map[string]interface{} {
	phaseList := input[attributeName].([]interface{})
	if len(phaseList) > 0 {
		phase := phaseList[0].(map[string]interface{})

		return map[string]interface{}{
			"startAt":      phase["start_at"].(string),
			"startNotice":  phase["start_notice"].(string),
			"startEarlyIf": phase["start_early_if"].(string),
			"warnAt":       phase["warn_at"].([]interface{}),
			"recipients":   phase["recipients"].([]interface{}),
		}
	}
	return nil
}

func setDraftInputAttribute(input map[string]interface{}, attributeName string) map[string]interface{} {
	draftInputs := input[attributeName].([]interface{})
	if len(draftInputs) > 0 {
		draft := draftInputs[0].(map[string]interface{})

		return map[string]interface{}{
			"startAt": draft["start_at"].(string),
		}
	}
	return nil
}

func setPreviewInputAttribute(input map[string]interface{}, attributeName string) map[string]interface{} {
	previewInputs := input[attributeName].([]interface{})
	if len(previewInputs) > 0 {
		preview := previewInputs[0].(map[string]interface{})

		return map[string]interface{}{
			"startAt":      preview["start_at"].(string),
			"startNotice":  preview["start_notice"].(string),
			"startEarlyIf": preview["start_early_if"].(string),
			"recipients":   preview["recipients"].([]interface{}),
		}
	}
	return nil
}
