package turbot

import (
	"reflect"
	"strings"
	"testing"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
)

const testGuardrailEncryptionType = "tmod:@turbot/aws-s3#/guardrail/types/bucketEncryptionAtRestGuardrail"

func testGuardrailEncryptionTypeDefinition() *apiClient.GuardrailType {
	prefix := "tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrail"
	guardrailType := &apiClient.GuardrailType{Uri: testGuardrailEncryptionType}
	guardrailType.PolicyTypes.Items = []apiClient.GuardrailSettingType{
		{Uri: prefix + "AcceptStrongerEncryption", Schema: map[string]interface{}{"type": "boolean", "default": true}},
		{Uri: prefix + "CustomerManagedKey", DefaultTemplate: "{{ $.defaultKey }}", DefaultTemplateInput: "{\n  defaultKey: policy(uri: \"aws-kms#/policy/types/keyDefaultCustomerManagedKey\")\n}\n"},
		{Uri: prefix + "EncryptionLevel", Schema: map[string]interface{}{"type": "string", "default": "AWS managed key"}},
	}
	return guardrailType
}

func testGuardrailRaw(extra map[string]interface{}) map[string]interface{} {
	raw := map[string]interface{}{
		"type":            testGuardrailEncryptionType,
		"title":           "tf-acc-unit",
		"mode":            "Check",
		"account_filter":  "title:tf-acc-no-such-account",
		"resource_filter": "$.Name:tf-acc-no-such-bucket",
	}
	for key, value := range extra {
		raw[key] = value
	}
	return raw
}

// testGuardrailUpdateData builds the ResourceData an Update sees: the old configuration as state,
// diffed against the new configuration.
func testGuardrailUpdateData(t *testing.T, oldRaw, newRaw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	r := resourceTurbotGuardrail()
	old := schema.TestResourceDataRaw(t, r.Schema, oldRaw)
	old.SetId("401101108545745")
	state := old.State()
	diff, err := r.Diff(state, terraform.NewResourceConfigRaw(newRaw), nil)
	if err != nil {
		t.Fatalf("diff: %s", err)
	}
	d, err := schema.InternalMap(r.Schema).Data(state, diff)
	if err != nil {
		t.Fatalf("data: %s", err)
	}
	return d
}

func settingsByType(t *testing.T, input map[string]interface{}) map[string]map[string]interface{} {
	t.Helper()
	byType := map[string]map[string]interface{}{}
	settings, _ := input["settings"].([]interface{})
	for _, setting := range settings {
		s := setting.(map[string]interface{})
		byType[s["type"].(string)] = s
	}
	return byType
}

func TestGuardrailSettingName(t *testing.T) {
	cases := map[string]string{
		"tmod:@turbot/aws-s3#/policy/types/bucketVersioningGuardrailVersioning": "versioning",
		"tmod:@turbot/aws-s3#/policy/types/bucketVersioningGuardrailMode":       "mode",
		"tmod:@turbot/aws-s3#/policy/types/bucketVersioningGuardrailAction":     "action",
		"tmod:@turbot/aws-s3#/policy/types/somethingElse":                       "somethingElse",
	}
	for uri, want := range cases {
		if got := guardrailSettingName("tmod:@turbot/aws-s3#/guardrail/types/bucketVersioningGuardrail", uri); got != want {
			t.Errorf("%s: got %q, want %q", uri, got, want)
		}
	}
}

