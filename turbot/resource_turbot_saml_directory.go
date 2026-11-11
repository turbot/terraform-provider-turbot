package turbot

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
)

// input properties which must be passed to a create/update call
var samlDirectoryInputProperties = []interface{}{"title", "description", "status", "entry_point", "issuer", "certificate", "profile_id_template", "name_id_format", "sign_requests", "allow_idp_initiated_sso", "allow_group_syncing", "profile_groups_attribute", "group_filter", "signature_private_key", "signature_algorithm", "pool_id", "parent", "tags"}

// exclude properties from input map to make a create call
func getSamlDirectoryProperties() []interface{} {
	excludedProperties := []string{"group_id_template", "tags", "profile_id_template"}
	return helpers.RemoveProperties(samlDirectoryInputProperties, excludedProperties)
}

func resourceTurbotSamlDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotSamlDirectoryCreate,
		Read:   resourceTurbotSamlDirectoryRead,
		Update: resourceTurbotSamlDirectoryUpdate,
		Delete: resourceTurbotSamlDirectoryDelete,
		Exists: resourceTurbotSamlDirectoryExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotSamlDirectoryImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resourcesamlDirectoryProperties
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent but he config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
			},
			// when doing a read, fetch the parent akas to use in suppressIfAkaMatches("parent_akas")()
			"parent_akas": {
				Type:     schema.TypeList,
				Computed: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"title": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"directory_type": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"entry_point": {
				Type:     schema.TypeString,
				Required: true,
			},
			"certificate": {
				Type:     schema.TypeString,
				Required: true,
			},
			"profile_id_template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"issuer": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"name_id_format": {
				Type:     schema.TypeString,
				Optional: true,
				Default:  "UNSPECIFIED",
			},
			// Guardrails takes and returns a Boolean, but sign_requests stays a string so that existing state keeps
			// its type. Computed, as is signature_algorithm, because the console saves both on every edit.
			"sign_requests": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateSignRequests,
			},
			// Guardrails stores the key as a secret and returns only a reference to it, so it is never read back.
			"signature_private_key": {
				Type:      schema.TypeString,
				Optional:  true,
				Sensitive: true,
			},
			"signature_algorithm": {
				Type:         schema.TypeString,
				Optional:     true,
				Computed:     true,
				ValidateFunc: validateSignatureAlgorithm,
			},
			"allow_group_syncing": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"allow_idp_initiated_sso": {
				Type:     schema.TypeBool,
				Optional: true,
				Default:  false,
			},
			"profile_groups_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id_template": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "group_id_template has been deprecated and is not used.",
			},
			"pool_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "pool_id has been deprecated and is not used.",
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceTurbotSamlDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotSamlDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	input := mapFromResourceData(d, samlDirectoryInputProperties)
	setSignRequestsInput(input)
	// set computed properties
	input["status"] = "ACTIVE"
	samlDirectory, err := client.CreateSamlDirectory(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(samlDirectory.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(samlDirectory.Turbot.Id)
	// assign Read query properties
	d.Set("status", strings.ToUpper(samlDirectory.Status))
	d.Set("parent", samlDirectory.Parent)
	d.Set("title", samlDirectory.Title)
	d.Set("description", samlDirectory.Description)
	setSigning(d, samlDirectory)
	return nil
}

func resourceTurbotSamlDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	samlDirectory, err := client.ReadSamlDirectory(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// saml directory was not found - clear id
			d.SetId("")
		}
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(samlDirectory.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign results back into ResourceData
	d.Set("parent", samlDirectory.Parent)
	d.Set("title", samlDirectory.Title)
	d.Set("description", samlDirectory.Description)
	d.Set("status", strings.ToUpper(samlDirectory.Status))
	d.Set("name_id_format", strings.ToUpper(samlDirectory.NameIdFormat))
	d.Set("profile_id_template", samlDirectory.ProfileIdTemplate)
	d.Set("entry_point", samlDirectory.EntryPoint)
	d.Set("certificate", samlDirectory.Certificate)
	setSigning(d, samlDirectory)
	d.Set("allow_group_syncing", samlDirectory.AllowGroupSyncing)
	d.Set("allow_idp_initiated_sso", samlDirectory.AllowIdpInitiatedSso)
	d.Set("profile_groups_attribute", samlDirectory.ProfileGroupsAttribute)
	d.Set("group_filter", samlDirectory.GroupFilter)
	d.Set("tags", samlDirectory.Turbot.Tags)
	return nil
}

func resourceTurbotSamlDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	input := mapFromResourceData(d, getSamlDirectoryProperties())
	setSignRequestsInput(input)
	input["id"] = d.Id()

	// update saml directory returns saml directory
	samlDirectory, err := client.UpdateSamlDirectory(input)
	if err != nil {
		return err
	}

	// assign Read query properties
	d.Set("parent", samlDirectory.Parent)
	d.Set("title", samlDirectory.Title)
	d.Set("description", samlDirectory.Description)
	setSigning(d, samlDirectory)
	// set parent_akas property by loading parent resource and fetching the akas
	return storeAkas(samlDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotSamlDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotSamlDirectoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotSamlDirectoryRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// setSignRequestsInput converts sign_requests to the Boolean that Guardrails takes.
func setSignRequestsInput(input map[string]interface{}) {
	if v, ok := input["signRequests"]; ok {
		input["signRequests"] = v.(string) == "true"
	}
}

// setSigning stores the signing settings that Guardrails returns. Guardrails stores signatureAlgorithm as "SHA-256",
// but the create and update inputs take "SHA256".
func setSigning(d *schema.ResourceData, samlDirectory *apiClient.SamlDirectory) {
	signRequests := ""
	if samlDirectory.SignRequests != nil {
		signRequests = strconv.FormatBool(*samlDirectory.SignRequests)
	}
	d.Set("sign_requests", signRequests)
	d.Set("signature_algorithm", strings.ReplaceAll(samlDirectory.SignatureAlgorithm, "-", ""))
}

func validateSignRequests(v interface{}, k string) ([]string, []error) {
	if s := v.(string); s != "true" && s != "false" {
		return nil, []error{fmt.Errorf("%q must be true or false, got %q", k, s)}
	}
	return nil, nil
}

func validateSignatureAlgorithm(v interface{}, k string) ([]string, []error) {
	switch s := v.(string); s {
	case "SHA1", "SHA256", "SHA512":
		return nil, nil
	default:
		return nil, []error{fmt.Errorf("%q must be one of SHA1, SHA256 or SHA512, got %q", k, s)}
	}
}
