package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/turbot/terraform-provider-turbot/apiClient"
	"github.com/turbot/terraform-provider-turbot/errors"
	"github.com/turbot/terraform-provider-turbot/helpers"
	"strings"
)

// these are the properties which must be passed to a create/update call
// each element in the array is either a map, defining an explicit mapping, or a string, which is the terraform property name
// this is automatically mapped to the turbot property name by converting snake -> camel case

var googleDirectoryInputProperties = []interface{}{
	// explicit mapping
	map[string]string{"hosted_name": "hostedDomain"},
	"title", "parent", "client_id", "description", "client_secret", "profile_id_template", "tags"}

// exclude properties from input map to make a update call
func getGoogleDirectoryUpdateProperties() []interface{} {
	excludedProperties := []string{"profile_id_template", "tags"}
	return helpers.RemoveProperties(googleDirectoryInputProperties, excludedProperties)
}

func resourceGoogleDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotGoogleDirectoryCreate,
		Read:   resourceTurbotGoogleDirectoryRead,
		Update: resourceTurbotGoogleDirectoryUpdate,
		Delete: resourceTurbotGoogleDirectoryDelete,
		Exists: resourceTurbotGoogleDirectoryExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotGoogleDirectoryImport,
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Required: true,
				// when doing a diff, the state file will contain the id of the parent but the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
			},
			// when doing a read, fetch the parent akas to use in suppressIfAkaMatches
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
			"profile_id_template": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:             schema.TypeString,
				Required:         true,
				DiffSuppressFunc: suppressIfClientSecretPresent,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"pgp_key": {
				Type:     schema.TypeString,
				ForceNew: true,
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
			"key_fingerprint": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"pool_id": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "pool_id has been deprecated and is not used.",
			},
			"hosted_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id_template": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "group_id_template has been deprecated and is not used",
			},
			"login_name_template": {
				Type:       schema.TypeString,
				Optional:   true,
				Deprecated: "login_name_template has been deprecated and is not used",
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceTurbotGoogleDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotGoogleDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	// build mutation input
	input := mapFromResourceData(d, googleDirectoryInputProperties)
	input["status"] = "ACTIVE"
	turbotMetadata, err := client.CreateGoogleDirectory(input)
	if err != nil {
		return err
	}
	// assign computed properties
	if err = storeClientSecret(d, input["clientSecret"].(string)); err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)
	d.Set("status", input["status"])
	d.Set("directory_type", input["directoryType"])
	return nil
}

func resourceTurbotGoogleDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	googleDirectory, err := client.ReadGoogleDirectory(id)
	if err != nil {
		if errors.NotFoundError(err) {
			// directory was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("parent", googleDirectory.Parent)
	d.Set("title", googleDirectory.Title)
	d.Set("directory_type", googleDirectory.DirectoryType)
	d.Set("status", strings.ToUpper(googleDirectory.Status))
	d.Set("profile_id_template", googleDirectory.ProfileIdTemplate)
	d.Set("description", googleDirectory.Description)
	d.Set("client_id", googleDirectory.ClientID)
	d.Set("pool_id", googleDirectory.PoolId)
	d.Set("group_id_template", googleDirectory.GroupIdTemplate)
	d.Set("login_name_template", googleDirectory.LoginNameTemplate)
	d.Set("hosted_name", googleDirectory.HostedDomain)
	d.Set("tags", googleDirectory.Turbot.Tags)
	// set parent_akas property by loading parent resource and fetching the akas
	return storeAkas(googleDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotGoogleDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build mutation payload
	input := mapFromResourceData(d, getGoogleDirectoryUpdateProperties())
	input["id"] = d.Id()
	// do update
	turbotMetadata, err := client.UpdateGoogleDirectory(input)
	if err != nil {
		return err
	}
	clientSecret := input["clientSecret"].(string)
	// set parent_akas property by loading parent resource and fetching the akas
	if err := storeAkas(turbotMetadata.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// store client secret, encrypting if a pgp key was provided
	return storeClientSecret(d, clientSecret)
}

func resourceTurbotGoogleDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotGoogleDirectoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotGoogleDirectoryRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// write client secret to ResourceData, encrypting if a pgp key was provided
func storeClientSecret(d *schema.ResourceData, clientSecret string) error {
	if pgpKey, ok := d.GetOk("pgp_key"); ok {
		fingerprint, encrypted, err := helpers.EncryptValue(pgpKey.(string), clientSecret)
		if err != nil {
			return err
		}
		d.Set("client_secret", encrypted)
		d.Set("key_fingerprint", fingerprint)
	} else {
		d.Set("client_secret", clientSecret)
	}
	return nil
}

func suppressIfClientSecretPresent(k, old, new string, d *schema.ResourceData) bool {
	// We do not read back client secret so suppress diff caused by empty value
	_, keyPresent := d.GetOk("pgp_key")
	if old != "" {
		if keyPresent || new == "" {
			return true
		}
	}
	return false
}
