package apiClient

import (
	"fmt"

	errorsHandler "github.com/turbot/terraform-provider-turbot/errors"
)

func (client *Client) CreatePolicySetting(input map[string]interface{}) (*PolicySetting, error) {
	query := createPolicySettingMutation()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{
		"input": input,
	}

	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleCreateError(err, input, "policy setting")
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) ReadPolicySetting(id string) (*PolicySetting, error) {
	query := readPolicySettingQuery()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{"id": id}

	// execute api call
	err := client.doRequest(query, variables, responseData)
	if err == nil {
		return &responseData.PolicySetting, nil
	}
	if !errorsHandler.ForbiddenError(err) {
		return nil, client.handleReadError(err, id, "policy setting")
	}

	// A Forbidden here does not necessarily mean the identity may not read the setting: the
	// primary query requests secretValue/secretValueSource, and Guardrails guards those fields
	// with a Turbot/Admin check at the setting's target resource even when the policy type holds
	// no secret. Retry without the secret fields — Turbot/Metadata or above can read everything
	// else. If the fallback fails too, the denial is on the setting itself (or the workspace
	// predates the fallback's type fields), so surface the original error unchanged.
	fallback := &policySettingWithoutSecretsResponse{}
	if fallbackErr := client.doRequest(readPolicySettingWithoutSecretsQuery(), variables, fallback); fallbackErr != nil {
		return nil, client.handleReadError(err, id, "policy setting")
	}
	return fallback.PolicySetting.toPolicySetting()
}

// policySettingWithoutSecretsResponse decodes the without-secrets read.
//
// Value and ValueSource are interface{} rather than PolicySetting's interface{}/string pair
// because for a secret policy type Guardrails substitutes a secret REFERENCE for both —
// {"secret": {"id": "..."}} — in place of the value (captured live). Decoding that object into
// PolicySetting.ValueSource, a string, fails with an opaque "cannot unmarshal object into Go
// struct field" error before any guard can run, so the shape has to be decoded first and
// inspected second.
type policySettingWithoutSecretsResponse struct {
	PolicySetting policySettingWithoutSecrets
}

type policySettingWithoutSecrets struct {
	Type struct {
		Uri         string
		Secret      bool
		SecretLevel string
	}
	Value              interface{}
	ValueSource        interface{}
	Default            bool
	Precedence         string
	Template           string
	TemplateInput      interface{}
	Input              string
	Note               string
	ValidFromTimestamp string
	ValidToTimestamp   string
	Turbot             TurbotPolicyMetadata
}

// isSecretReference reports whether v is the reference Guardrails substitutes for a secret value
// in the plain value/valueSource fields — {"secret": {"id": "..."}} — rather than the value itself.
//
// This is the load-bearing secret check, because it reads the data actually returned instead of
// trusting the policy type's secret metadata. That metadata cannot carry the check on its own:
// PolicyType.secret and PolicyType.secretLevel are both nullable, and in practice secretLevel is
// null on EVERY policy type on the workspaces checked (5000 of 5000), while `secret: true` marks
// only a handful. A guard keyed on secretLevel would therefore either refuse every setting (if it
// demanded "NONE") or catch nothing (if it demanded "SECRET"/"CONFIDENTIAL"). The type fields are
// still checked below as a second, independent signal.
func isSecretReference(v interface{}) bool {
	m, ok := v.(map[string]interface{})
	if !ok || len(m) != 1 {
		return false
	}
	_, ok = m["secret"]
	return ok
}

// toPolicySetting converts a without-secrets read into a PolicySetting, refusing anything whose
// value is not plainly readable rather than storing a stand-in for it.
func (s policySettingWithoutSecrets) toPolicySetting() (*PolicySetting, error) {
	// Refuse a secret setting: the reference is not the value, so writing it to state would
	// produce a permanent diff and Terraform would try to "correct" a value this identity cannot
	// read. Both the returned shape and the type metadata are checked, so a secret is refused
	// whichever signal is present.
	if isSecretReference(s.Value) || isSecretReference(s.ValueSource) ||
		s.Type.Secret || s.Type.SecretLevel == "SECRET" || s.Type.SecretLevel == "CONFIDENTIAL" {
		return nil, fmt.Errorf("error reading policy setting: policy type %s is secret - reading its value requires Turbot/Admin on resource %s (Turbot/Owner at the Turbot root for SECRET-level policy types)", s.Type.Uri, s.Turbot.ResourceId)
	}

	// valueSource is a YAML string for every non-secret setting. An unrecognised shape is refused
	// rather than silently dropped from state - failing closed on anything not understood.
	valueSource, ok := s.ValueSource.(string)
	if s.ValueSource != nil && !ok {
		return nil, fmt.Errorf("error reading policy setting: unexpected valueSource shape %T for policy type %s on resource %s - refusing to store it", s.ValueSource, s.Type.Uri, s.Turbot.ResourceId)
	}

	setting := &PolicySetting{
		Value:              s.Value,
		ValueSource:        valueSource,
		Default:            s.Default,
		Precedence:         s.Precedence,
		Template:           s.Template,
		TemplateInput:      s.TemplateInput,
		Input:              s.Input,
		Note:               s.Note,
		ValidFromTimestamp: s.ValidFromTimestamp,
		ValidToTimestamp:   s.ValidToTimestamp,
		Turbot:             s.Turbot,
	}
	setting.Type.Uri = s.Type.Uri
	setting.Type.Secret = s.Type.Secret
	setting.Type.SecretLevel = s.Type.SecretLevel
	return setting, nil
}