func TestGuardrailCreateInputFillsDefaultsForSettingsWithNoBlock(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTurbotGuardrail().Schema, testGuardrailRaw(map[string]interface{}{
		"setting": []interface{}{map[string]interface{}{"type": "encryptionLevel", "value": "S3 managed key"}},
	}))

	input, err := guardrailCreateInput(d, testGuardrailEncryptionTypeDefinition())
	if err != nil {
		t.Fatal(err)
	}

	settings := settingsByType(t, input)
	want := map[string]map[string]interface{}{
		"acceptStrongerEncryption": {"type": "acceptStrongerEncryption", "value": true},
		"customerManagedKey": {"type": "customerManagedKey", "template": "{{ $.defaultKey }}",
			"templateInput": "{\n  defaultKey: policy(uri: \"aws-kms#/policy/types/keyDefaultCustomerManagedKey\")\n}\n"},
		"encryptionLevel": {"type": "encryptionLevel", "valueSource": "S3 managed key"},
	}
	if !reflect.DeepEqual(settings, want) {
		t.Fatalf("settings:\n got %#v\nwant %#v", settings, want)
	}
	if !reflect.DeepEqual(input["mode"], map[string]interface{}{"value": "Check"}) {
		t.Errorf("mode: %#v", input["mode"])
	}
	if input["accountFilter"] != "title:tf-acc-no-such-account" || input["resourceFilter"] != "$.Name:tf-acc-no-such-bucket" {
		t.Errorf("targeting: %#v", input)
	}
	for _, absent := range []string{"action", "accountFilterId", "resourceFilterId", "gracePeriod", "tags", "akas"} {
		if _, ok := input[absent]; ok {
			t.Errorf("%s must not be sent when unset, got %#v", absent, input[absent])
		}
	}
}

func TestGuardrailCreateInputCalculatedSettingAndAction(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTurbotGuardrail().Schema, testGuardrailRaw(map[string]interface{}{
		"setting": []interface{}{map[string]interface{}{"type": "acceptStrongerEncryption", "template": "false", "template_input": "{ resource { turbot { id } } }"}},
		"action":  []interface{}{map[string]interface{}{"template": "x", "template_input": "{ resource { turbot { id } } }"}},
	}))

	input, err := guardrailCreateInput(d, testGuardrailEncryptionTypeDefinition())
	if err != nil {
		t.Fatal(err)
	}

	got := settingsByType(t, input)["acceptStrongerEncryption"]
	want := map[string]interface{}{"type": "acceptStrongerEncryption", "template": "false", "templateInput": "{ resource { turbot { id } } }"}
	if !reflect.DeepEqual(got, want) {
		t.Errorf("calculated setting: got %#v, want %#v", got, want)
	}
	if !reflect.DeepEqual(input["action"], map[string]interface{}{"template": "x", "templateInput": "{ resource { turbot { id } } }"}) {
		t.Errorf("calculated action: %#v", input["action"])
	}
}

func TestGuardrailCreateInputFixedActionIsAValue(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTurbotGuardrail().Schema, testGuardrailRaw(map[string]interface{}{
		"action": []interface{}{map[string]interface{}{"value": "tmod:@turbot/aws-s3#/action/types/bucketUpdateEncryptionAtRest"}},
	}))

	input, err := guardrailCreateInput(d, testGuardrailEncryptionTypeDefinition())
	if err != nil {
		t.Fatal(err)
	}
	if !reflect.DeepEqual(input["action"], map[string]interface{}{"value": "tmod:@turbot/aws-s3#/action/types/bucketUpdateEncryptionAtRest"}) {
		t.Errorf("action: %#v", input["action"])
	}
}

func TestGuardrailCreateInputRejectsUnknownSetting(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTurbotGuardrail().Schema, testGuardrailRaw(map[string]interface{}{
		"setting": []interface{}{map[string]interface{}{"type": "EncryptionLevel", "value": "None"}},
	}))

	_, err := guardrailCreateInput(d, testGuardrailEncryptionTypeDefinition())

	if err == nil || !strings.Contains(err.Error(), "acceptStrongerEncryption, customerManagedKey, encryptionLevel") {
		t.Fatalf("expected an unknown-setting error listing the valid names, got %v", err)
	}
}

func TestGuardrailCreateInputRejectsValueAndTemplate(t *testing.T) {
	d := schema.TestResourceDataRaw(t, resourceTurbotGuardrail().Schema, testGuardrailRaw(map[string]interface{}{
		"setting": []interface{}{map[string]interface{}{"type": "encryptionLevel", "value": "None", "template": "x", "template_input": "{ a }"}},
	}))

	if _, err := guardrailCreateInput(d, testGuardrailEncryptionTypeDefinition()); err == nil {
		t.Fatal("a setting with both a value and a template must be rejected")
	}
}

