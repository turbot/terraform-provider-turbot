package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	apiClient "github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
)

func Provider() terraform.ResourceProvider {
	return &schema.Provider{
		Schema: map[string]*schema.Schema{
			"access_key_id": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TURBOT_ACCESS_KEY_ID", nil),
			},
			"secret_access_key": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TURBOT_SECRET_ACCESS_KEY", nil),
			},
			"workspace": {
				Type:        schema.TypeString,
				Required:    true,
				DefaultFunc: schema.EnvDefaultFunc("TURBOT_WORKSPACE", nil),
			},
		},

		ResourcesMap: map[string]*schema.Resource{
			"turbot_policy_setting":  resourceTurbotPolicySetting(),
			"turbot_mod":             resourceTurbotMod(),
			"turbot_folder":          resourceTurbotFolder(),
			"turbot_resource":        resourceTurbotResource(),
			"turbot_local_directory": resourceTurbotLocalDirectory(),
			"turbot_profile":         resourceTurbotProfile(),
			"turbot_saml_directory": resourceTurbotSamlDirectory().
		},
		DataSourcesMap: map[string]*schema.Resource{
			"turbot_policy_value": dataSourceTurbotPolicyValue(),
			"turbot_resource":     dataSourceTurbotResource(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	client, err := apiClient.CreateClient(
		d.Get("access_key_id").(string),
		d.Get("secret_access_key").(string),
		d.Get("workspace").(string))
	if err != nil {
		return nil, fmt.Errorf("failed to create client: %s", err.Error())
	}
	log.Println("[INFO] Turbot API client initialized, now validating...", client)
	if err = client.Validate(); err != nil {
		return nil, fmt.Errorf("failed to validate client: %s", err.Error())
	}
	return client, nil
}