func (client *Client) UpdatePolicySetting(input map[string]interface{}) (*PolicySetting, error) {
	query := updatePolicySettingMutation()
	responseData := &PolicySettingResponse{}

	variables := map[string]interface{}{
		"input": input,
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return nil, client.handleUpdateError(err, input, "policy setting")
	}
	return &responseData.PolicySetting, nil
}

func (client *Client) DeletePolicySetting(id string) error {
	query := deletePolicySettingMutation()
	responseData := &PolicySettingResponse{}
	variables := map[string]interface{}{
		"input": map[string]string{
			"id": id,
		},
	}
	// execute api call
	if err := client.doRequest(query, variables, responseData); err != nil {
		return fmt.Errorf("error deleting policy: %s", err.Error())
	}
	return nil
}

func (client *Client) FindPolicySetting(policyTypeUri, resourceAka string) (PolicySetting, error) {
	responseData := &FindPolicySettingResponse{}

	query := findPolicySettingQuery()
	variables := map[string]interface{}{
		"filter": []string{fmt.Sprintf("policyType:%s resource:%s", policyTypeUri, resourceAka)},
	}

	// execute api call
	err := client.doRequest(query, variables, responseData)
	if err != nil && errorsHandler.ForbiddenError(err) {
		// The Forbidden may be the per-field secretValue guard firing on a matched item rather
		// than a denial of the list itself — see ReadPolicySetting. Retry without the secret
		// fields; if that fails too, surface the original error.
		//
		// Unlike the read, this path does NOT refuse a secret policy type. The caller
		// (resourceTurbotPolicySettingCreate) uses the result only to detect that a setting
		// already exists, and gates on Value != nil. For a secret type the plain field carries
		// the secret REFERENCE {"secret": {"id": "..."}} rather than the value — non-nil, so an
		// existing secret setting is still detected as existing. Nothing here is stored in state,
		// so the reference is never written anywhere; it only has to be distinguishable from
		// absent, and it is.
		fallbackData := &findPolicySettingWithoutSecretsResponse{}
		if fallbackErr := client.doRequest(findPolicySettingWithoutSecretsQuery(), variables, fallbackData); fallbackErr == nil {
			return fallbackData.firstDefault(), nil
		}
	}
	if err != nil {
		return PolicySetting{}, client.handleReadError(err, policyTypeUri, "policy setting")
	}

	for _, setting := range responseData.PolicySettings.Items {
		if setting.Default {
			return setting, nil
		}
	}
	return PolicySetting{}, nil
}

// findPolicySettingWithoutSecretsResponse mirrors FindPolicySettingResponse with a tolerant
// ValueSource, for the same reason as policySettingWithoutSecretsResponse: a matched item of a
// secret policy type carries an object there, which will not decode into a string.
type findPolicySettingWithoutSecretsResponse struct {
	PolicySettings struct {
		Items []struct {
			Value       interface{}
			ValueSource interface{}
			Default     bool
			Precedence  string
			Turbot      TurbotPolicyMetadata
		}
	}
}

// firstDefault applies the same selection as FindPolicySetting's primary path. Only Value (the
// caller's exists check) and Turbot.Id (its error message) are carried over; ValueSource is
// deliberately dropped rather than coerced, since nothing downstream reads it here.
func (r *findPolicySettingWithoutSecretsResponse) firstDefault() PolicySetting {
	for _, item := range r.PolicySettings.Items {
		if item.Default {
			return PolicySetting{
				Value:      item.Value,
				Precedence: item.Precedence,
				Default:    item.Default,
				Turbot:     item.Turbot,
			}
		}
	}
	return PolicySetting{}
}