func TestGuardrailUpdateInputSendsOnlyChangedFields(t *testing.T) {
	setting := []interface{}{map[string]interface{}{"type": "encryptionLevel", "value": "S3 managed key"}}
	d := testGuardrailUpdateData(t,
		testGuardrailRaw(map[string]interface{}{"setting": setting, "grace_period": "14d"}),
		testGuardrailRaw(map[string]interface{}{"setting": setting, "grace_period": "14d", "title": "tf-acc-unit-renamed"}))

	input, err := guardrailUpdateInput(d, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"id": "401101108545745", "title": "tf-acc-unit-renamed"}
	if !reflect.DeepEqual(input, want) {
		t.Fatalf("got %#v, want %#v", input, want)
	}
}

func TestGuardrailUpdateInputRemovedTagIsNull(t *testing.T) {
	d := testGuardrailUpdateData(t,
		testGuardrailRaw(map[string]interface{}{"tags": map[string]interface{}{"keep": "a", "drop": "b"}}),
		testGuardrailRaw(map[string]interface{}{"tags": map[string]interface{}{"keep": "a"}}))

	input, err := guardrailUpdateInput(d, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"keep": "a", "drop": nil}
	if !reflect.DeepEqual(input["tags"], want) {
		t.Fatalf("tags: got %#v, want %#v", input["tags"], want)
	}
}

func TestGuardrailUpdateInputModeGracePeriodAndAkas(t *testing.T) {
	d := testGuardrailUpdateData(t,
		testGuardrailRaw(map[string]interface{}{"grace_period": "14d", "akas": []interface{}{"tf-acc-a"}}),
		testGuardrailRaw(map[string]interface{}{"mode": "Draft", "akas": []interface{}{"tf-acc-b"}}))

	input, err := guardrailUpdateInput(d, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{
		"id":          "401101108545745",
		"mode":        map[string]interface{}{"value": "Draft"},
		"gracePeriod": nil,
		"akas":        []interface{}{"tf-acc-b"},
	}
	if !reflect.DeepEqual(input, want) {
		t.Fatalf("got %#v, want %#v", input, want)
	}
}

func TestGuardrailUpdateInputRemovedSettingBlockWritesDefault(t *testing.T) {
	d := testGuardrailUpdateData(t,
		testGuardrailRaw(map[string]interface{}{"setting": []interface{}{
			map[string]interface{}{"type": "encryptionLevel", "value": "S3 managed key"},
			map[string]interface{}{"type": "acceptStrongerEncryption", "value": "false"},
			map[string]interface{}{"type": "customerManagedKey", "value": "alias/x"},
		}}),
		testGuardrailRaw(map[string]interface{}{"setting": []interface{}{
			map[string]interface{}{"type": "encryptionLevel", "value": "None"},
		}}))

	input, err := guardrailUpdateInput(d, testGuardrailEncryptionTypeDefinition())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]map[string]interface{}{
		"encryptionLevel":          {"type": "encryptionLevel", "valueSource": "None"},
		"acceptStrongerEncryption": {"type": "acceptStrongerEncryption", "value": true},
		"customerManagedKey": {"type": "customerManagedKey", "template": "{{ $.defaultKey }}",
			"templateInput": "{\n  defaultKey: policy(uri: \"aws-kms#/policy/types/keyDefaultCustomerManagedKey\")\n}\n"},
	}
	if got := settingsByType(t, input); !reflect.DeepEqual(got, want) {
		t.Fatalf("settings:\n got %#v\nwant %#v", got, want)
	}
}

