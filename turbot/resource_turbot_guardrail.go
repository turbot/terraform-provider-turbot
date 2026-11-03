package turbot

import (
	"fmt"
	"reflect"
	"sort"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

func guardrailValueSchema() map[string]*schema.Schema {
	return map[string]*schema.Schema{
		"value": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppressIfYamlEquivalent,
		},
		"template": {
			Type:     schema.TypeString,
			Optional: true,
		},
		"template_input": {
			Type:             schema.TypeString,
			Optional:         true,
			DiffSuppressFunc: suppressIfTemplateInputEquivalent,
		},
	}
}

func guardrailSettingSchema() *schema.Resource {
	fields := guardrailValueSchema()
	fields["type"] = &schema.Schema{
		Type:     schema.TypeString,
		Required: true,
	}
	return &schema.Resource{Schema: fields}
}

func resourceTurbotGuardrail() *schema.Resource {
	return &schema.Resource{
		Create:        resourceTurbotGuardrailCreate,
		Read:          resourceTurbotGuardrailRead,
		Update:        resourceTurbotGuardrailUpdate,
		Delete:        resourceTurbotGuardrailDelete,
		CustomizeDiff: resourceTurbotGuardrailCustomizeDiff,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotGuardrailImport,
		},
		Schema: map[string]*schema.Schema{
			"type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"akas": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
				DiffSuppressFunc: suppressIfAkaRemoved(),
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
			"mode": {
				Type:     schema.TypeString,
				Required: true,
			},
			"grace_period": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"action": {
				Type:     schema.TypeList,
				Optional: true,
				MaxItems: 1,
				Elem:     &schema.Resource{Schema: guardrailValueSchema()},
			},
			"account_filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"account_filter_id"},
			},
			"account_filter_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"account_filter"},
			},
			"resource_filter": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resource_filter_id"},
			},
			"resource_filter_id": {
				Type:          schema.TypeString,
				Optional:      true,
				ConflictsWith: []string{"resource_filter"},
			},
			"setting": {
				Type:     schema.TypeSet,
				Optional: true,
				Elem:     guardrailSettingSchema(),
			},
		},
	}
}

func resourceTurbotGuardrailCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	guardrailType, err := client.ReadGuardrailType(d.Get("type").(string))
	if err != nil {
		return err
	}
	input, err := guardrailCreateInput(d, guardrailType)
	if err != nil {
		return err
	}

	id, err := client.CreateGuardrail(input)
	if err != nil {
		return err
	}
	d.SetId(id)
	return resourceTurbotGuardrailRead(d, meta)
}

func resourceTurbotGuardrailRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	guardrail, err := client.ReadGuardrail(d.Id())
	if err != nil {
		if apiClient.IsGuardrailNotFound(err) {
			d.SetId("")
			return nil
		}
		return err
	}
	return setGuardrailState(d, guardrail)
}

func resourceTurbotGuardrailUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	var guardrailType *apiClient.GuardrailType
	if d.HasChange("setting") {
		var err error
		if guardrailType, err = client.ReadGuardrailType(d.Get("type").(string)); err != nil {
			return err
		}
	}
	input, err := guardrailUpdateInput(d, guardrailType)
	if err != nil {
		return err
	}

	if len(input) > 1 {
		if err := client.UpdateGuardrail(input); err != nil {
			return err
		}
	}
	return resourceTurbotGuardrailRead(d, meta)
}

func resourceTurbotGuardrailDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	if err := client.DeleteGuardrail(d.Id()); err != nil {
		return err
	}
	d.SetId("")
	return nil
}

func resourceTurbotGuardrailImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotGuardrailRead(d, meta); err != nil {
		return nil, err
	}
	if d.Id() == "" {
		return nil, fmt.Errorf("guardrail not found")
	}
	return []*schema.ResourceData{d}, nil
}

// The API ignores an empty or null account filter and an omitted action on update, so neither can
// be removed from an existing guardrail; a plan that removes one would never converge.
func resourceTurbotGuardrailCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {
	if settings, ok := d.Get("setting").(*schema.Set); ok {
		seen := map[string]bool{}
		for _, element := range settings.List() {
			name, _ := element.(map[string]interface{})["type"].(string)
			if seen[name] {
				return fmt.Errorf("setting %q is set more than once", name)
			}
			seen[name] = true
		}
	}
	if d.Id() == "" {
		return nil
	}

	oldAction, newAction := d.GetChange("action")
	if len(oldAction.([]interface{})) > 0 && len(newAction.([]interface{})) == 0 {
		return fmt.Errorf("the action block cannot be removed from a guardrail: the Guardrails API keeps a guardrail's action once set. Set the action you want, or replace the guardrail")
	}

	if !d.NewValueKnown("account_filter") || !d.NewValueKnown("account_filter_id") {
		return nil
	}
	oldFilter, newFilter := d.GetChange("account_filter")
	oldFilterId, newFilterId := d.GetChange("account_filter_id")
	if newFilter.(string) != "" || newFilterId.(string) != "" {
		return nil
	}
	for _, removed := range []struct{ name, old string }{{"account_filter", oldFilter.(string)}, {"account_filter_id", oldFilterId.(string)}} {
		if removed.old != "" {
			return fmt.Errorf("%s cannot be removed from a guardrail: the Guardrails API cannot clear account targeting. Set account_filter or account_filter_id, or replace the guardrail", removed.name)
		}
	}
	return nil
}

