package turbot

import (
"github.com/hashicorp/terraform/helper/schema"
"github.com/terraform-providers/terraform-provider-turbot/apiClient"
"github.com/terraform-providers/terraform-provider-turbot/helpers"
"strings"
)

// input properties which must be passed to a create/update call
var ldapDirectoryInputProperties = []interface{}{
	"title",
	"description",
	"profileIdTemplate",
	"groupProfileIdTemplate",
	"url",
	"distinguishedName",
	"password",
	"base",
	"userObjectFilter",
	"disabledUserFilter",
	"userMatchFilter",
	"userSearchFilter",
	"userSearchAttributes",
	"groupObjectFilter",
	"groupSearchFilter",
	"groupSyncFilter",
	"userCanonicalNameAttribute",
	"userEmailAttribute",
	"userDisplayNameAttribute",
	"userGivenNameAttribute",
	"userFamilyNameAttribute",
	"groupCanonicalNameAttribute",
	"tlsEnabled",
	"tlsServerCertificate",
	"groupMemberOfAttribute",
	"groupMembershipAttribute",
	"connectivityTestFilter",
	"rejectUnauthorized",
	"disabledGroupFilter",
}

// exclude properties from input map to make a update call
func getldapDirectoryUpdateProperties() []interface{} {
	excludedProperties := []string{"profile_id_template"}
	return helpers.RemoveProperties(ldapDirectoryInputProperties, excludedProperties)
}

func resourceTurbotLdapDirectory() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotLdapDirectoryCreate,
		Read:   resourceTurbotLdapDirectoryRead,
		Update: resourceTurbotLdapDirectoryUpdate,
		Delete: resourceTurbotLdapDirectoryDelete,
		Exists: resourceTurbotLdapDirectoryExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotLdapDirectoryImport,
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
			"distinguished_name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"password": {
				Type:     schema.TypeString,
				Required: true,
			},
			"url": {
				Type:     schema.TypeString,
				Required: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_profile_id_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled_user_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_search_attributes": {
				Type:     schema.TypeList,
				Optional: true,
				Elem:     &schema.Schema{Type: schema.TypeString},
				Default: []string{"*"},
			},
			"user_search_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_object_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_search_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_sync_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_object_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_match_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_canonical_name_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_email_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_display_name_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_given_name_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"user_family_name_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_enabled": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"tls_server_certificate": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_member_of_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_membership_attribute": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"connectivity_test_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"reject_unauthorized": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"disabled_group_filter": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"tags": {
				Type:     schema.TypeMap,
				Optional: true,
			},
		},
	}
}

func resourceTurbotLdapDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotLdapDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build mutation input

	input := mapFromResourceData(d, ldapDirectoryInputProperties)
	input["status"] = "NEW"

	ldapDirectory, err := client.CreateLdapDirectory(input)
	if err != nil {
		return err
	}

	// set parent_akas property by loading resource and fetching the akas
	if err := storeAkas(ldapDirectory.Turbot.ParentId, "parent_akas", d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(ldapDirectory.Turbot.Id)
	// assign properties coming back from create graphQl API
	d.Set("parent", ldapDirectory.Parent)
	d.Set("title", ldapDirectory.Title)
	d.Set("status", strings.ToUpper(ldapDirectory.Status))
	d.Set("directory_type", ldapDirectory.DirectoryType)
	// Set the values from Resource Data
	return nil
}

func resourceTurbotLdapDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()

	ldapDirectory, err := client.ReadLdapDirectory(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// local directory was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData
	d.Set("parent", ldapDirectory.Parent)
	d.Set("title", ldapDirectory.Title)
	d.Set("description", ldapDirectory.Description)
	d.Set("status", strings.ToUpper(ldapDirectory.Status))
	d.Set("profile_id_template", ldapDirectory.ProfileIdTemplate)
	d.Set("directory_type", ldapDirectory.DirectoryType)
	d.Set("tags", ldapDirectory.Turbot.Tags)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(ldapDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotLdapDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// build mutation payload
	input := mapFromResourceData(d, getldapDirectoryUpdateProperties())
	input["id"] = d.Id()
	// do update
	ldapDirectory, err := client.UpdateLdapDirectory(input)
	if err != nil {
		return err
	}

	// assign properties coming back from update graphQl API
	d.Set("parent", ldapDirectory.)
	d.Set("title", ldapDirectory.Title)
	d.Set("status", strings.ToUpper(ldapDirectory.Status))
	d.Set("directory_type", ldapDirectory.)
	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(ldapDirectory.Turbot.ParentId, "parent_akas", d, meta)
}

func resourceTurbotLdapDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	err := client.DeleteLdapDirectory(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")
	return nil
}

func resourceTurbotLdapDirectoryImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotLdapDirectoryRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}