func TestGuardrailUpdateInputUnchangedSettingIsNotSent(t *testing.T) {
	d := testGuardrailUpdateData(t,
		testGuardrailRaw(map[string]interface{}{"setting": []interface{}{
			map[string]interface{}{"type": "encryptionLevel", "value": "S3 managed key"},
			map[string]interface{}{"type": "acceptStrongerEncryption", "value": "false"},
		}}),
		testGuardrailRaw(map[string]interface{}{"setting": []interface{}{
			map[string]interface{}{"type": "encryptionLevel", "value": "S3 managed key"},
			map[string]interface{}{"type": "acceptStrongerEncryption", "value": "true"},
		}}))

	input, err := guardrailUpdateInput(d, testGuardrailEncryptionTypeDefinition())
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]map[string]interface{}{"acceptStrongerEncryption": {"type": "acceptStrongerEncryption", "valueSource": "true"}}
	if got := settingsByType(t, input); !reflect.DeepEqual(got, want) {
		t.Fatalf("settings:\n got %#v\nwant %#v", got, want)
	}
}

func TestGuardrailUpdateInputResourceFilterRemovalClears(t *testing.T) {
	newRaw := testGuardrailRaw(nil)
	delete(newRaw, "resource_filter")
	d := testGuardrailUpdateData(t, testGuardrailRaw(nil), newRaw)

	input, err := guardrailUpdateInput(d, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"id": "401101108545745", "resourceFilter": ""}
	if !reflect.DeepEqual(input, want) {
		t.Fatalf("got %#v, want %#v", input, want)
	}
}

func TestGuardrailUpdateInputSwitchesToTypedFilters(t *testing.T) {
	newRaw := testGuardrailRaw(map[string]interface{}{"account_filter_id": "401101111592184", "resource_filter_id": "401101112555773"})
	delete(newRaw, "account_filter")
	delete(newRaw, "resource_filter")
	d := testGuardrailUpdateData(t, testGuardrailRaw(nil), newRaw)

	input, err := guardrailUpdateInput(d, nil)
	if err != nil {
		t.Fatal(err)
	}
	want := map[string]interface{}{"id": "401101108545745", "accountFilterId": "401101111592184", "resourceFilterId": "401101112555773"}
	if !reflect.DeepEqual(input, want) {
		t.Fatalf("got %#v, want %#v", input, want)
	}
}

func TestGuardrailFilterAttributesConflict(t *testing.T) {
	for _, pair := range [][2]string{{"account_filter", "account_filter_id"}, {"resource_filter", "resource_filter_id"}} {
		raw := testGuardrailRaw(map[string]interface{}{pair[0]: "x", pair[1]: "401101111592184"})
		_, errs := resourceTurbotGuardrail().Validate(terraform.NewResourceConfigRaw(raw))
		if len(errs) == 0 {
			t.Errorf("%s and %s together must be rejected", pair[0], pair[1])
		}
	}
	if _, errs := resourceTurbotGuardrail().Validate(terraform.NewResourceConfigRaw(testGuardrailRaw(nil))); len(errs) != 0 {
		t.Errorf("a filter of each kind is valid: %v", errs)
	}
}

func testGuardrailPlan(t *testing.T, oldRaw, newRaw map[string]interface{}) error {
	t.Helper()
	r := resourceTurbotGuardrail()
	old := schema.TestResourceDataRaw(t, r.Schema, oldRaw)
	old.SetId("401101108545745")
	_, err := r.Diff(old.State(), terraform.NewResourceConfigRaw(newRaw), nil)
	return err
}

func TestGuardrailPlanRejectsAccountFilterRemoval(t *testing.T) {
	newRaw := testGuardrailRaw(nil)
	delete(newRaw, "account_filter")

	err := testGuardrailPlan(t, testGuardrailRaw(nil), newRaw)

	if err == nil || !strings.Contains(err.Error(), "account_filter cannot be removed") {
		t.Fatalf("expected a plan error naming account_filter, got %v", err)
	}
}

