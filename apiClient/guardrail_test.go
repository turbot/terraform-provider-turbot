package apiClient

import (
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadGuardrailNullIsNotFound(t *testing.T) {
	client := typedFilterClient(t, func(string) string { return `{"data":{"guardrail":null}}` })

	_, err := client.ReadGuardrail("1")

	assert.True(t, IsGuardrailNotFound(err))
}

func TestReadGuardrailOfAnotherResourceIsNotFound(t *testing.T) {
	client := typedFilterClient(t, func(string) string {
		return `{"data":{"guardrail":{"title":null,"type":null,"settings":{"items":[]},"turbot":{"id":"1"}}}}`
	})

	_, err := client.ReadGuardrail("1")

	assert.True(t, IsGuardrailNotFound(err))
}

func TestReadGuardrailOtherErrorIsNotNotFound(t *testing.T) {
	client := typedFilterClient(t, func(string) string {
		return `{"errors":[{"message":"Forbidden: Insufficient permissions"}]}`
	})

	_, err := client.ReadGuardrail("1")

	assert.Error(t, err)
	assert.False(t, IsGuardrailNotFound(err))
}

func TestReadGuardrailReturnsSettings(t *testing.T) {
	client := typedFilterClient(t, func(string) string {
		return `{"data":{"guardrail":{"title":"t","gracePeriod":null,"accountFilter":"title:x","type":{"uri":"tmod:@turbot/aws-s3#/guardrail/types/bucketVersioningGuardrail","mode":{"uri":"m"},"action":null,"policyTypes":{"items":[{"uri":"v","schema":{"default":"Enabled"},"defaultTemplate":null,"defaultTemplateInput":null}]}},"settings":{"items":[{"type":{"uri":"m"},"value":"Check","valueSource":"Check\n","isCalculated":false,"template":null,"templateInput":null}]},"turbot":{"id":"1","akas":null,"tags":{}}}}}`
	})

	guardrail, err := client.ReadGuardrail("1")

	assert.NoError(t, err)
	assert.Equal(t, "title:x", guardrail.AccountFilter)
	assert.Equal(t, "", guardrail.GracePeriod)
	assert.Equal(t, "Enabled", guardrail.Type.PolicyTypes.Items[0].Schema["default"])
	assert.Equal(t, "Check", guardrail.Settings.Items[0].Value)
}

func TestReadGuardrailTypeNullIsAnError(t *testing.T) {
	client := typedFilterClient(t, func(string) string { return `{"data":{"guardrailType":null}}` })

	_, err := client.ReadGuardrailType("tmod:@turbot/x#/guardrail/types/missing")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "not found")
}

func TestDeleteGuardrailOfMissingGuardrailSucceeds(t *testing.T) {
	client := typedFilterClient(t, func(body string) string {
		if strings.Contains(body, "deleteGuardrail") {
			return `{"data":{"deleteGuardrail":null},"errors":[{"message":"Bad Request: Failed to get resourceId from aka 1: Not Found"}]}`
		}
		return `{"data":{"guardrail":null}}`
	})

	assert.NoError(t, client.DeleteGuardrail("1"))
}

func TestDeleteGuardrailRefusalSurfaces(t *testing.T) {
	client := typedFilterClient(t, func(body string) string {
		if strings.Contains(body, "deleteGuardrail") {
			return `{"data":{"deleteGuardrail":null},"errors":[{"message":"Forbidden: Insufficient permissions"}]}`
		}
		return `{"data":{"guardrail":{"title":"t","type":{"uri":"g"},"settings":{"items":[]},"turbot":{"id":"1"}}}}`
	})

	err := client.DeleteGuardrail("1")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "Insufficient permissions")
}