func guardrailCreateInput(d *schema.ResourceData, guardrailType *apiClient.GuardrailType) (map[string]interface{}, error) {
	input := map[string]interface{}{
		"type":  d.Get("type"),
		"title": d.Get("title"),
		"mode":  map[string]interface{}{"value": d.Get("mode")},
	}
	for terraformName, turbotName := range map[string]string{
		"description":        "description",
		"grace_period":       "gracePeriod",
		"account_filter":     "accountFilter",
		"account_filter_id":  "accountFilterId",
		"resource_filter":    "resourceFilter",
		"resource_filter_id": "resourceFilterId",
	} {
		if value := d.Get(terraformName).(string); value != "" {
			input[turbotName] = value
		}
	}
	if akas := d.Get("akas").([]interface{}); len(akas) > 0 {
		input["akas"] = akas
	}
	if tags := d.Get("tags").(map[string]interface{}); len(tags) > 0 {
		input["tags"] = tags
	}
	if action := d.Get("action").([]interface{}); len(action) > 0 {
		actionInput, err := guardrailActionInput(action[0].(map[string]interface{}))
		if err != nil {
			return nil, err
		}
		input["action"] = actionInput
	}

	blocks := guardrailSettingBlocks(d.Get("setting"))
	if err := validateGuardrailSettingNames(blocks, guardrailType); err != nil {
		return nil, err
	}
	var settings []interface{}
	for _, settingType := range guardrailType.PolicyTypes.Items {
		name := guardrailSettingName(guardrailType.Uri, settingType.Uri)
		if block, ok := blocks[name]; ok {
			setting, err := guardrailBlockInput(name, block)
			if err != nil {
				return nil, err
			}
			settings = append(settings, setting)
		} else if setting, ok := guardrailSettingDefault(name, &settingType); ok {
			settings = append(settings, setting)
		}
	}
	if len(settings) > 0 {
		input["settings"] = settings
	}
	return input, nil
}

// guardrailType is only read when settings changed: a setting whose block was removed is written
// back to the type's default.
func guardrailUpdateInput(d *schema.ResourceData, guardrailType *apiClient.GuardrailType) (map[string]interface{}, error) {
	input := map[string]interface{}{"id": d.Id()}

	for terraformName, turbotName := range map[string]string{"title": "title", "description": "description", "mode": "mode"} {
		if !d.HasChange(terraformName) {
			continue
		}
		if terraformName == "mode" {
			input[turbotName] = map[string]interface{}{"value": d.Get("mode")}
		} else {
			input[turbotName] = d.Get(terraformName)
		}
	}
	if d.HasChange("akas") {
		akas := d.Get("akas").([]interface{})
		if akas == nil {
			akas = []interface{}{}
		}
		input["akas"] = akas
	}
	if d.HasChange("tags") {
		if tags := mergedTagsUpdate(d); len(tags) > 0 {
			input["tags"] = tags
		}
	}
	if d.HasChange("grace_period") {
		if gracePeriod := d.Get("grace_period").(string); gracePeriod != "" {
			input["gracePeriod"] = gracePeriod
		} else {
			input["gracePeriod"] = nil
		}
	}
	if d.HasChange("action") {
		if action := d.Get("action").([]interface{}); len(action) > 0 {
			actionInput, err := guardrailActionInput(action[0].(map[string]interface{}))
			if err != nil {
				return nil, err
			}
			input["action"] = actionInput
		}
	}
	if d.HasChange("account_filter") || d.HasChange("account_filter_id") {
		if filter := d.Get("account_filter").(string); filter != "" {
			input["accountFilter"] = filter
		} else if filterId := d.Get("account_filter_id").(string); filterId != "" {
			input["accountFilterId"] = filterId
		}
	}
	if d.HasChange("resource_filter") || d.HasChange("resource_filter_id") {
		if filterId := d.Get("resource_filter_id").(string); filterId != "" {
			input["resourceFilterId"] = filterId
		} else {
			input["resourceFilter"] = d.Get("resource_filter").(string)
		}
	}
	if d.HasChange("setting") {
		settings, err := guardrailUpdateSettings(d, guardrailType)
		if err != nil {
			return nil, err
		}
		if len(settings) > 0 {
			input["settings"] = settings
		}
	}
	return input, nil
}

