package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
)

var controlProperties = map[string]string{
	"control_id":   "id",
	"resource":     "resource",
	"control_type": "controlType",
	"note":         "note",
	"to_timestamp": "toTimestamp",
	"until_states": "untilStates",
}

func getControlUpdateProperties() map[string]string {
	excludedProperties := []string{"control_id", "resource", "control_type"}

	// Remove the excluded properties from controlProperties
	for _, key := range excludedProperties {
		delete(controlProperties, key)
	}

	return controlProperties
}

func getControlDeleteProperties() map[string]string {
	excludedProperties := []string{"to_timestamp", "note", "until_states", "resource", "control_type"}

	// Remove the excluded properties from controlProperties
	for _, key := range excludedProperties {
		delete(controlProperties, key)
	}

	return controlProperties
}

func resourceTurbotControlMute() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotControlMuteCreate,
		Read:   resourceTurbotControlMuteRead,
		Update: resourceTurbotControlMuteUpdate,
		Delete: resourceTurbotControlMuteDelete,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotControlMuteImport,
		},
		Schema: map[string]*schema.Schema{
			"control_id": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID of the control to mute.",
			},
			"resource": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID or AKA of the resource where the control is available.",
			},
			"control_type": {
				Type:        schema.TypeString,
				Optional:    true,
				Computed:    true,
				Description: "The ID or AKA of the control type to be muted.",
			},
			"note": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "An optional note explaining the reason for muting the control.",
			},
			"to_timestamp": {
				Type:        schema.TypeString,
				Optional:    true,
				Description: "The timestamp specifying when the mute should end, in ISO 8601 format.",
			},
			"until_states": {
				Type:        schema.TypeList,
				Optional:    true,
				Description: "A list of control states specifying where the mute will not apply.",
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"state": {
				Type:        schema.TypeString,
				Computed:    true,
				Description: "The current state of the control.",
			},
		},
	}
}

// Mute a control
func resourceTurbotControlMuteCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build input map to pass to mutation
	if value, ok := d.GetOk("control_id"); ok {
		d.Set("control_id", value.(string))
		delete(controlProperties, "resource")
		delete(controlProperties, "control_type")
	}

	if value, ok := d.GetOk("resource"); ok {
		d.Set("resource", value.(string))
		delete(controlProperties, "id")
	}

	if value, ok := d.GetOk("control_type"); ok {
		d.Set("control_type", value.(string))
		delete(controlProperties, "id")
	}

	// build mutation input data by parsing the resource schema and
	// excluding top level properties which must not be sent to update - `type`
	input := mapFromResourceDataWithPropertyMap(d, controlProperties)

	muteControl, err := client.MuteControl(input)
	if err != nil {
		return err
	}

	// Set the id
	d.SetId(muteControl.Turbot["id"])

	d.Set("control_id", muteControl.Turbot["id"])
	d.Set("resource", muteControl.Turbot["resourceId"])
	d.Set("state", muteControl.State)
	d.Set("control_type", muteControl.Type.Uri)

	// Check control mute status
	muteConfig := muteControl.Mute.(map[string]interface{})
	if muteConfig["note"] != nil {
		d.Set("note", muteConfig["note"].(string))
	}

	if muteConfig["toTimestamp"] != nil {
		d.Set("to_timestamp", muteConfig["toTimestamp"].(string))
	}

	if muteConfig["untilStates"] != nil {
		d.Set("until_states", muteConfig["untilStates"].([]interface{}))
	}

	return nil
}

// Read control mute configuration
func resourceTurbotControlMuteRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// Get the control id
	controlId := d.Get("control_id").(string)

	control, err := client.ReadControl(controlId)
	if err != nil {
		if errors.NotFoundError(err) {
			// control was not found - clear id
			d.Set("id", "")
		}
		return err
	}

	d.SetId(control.Turbot["id"])
	d.Set("control_id", control.Turbot["id"])
	d.Set("control_type", control.Type.Uri)
	d.Set("resource", control.Turbot["resourceId"])
	d.Set("state", control.State)

	muteConfig := control.Mute.(map[string]interface{})
	if muteConfig["toTimestamp"] != nil {
		d.Set("to_timestamp", muteConfig["toTimestamp"].(string))
	}
	if muteConfig["untilStates"] != nil {
		d.Set("until_states", muteConfig["untilStates"].([]interface{}))
	}
	if muteConfig["note"] != nil {
		d.Set("note", muteConfig["note"].(string))
	}

	return nil
}

// Update control mute configuration
func resourceTurbotControlMuteUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	controlId := d.Get("control_id").(string)

	// build map of control properties
	input := mapFromResourceDataWithPropertyMap(d, getControlUpdateProperties())
	input["id"] = controlId

	muteControl, err := client.MuteControl(input)
	if err != nil {
		return err
	}
	// Set the id
	d.SetId(controlId)

	d.Set("control_id", muteControl.Turbot["id"])
	d.Set("resource", muteControl.Turbot["resourceId"])
	d.Set("state", muteControl.State)
	d.Set("control_type", muteControl.Type.Uri)

	// Check control mute status
	muteConfig := muteControl.Mute.(map[string]interface{})
	if muteConfig["note"] != nil {
		d.Set("note", muteConfig["note"].(string))
	}

	if muteConfig["toTimestamp"] != nil {
		d.Set("to_timestamp", muteConfig["toTimestamp"].(string))
	}

	if muteConfig["untilStates"] != nil {
		d.Set("until_states", muteConfig["untilStates"].([]interface{}))
	}

	return nil
}

// Unmute a control
func resourceTurbotControlMuteDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	controlId := d.Get("control_id").(string)

	// build map of control properties
	input := mapFromResourceDataWithPropertyMap(d, getControlDeleteProperties())
	input["id"] = controlId

	_, err := client.UnMuteControl(input)
	if err != nil {
		return err
	}

	return nil
}

func resourceTurbotControlMuteImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotControlMuteRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
