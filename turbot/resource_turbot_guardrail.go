package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
)

// properties which must be passed to a create/update call
var guardrailProperties = map[string]string{
	"title":       "title",
	"description": "description",
	"targets":     "targets",
	"controls":    "controlTypes",
	"akas":        "akas",
	"tags":        "tags",
}

func getGuardrailUpdateProperties() map[string]string {
	excludedProperties := []string{"controls"}

	// Remove the excluded properties from guardrailProperties
	for _, key := range excludedProperties {
		delete(guardrailProperties, key)
	}

	return guardrailProperties
}

func resourceTurbotGuardrail() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotGuardrailCreate,
		Read:   resourceTurbotGuardrailRead,
		Update: resourceTurbotGuardrailUpdate,
		Delete: resourceTurbotGuardrailDelete,
		Exists: resourceTurbotGuardrailExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotGuardrailImport,
		},
		Schema: map[string]*schema.Schema{
			"title": {
				Description: "The title of the guardrail.",
				Type:        schema.TypeString,
				Required:    true,
			},
			"description": {
				Description: "The description of the guardrail.",
				Type:        schema.TypeString,
				Optional:    true,
			},
			"targets": {
				Description: "The targets where the guardrail will be applied.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"controls": {
				Description: "The control types associated with the guardrail.",
				Type:        schema.TypeList,
				Required:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"tags": {
				Description: "The tags for the guardrail.",
				Type:        schema.TypeMap,
				Optional:    true,
			},
			"akas": {
				Description: "The akas of the guardrail.",
				Type:        schema.TypeList,
				Optional:    true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: suppressIfAkaRemoved(),
			},
		},
	}
}

func resourceTurbotGuardrailExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotGuardrailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build map of folder properties
	input := mapFromResourceDataWithPropertyMap(d, guardrailProperties)

	guardrail, err := client.CreateGuardrail(input)
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(guardrail.Turbot.Id)

	// TODO Remove Read call once schema changes are In.
	return resourceTurbotGuardrailRead(d, meta)
}

func resourceTurbotGuardrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	// build map of folder properties
	input := mapFromResourceDataWithPropertyMap(d, getGuardrailUpdateProperties())
	input["id"] = id

	_, err := client.UpdateGuardrail(input)
	if err != nil {
		return err
	}
	// set 'Read' Properties
	// TODO Remove Read call once schema changes are In.
	return resourceTurbotGuardrailRead(d, meta)
}

func resourceTurbotGuardrailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	guardrail, err := client.ReadGuardrail(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	d.Set("title", guardrail.Turbot.Title)
	d.Set("description", guardrail.Description)
	d.Set("akas", guardrail.Turbot.Akas)
	d.Set("tags", guardrail.Turbot.Tags)

	if len(guardrail.Targets.Items) > 0 {
		targets := []string{}
		for _, target := range guardrail.Targets.Items {
			targets = append(targets, target.Uri)
		}
		d.Set("targets", targets)
	}

	if len(guardrail.ControlTypes.Items) > 0 {
		controls := []string{}
		for _, control := range guardrail.ControlTypes.Items {
			controls = append(controls, control.Uri)
		}
		d.Set("controls", controls)
	}

	return nil
}

func resourceTurbotGuardrailDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotGuardrailImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotGuardrailRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
