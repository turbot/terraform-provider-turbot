package turbot

import (
	"fmt"
	"log"
	"time"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/httpclient"
	"github.com/hashicorp/terraform/terraform"
	"github.com/turbot/terraform-provider-turbot/apiClient"
)

// ProviderVersion is set from the release build's version.
var ProviderVersion = "dev"

func Provider() terraform.ResourceProvider {
	p := &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"secret_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"workspace": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"profile": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"credentials_file": {
				Type:     schema.TypeString,
				Optional: true,
			},
			// request_timeout bounds a single GraphQL request. Accepts a Go duration string
			// ("30s", "10m"). Empty falls back to apiClient.DefaultRequestTimeout. Guardrails
			// operations can be genuinely slow, so this exists to stop an indefinitely-hung
			// connection, not to enforce a tight latency budget.
			"request_timeout": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"turbot_control_mute":            resourceTurbotControlMute(),
			"turbot_file":                    resourceTurbotFile(),
			"turbot_folder":                  resourceTurbotFolder(),
			"turbot_google_directory":        resourceGoogleDirectory(),
			"turbot_grant":                   resourceTurbotGrant(),
			"turbot_grant_activation":        resourceTurbotGrantActivation(),
			"turbot_guardrail":               resourceTurbotGuardrail(),
			"turbot_ldap_directory":          resourceTurbotLdapDirectory(),
			"turbot_local_directory":         resourceTurbotLocalDirectory(),
			"turbot_local_directory_user":    resourceTurbotLocalDirectoryUser(),
			"turbot_mod":                     resourceTurbotMod(),
			"turbot_policy_pack":             resourceTurbotPolicyPack(),
			"turbot_policy_pack_attachment":  resourceTurbotPolicyPackAttachment(),
			"turbot_policy_setting":          resourceTurbotPolicySetting(),
			"turbot_profile":                 resourceTurbotProfile(),
			"turbot_resource":                resourceTurbotResource(),
			"turbot_saml_directory":          resourceTurbotSamlDirectory(),
			"turbot_shadow_resource":         resourceTurbotShadowResource(),
			"turbot_smart_folder":            resourceTurbotSmartFolder(),
			"turbot_smart_folder_attachment": resourceTurbotSmartFolderAttachemnt(),
			"turbot_turbot_directory":        resourceTurbotTurbotDirectory(),
			"turbot_typed_filter":            resourceTurbotTypedFilter(),
			"turbot_watch":                   resourceTurbotWatch(),
			//"turbot_group_profile":           resourceTurbotGroupProfile(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"turbot_control":      dataSourceTurbotControl(),
			"turbot_policy_value": dataSourceTurbotPolicyValue(),
			"turbot_resource":     dataSourceTurbotResource(),
		},
	}
	p.ConfigureFunc = func(d *schema.ResourceData) (interface{}, error) {
		return providerConfigure(d, userAgent(p.TerraformVersion))
	}
	return p
}

// userAgent names Terraform and this provider, so Guardrails can tell a Terraform change from
// any other client's.
func userAgent(terraformVersion string) string {
	return fmt.Sprintf("%s terraform-provider-turbot/%s", httpclient.TerraformUserAgent(terraformVersion), ProviderVersion)
}

func providerConfigure(d *schema.ResourceData, userAgent string) (interface{}, error) {
	config := apiClient.ClientConfig{
		UserAgent: userAgent,
		Credentials: apiClient.ClientCredentials{
			AccessKey: d.Get("access_key").(string),
			SecretKey: d.Get("secret_key").(string),
			Workspace: d.Get("workspace").(string),
		},
		Profile:         d.Get("profile").(string),
		CredentialsPath: d.Get("credentials_file").(string),
	}

	// Parse the optional request_timeout duration. An invalid value is a config error rather than
	// a silent fallback: silently ignoring a mistyped timeout would hide the operator's intent.
	if raw := d.Get("request_timeout").(string); raw != "" {
		timeout, err := time.ParseDuration(raw)
		if err != nil {
			return nil, fmt.Errorf("invalid request_timeout %q: %s", raw, err.Error())
		}
		if timeout < 0 {
			return nil, fmt.Errorf("invalid request_timeout %q: must not be negative", raw)
		}
		config.RequestTimeout = timeout
	}

	client, err := apiClient.CreateClient(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %s", err.Error())
	}
	log.Println("[INFO] Turbot API client initialized, now validating...", client)
	if err = client.Validate(); err != nil {
		return nil, err
	}
	return client, nil
}
