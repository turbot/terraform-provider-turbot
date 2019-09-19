package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// these are the properties which must be passed to a create/update call
var googleDirectoryProperties = []string{"title", "pool_id", "profile_id_template", "group_id_template", "login_name_template", "client_secret", "hosted_name", "description"}

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
				// when doing a diff, the state file will contain the id of the parent bu tthe config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: supressIfParentAkaMatches,
			},
			// when doing a read, fetch the parent akas to use in supressIfParentAkaMatches()
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
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"status": {
				Type:     schema.TypeString,
				Required: true,
			},
			"directory_type": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_id": {
				Type:     schema.TypeString,
				Required: true,
			},
			"client_secret": {
				Type:     schema.TypeString,
				Required: true,
				DiffSuppressFunc: suppressIfClientSecret,
			},
			"pool_id": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"group_id_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"login_name_template": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"hosted_name": {
				Type:     schema.TypeString,
				Optional: true,
			},
		},
	}
}

func resourceTurbotGoogleDirectoryExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotGoogleDirectoryCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	// build map of local directory properties
	data := mapFromResourceData(d, googleDirectoryProperties)
	data["status"] = "New"
	data["directoryType"] = "google"
	data["clientID"] = d.Get("client_id").(string)
	turbotMetadata, err := client.CreateGoogleDirectory(parentAka, data)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}

	// assign the id
	d.SetId(turbotMetadata.Id)
	d.Set("status", data["status"])
	d.Set("directoryType", data["directoryType"])
	d.Set("clientSecret", "sensitive")
	return nil
}

func resourceTurbotGoogleDirectoryUpdate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	id := d.Id()

	// build map of local directory properties
	data := mapFromResourceData(d, googleDirectoryProperties)
	data["clientID"] = d.Get("client_id").(string)
	// create folder returns turbot resource metadata containing the id
	turbotMetadata, err := client.UpdateGoogleDirectory(id, parentAka, data)
	if err != nil {
		return err
	}
	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(turbotMetadata.ParentId, d, meta); err != nil {
		return err
	}
	return nil
}

func resourceTurbotGoogleDirectoryRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	googleDirectory, err := client.ReadGoogleDirectory(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// local directoery was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(googleDirectory.Turbot.ParentId, d, meta); err != nil {
		return err
	}
	d.Set("parent", googleDirectory.Parent)
	d.Set("title", googleDirectory.Title)
	d.Set("profile_id_template", googleDirectory.ProfileIdTemplate)
	d.Set("description", googleDirectory.Description)
	d.Set("status", googleDirectory.Status)
	d.Set("directory_type", googleDirectory.DirectoryType)
	d.Set("client_id", googleDirectory.ClientID)
	//d.Set("client_secret", googleDirectory.ClientSecret)
	d.Set("pool_id", googleDirectory.PoolId)
	d.Set("group_id_template", googleDirectory.GroupIdTemplate)
	d.Set("login_name_template", googleDirectory.LoginNameTemplate)
	d.Set("hosted_name", googleDirectory.HostedName)
	return nil
}

func resourceTurbotGoogleDirectoryDelete(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
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
