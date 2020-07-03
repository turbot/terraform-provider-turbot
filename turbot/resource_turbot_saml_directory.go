package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"github.com/terraform-providers/terraform-provider-turbot/helpers"
	"strings"
)

// input properties which must be passed to a create/update call
var samlDirectoryInputProperties = []interface{}{"title", "description", "status", "entry_point", "issuer", "certificate", "profile_id_template", "name_id_format", "sign_requests", "allow_group_syncing", "profile_groups_attribute", "group_filter", "signature_private_key", "signature_algorithm", "pool_id", "parent", "tags"}

// exclude properties from input map to make a create call
func getSamlDirectoryProperties() []interface{} {
	excludedProperties := []string{"group_id_template", "tags", "profile_id_template"}
	return helpers.RemoveProperties(localDirectoryInputProperties, excludedProperties)
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
				Required: true,
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
			"sign_requests": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signature_private_key": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"signature_algorithm": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"allow_group_syncing": {
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
				Deprecated: "use '' argument instead",
			},
			"pool_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "use '' argument instead",
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
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
	return nil
}

func resourceTurbotSamlDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	samlDirectory, err := client.ReadSamlDirectory(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
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
	d.Set("sign_requests", samlDirectory.SignRequests)
	d.Set("signature_private_key", samlDirectory.SignaturePrivateKey)
	d.Set("signature_algorithm", samlDirectory.SignatureAlgorithm)
	d.Set("allow_group_syncing", samlDirectory.AllowGroupSyncing)
	d.Set("profile_groups_attribute", samlDirectory.ProfileGroupsAttribute)
	d.Set("group_filter", samlDirectory.GroupFilter)
	return nil
}

func resourceTurbotSamlDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	input := mapFromResourceData(d, getSamlDirectoryProperties())
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