func TestGuardrailPlanRejectsAccountFilterIdRemoval(t *testing.T) {
	oldRaw := testGuardrailRaw(map[string]interface{}{"account_filter_id": "401101111592184"})
	delete(oldRaw, "account_filter")
	newRaw := testGuardrailRaw(nil)
	delete(newRaw, "account_filter")

	err := testGuardrailPlan(t, oldRaw, newRaw)

	if err == nil || !strings.Contains(err.Error(), "account_filter_id cannot be removed") {
		t.Fatalf("expected a plan error naming account_filter_id, got %v", err)
	}
}

func TestGuardrailPlanAllowsAccountFilterSwitch(t *testing.T) {
	newRaw := testGuardrailRaw(map[string]interface{}{"account_filter_id": "401101111592184"})
	delete(newRaw, "account_filter")

	if err := testGuardrailPlan(t, testGuardrailRaw(nil), newRaw); err != nil {
		t.Fatalf("switching account_filter to account_filter_id is allowed: %s", err)
	}
}

func TestGuardrailPlanRejectsActionRemoval(t *testing.T) {
	oldRaw := testGuardrailRaw(map[string]interface{}{"action": []interface{}{map[string]interface{}{"value": "tmod:@turbot/aws-s3#/action/types/bucketUpdateEncryptionAtRest"}}})

	err := testGuardrailPlan(t, oldRaw, testGuardrailRaw(nil))

	if err == nil || !strings.Contains(err.Error(), "action block cannot be removed") {
		t.Fatalf("expected a plan error for removing the action, got %v", err)
	}
}

const testGuardrailReadResponse = `{"data":{"guardrail":{
	"title":"tf-acc-unit","description":null,"gracePeriod":"14d",
	"accountFilter":"title:tf-acc-no-such-account","accountFilterId":null,"resourceFilter":null,"resourceFilterId":null,
	"type":{"uri":"tmod:@turbot/aws-s3#/guardrail/types/bucketEncryptionAtRestGuardrail","policyTypes":{"items":[
		{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailAcceptStrongerEncryption","schema":{"type":"boolean","default":true},"defaultTemplate":null,"defaultTemplateInput":null},
		{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailCustomerManagedKey","schema":{},"defaultTemplate":"{{ $.defaultKey }}","defaultTemplateInput":"{\n  defaultKey: policy(uri: \"aws-kms#/policy/types/keyDefaultCustomerManagedKey\")\n}\n"},
		{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailEncryptionLevel","schema":{"type":"string","default":"AWS managed key"},"defaultTemplate":null,"defaultTemplateInput":null}]}},
	"settings":{"items":[
		{"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailAcceptStrongerEncryption"},"value":true,"valueSource":"true\n","isCalculated":false,"template":null,"templateInput":null},
		{"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailCustomerManagedKey"},"value":null,"valueSource":null,"isCalculated":true,"template":"{{ $.defaultKey }}","templateInput":"{\n  defaultKey: policy(uri: \"aws-kms#/policy/types/keyDefaultCustomerManagedKey\")\n}\n"},
		{"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailEncryptionLevel"},"value":"S3 managed key","valueSource":"S3 managed key\n","isCalculated":false,"template":null,"templateInput":null},
		{"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailMode"},"value":"Check","valueSource":"Check\n","isCalculated":false,"template":null,"templateInput":null},
		{"type":{"uri":"tmod:@turbot/aws-s3#/policy/types/bucketEncryptionAtRestGuardrailAction"},"value":"tmod:@turbot/aws-s3#/action/types/bucketUpdateEncryptionAtRest","valueSource":"tmod:@turbot/aws-s3#/action/types/bucketUpdateEncryptionAtRest\n","isCalculated":false,"template":null,"templateInput":null}]},
	"turbot":{"id":"401101108545745","akas":null,"tags":{"keep":"a"}}}}}`

func testGuardrailData(t *testing.T, raw map[string]interface{}) *schema.ResourceData {
	t.Helper()
	d := schema.TestResourceDataRaw(t, resourceTurbotGuardrail().Schema, raw)
	d.SetId("401101108545745")
	return d
}