func guardrailUpdateSettings(d *schema.ResourceData, guardrailType *apiClient.GuardrailType) ([]interface{}, error) {
	oldSettings, newSettings := d.GetChange("setting")
	oldBlocks := guardrailSettingBlocks(oldSettings)
	newBlocks := guardrailSettingBlocks(newSettings)
	if err := validateGuardrailSettingNames(newBlocks, guardrailType); err != nil {
		return nil, err
	}
	settingTypes := guardrailSettingTypesByName(guardrailType)

	var settings []interface{}
	for _, name := range sortedKeys(newBlocks) {
		if oldBlock, ok := oldBlocks[name]; ok && guardrailBlocksEqual(oldBlock, newBlocks[name]) {
			continue
		}
		setting, err := guardrailBlockInput(name, newBlocks[name])
		if err != nil {
			return nil, err
		}
		settings = append(settings, setting)
	}
	for _, name := range sortedKeys(oldBlocks) {
		if _, kept := newBlocks[name]; kept {
			continue
		}
		setting, ok := guardrailSettingDefault(name, settingTypes[name])
		if !ok {
			return nil, fmt.Errorf("setting %q has no default in guardrail type %s, so its block cannot be removed: give it a value instead", name, guardrailType.Uri)
		}
		settings = append(settings, setting)
	}
	return settings, nil
}

func setGuardrailState(d *schema.ResourceData, guardrail *apiClient.Guardrail) error {
	d.Set("type", guardrail.Type.Uri)
	d.Set("title", guardrail.Title)
	d.Set("description", guardrail.Description)
	d.Set("grace_period", guardrail.GracePeriod)
	d.Set("account_filter", guardrail.AccountFilter)
	d.Set("account_filter_id", guardrail.AccountFilterId)
	d.Set("resource_filter", guardrail.ResourceFilter)
	d.Set("resource_filter_id", guardrail.ResourceFilterId)
	d.Set("tags", guardrail.Turbot.Tags)
	d.Set("akas", guardrail.Turbot.Akas)

	priorBlocks := guardrailSettingBlocks(d.Get("setting"))
	var priorAction map[string]interface{}
	if action := d.Get("action").([]interface{}); len(action) > 0 && action[0] != nil {
		priorAction = action[0].(map[string]interface{})
	}
	settingTypes := guardrailSettingTypesByName(guardrail.Type)

	settings := []interface{}{}
	action := []interface{}{}
	for _, row := range guardrail.Settings.Items {
		name := guardrailSettingName(guardrail.Type.Uri, row.Type.Uri)
		switch name {
		case "mode":
			d.Set("mode", helpers.InterfaceToString(row.Value))
		case "action":
			action = append(action, guardrailBlockFromRow(row, priorAction))
		default:
			// A setting with no block takes the type's default, so a default-valued setting enters
			// state only when the configuration has a block for it. A non-default value with no
			// block is drift, and shows as a block to remove.
			priorBlock, inPrior := priorBlocks[name]
			if !inPrior && guardrailRowIsDefault(row, settingTypes[name]) {
				continue
			}
			block := guardrailBlockFromRow(row, priorBlock)
			block["type"] = name
			settings = append(settings, block)
		}
	}
	if err := d.Set("setting", settings); err != nil {
		return err
	}
	return d.Set("action", action)
}

// guardrailSettingName is the short name the API accepts as a setting's type: the setting's policy
// type name with the guardrail type's name stripped and its first letter lowered, so
// bucketVersioningGuardrailVersioning on bucketVersioningGuardrail is "versioning".
func guardrailSettingName(guardrailTypeUri, settingTypeUri string) string {
	typeName := guardrailTypeUri[strings.LastIndex(guardrailTypeUri, "/")+1:]
	segment := settingTypeUri[strings.LastIndex(settingTypeUri, "/")+1:]
	name := strings.TrimPrefix(segment, typeName)
	if name == "" || name == segment {
		return segment
	}
	return strings.ToLower(name[:1]) + name[1:]
}

func guardrailSettingTypesByName(guardrailType *apiClient.GuardrailType) map[string]*apiClient.GuardrailSettingType {
	byName := map[string]*apiClient.GuardrailSettingType{}
	if guardrailType == nil {
		return byName
	}
	for i := range guardrailType.PolicyTypes.Items {
		settingType := &guardrailType.PolicyTypes.Items[i]
		byName[guardrailSettingName(guardrailType.Uri, settingType.Uri)] = settingType
	}
	return byName
}

func validateGuardrailSettingNames(blocks map[string]map[string]interface{}, guardrailType *apiClient.GuardrailType) error {
	settingTypes := guardrailSettingTypesByName(guardrailType)
	for _, name := range sortedKeys(blocks) {
		if _, ok := settingTypes[name]; !ok {
			return fmt.Errorf("unknown setting %q for guardrail type %s. Valid settings: %s", name, guardrailType.Uri, strings.Join(sortedKeys(settingTypes), ", "))
		}
	}
	return nil
}

