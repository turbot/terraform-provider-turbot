package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/hashicorp/terraform/terraform"
	apiClient "github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
	"net/url"
	"path"
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
			"turbot_policy_setting": resourceTurbotPolicySetting(),
			"turbot_mod":            resourceTurbotMod(),
			"turbot_folder":         resourceTurbotFolder(),
		},
		DataSourcesMap: map[string]*schema.Resource{
			"turbot_policy_value": dataSourceTurbotPolicy(),
		},

		ConfigureFunc: providerConfigure,
	}
}

func providerConfigure(d *schema.ResourceData) (interface{}, error) {
	// build api url
	u, err := url.Parse(d.Get("workspace").(string))
	u.Path = path.Join(u.Path, "api/v5/graphql")
	apiUrl := u.String()

	client := apiClient.CreateClient(
		d.Get("access_key_id").(string),
		d.Get("secret_access_key").(string),
		apiUrl)

	log.Println("[INFO] Turbot API client initialized, now validating...", client)
	err = client.Validate()
	if err != nil {
		return nil, err
	}
	return client, nil
}