func TestGuardrailReadOmitsDefaultSettingsWithNoBlock(t *testing.T) {
	client := typedFilterTestClient(t, testGuardrailReadResponse)
	d := testGuardrailData(t, map[string]interface{}{})

	if err := resourceTurbotGuardrailRead(d, client); err != nil {
		t.Fatal(err)
	}

	settings := guardrailSettingBlocks(d.Get("setting"))
	if len(settings) != 1 || settings["encryptionLevel"]["value"] != "S3 managed key" {
		t.Fatalf("only the non-default setting belongs in state, got %#v", settings)
	}
	if d.Get("mode") != "Check" || d.Get("grace_period") != "14d" || d.Get("tags.keep") != "a" {
		t.Errorf("mode %v, grace_period %v, tags %v", d.Get("mode"), d.Get("grace_period"), d.Get("tags"))
	}
	if d.Get("action.0.value") != "tmod:@turbot/aws-s3#/action/types/bucketUpdateEncryptionAtRest" {
		t.Errorf("action: %#v", d.Get("action"))
	}
}

func TestGuardrailReadKeepsDefaultSettingsWithABlock(t *testing.T) {
	client := typedFilterTestClient(t, testGuardrailReadResponse)
	d := testGuardrailData(t, map[string]interface{}{"setting": []interface{}{
		map[string]interface{}{"type": "acceptStrongerEncryption", "value": "true"},
		map[string]interface{}{"type": "customerManagedKey", "template": "{{ $.defaultKey }}", "template_input": "{\n  defaultKey: policy(uri: \"aws-kms#/policy/types/keyDefaultCustomerManagedKey\")\n}\n"},
	}})

	if err := resourceTurbotGuardrailRead(d, client); err != nil {
		t.Fatal(err)
	}

	settings := guardrailSettingBlocks(d.Get("setting"))
	if len(settings) != 3 {
		t.Fatalf("settings with a block stay in state, got %#v", settings)
	}
	if settings["acceptStrongerEncryption"]["value"] != "true" || settings["customerManagedKey"]["template"] != "{{ $.defaultKey }}" {
		t.Errorf("settings: %#v", settings)
	}
}

func TestGuardrailReadKeepsStateOnOtherError(t *testing.T) {
	client := typedFilterTestClient(t, `{"errors":[{"message":"Forbidden: Insufficient permissions"}]}`)
	d := testGuardrailData(t, testGuardrailRaw(nil))

	err := resourceTurbotGuardrailRead(d, client)

	if err == nil {
		t.Fatal("expected the read error to fail the refresh")
	}
	if d.Id() != "401101108545745" {
		t.Fatalf("a non-not-found error must keep the guardrail in state, id is %q", d.Id())
	}
}

func TestGuardrailReadDropsStateWhenGone(t *testing.T) {
	for name, body := range map[string]string{
		"null":            `{"data":{"guardrail":null}}`,
		"not a guardrail": `{"data":{"guardrail":{"title":null,"type":null,"settings":{"items":[]},"turbot":{"id":"401101108545745"}}}}`,
	} {
		t.Run(name, func(t *testing.T) {
			client := typedFilterTestClient(t, body)
			d := testGuardrailData(t, testGuardrailRaw(nil))

			if err := resourceTurbotGuardrailRead(d, client); err != nil {
				t.Fatalf("a missing guardrail is not a refresh error: %s", err)
			}
			if d.Id() != "" {
				t.Fatalf("a deleted guardrail must drop from state, id is %q", d.Id())
			}
		})
	}
}

func TestGuardrailPlanRejectsARepeatedSetting(t *testing.T) {
	raw := testGuardrailRaw(map[string]interface{}{"setting": []interface{}{
		map[string]interface{}{"type": "encryptionLevel", "value": "None"},
		map[string]interface{}{"type": "encryptionLevel", "value": "S3 managed key"},
	}})

	err := testGuardrailPlan(t, testGuardrailRaw(nil), raw)

	if err == nil || !strings.Contains(err.Error(), "more than once") {
		t.Fatalf("expected a repeated setting to be rejected, got %v", err)
	}
}