func guardrailSettingBlocks(settings interface{}) map[string]map[string]interface{} {
	blocks := map[string]map[string]interface{}{}
	set, ok := settings.(*schema.Set)
	if !ok || set == nil {
		return blocks
	}
	for _, element := range set.List() {
		block := element.(map[string]interface{})
		blocks[block["type"].(string)] = block
	}
	return blocks
}

func guardrailBlockInput(name string, block map[string]interface{}) (map[string]interface{}, error) {
	input, err := guardrailValueInput(block, "valueSource")
	if err != nil {
		return nil, fmt.Errorf("setting %q: %s", name, err)
	}
	input["type"] = name
	return input, nil
}

func guardrailActionInput(block map[string]interface{}) (map[string]interface{}, error) {
	input, err := guardrailValueInput(block, "value")
	if err != nil {
		return nil, fmt.Errorf("action: %s", err)
	}
	return input, nil
}

func guardrailValueInput(block map[string]interface{}, valueField string) (map[string]interface{}, error) {
	value, _ := block["value"].(string)
	template, _ := block["template"].(string)
	templateInput, _ := block["template_input"].(string)
	switch {
	case template != "" && value != "":
		return nil, fmt.Errorf("set either value, or template and template_input, not both")
	case template != "":
		return map[string]interface{}{"template": template, "templateInput": templateInput}, nil
	case templateInput != "":
		return nil, fmt.Errorf("template_input needs a template")
	case value == "":
		return nil, fmt.Errorf("set value, or template and template_input")
	}
	return map[string]interface{}{valueField: value}, nil
}

func guardrailSettingDefault(name string, settingType *apiClient.GuardrailSettingType) (map[string]interface{}, bool) {
	if settingType == nil {
		return nil, false
	}
	if value, ok := settingType.Schema["default"]; ok {
		return map[string]interface{}{"type": name, "value": value}, true
	}
	if settingType.DefaultTemplate != "" {
		return map[string]interface{}{"type": name, "template": settingType.DefaultTemplate, "templateInput": settingType.DefaultTemplateInput}, true
	}
	return nil, false
}

func guardrailRowIsDefault(row apiClient.GuardrailSetting, settingType *apiClient.GuardrailSettingType) bool {
	if settingType == nil {
		return false
	}
	if value, ok := settingType.Schema["default"]; ok {
		return row.Template == "" && reflect.DeepEqual(row.Value, value)
	}
	if settingType.DefaultTemplate == "" || row.Template != settingType.DefaultTemplate {
		return false
	}
	rowInput, _ := helpers.InterfaceToStringOrYaml(row.TemplateInput)
	defaultInput, _ := helpers.InterfaceToStringOrYaml(settingType.DefaultTemplateInput)
	return yamlEquivalent(rowInput, defaultInput)
}

// The prior value is kept when the server's is equivalent, so a value written in a different but
// equal form does not show as a change.
func guardrailBlockFromRow(row apiClient.GuardrailSetting, prior map[string]interface{}) map[string]interface{} {
	priorString := func(key string) string {
		value, _ := prior[key].(string)
		return value
	}
	block := map[string]interface{}{"value": "", "template": "", "template_input": ""}
	if row.Template != "" {
		templateInput, _ := helpers.InterfaceToStringOrYaml(row.TemplateInput)
		if yamlEquivalent(priorString("template_input"), templateInput) {
			templateInput = priorString("template_input")
		}
		block["template"] = row.Template
		block["template_input"] = templateInput
		return block
	}
	value := strings.TrimSuffix(row.ValueSource, "\n")
	if row.ValueSource == "" {
		value = helpers.InterfaceToString(row.Value)
	}
	if priorValue := priorString("value"); priorValue != "" && yamlEquivalent(priorValue, value) {
		value = priorValue
	}
	block["value"] = value
	return block
}

func guardrailBlocksEqual(a, b map[string]interface{}) bool {
	return yamlEquivalent(a["value"].(string), b["value"].(string)) &&
		a["template"] == b["template"] &&
		yamlEquivalent(a["template_input"].(string), b["template_input"].(string))
}

func suppressIfYamlEquivalent(_, old, new string, _ *schema.ResourceData) bool {
	return old != "" && yamlEquivalent(old, new)
}

func yamlEquivalent(a, b string) bool {
	if a == b {
		return true
	}
	if a == "" || b == "" {
		return false
	}
	equal, err := helpers.YamlStringsAreEqual(a, b)
	return err == nil && equal
}

func sortedKeys[V any](m map[string]V) []string {
	keys := make([]string, 0, len(m))
	for key := range m {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}
