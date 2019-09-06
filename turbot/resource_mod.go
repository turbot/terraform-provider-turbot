package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
)

func resourceTurbotMod() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotModInstall,
		Read:   resourceTurbotModRead,
		Update: resourceTurbotModInstall,
		Delete: resourceTurbotModUninstall,
		Exists: resourceTurbotModExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotModImport,
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
			"org": {
				Type:     schema.TypeString,
				Required: true,
			},
			"mod": {
				Type:     schema.TypeString,
				Required: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				// default the version to any version
				Default:          "*",
				DiffSuppressFunc: supressIfLatestCompatibleVersionInstalled,
			},
			// TODO
			"latest_compatible_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTurbotModExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	// Exists - This is called to verify a resource still exists. It is called prior to Read,
	// and lowers the burden of Read to be able to assume the resource exists.
	log.Println("resourceTurbotModExists")
	client := meta.(*apiclient.Client)
	id := d.Id()

	_, err := client.ReadMod(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			return false, nil
		}
		return false, err
	}
	return true, nil
}

func resourceTurbotModInstall(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	org := d.Get("org").(string)
	modName := d.Get("mod").(string)
	version := d.Get("version").(string)

	// install mod returns turbot resource metadata containing the id
	turbotMetadata, err := client.InstallMod(parentAka, org, modName, version)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(turbotMetadata.Id)

	return nil
}

func resourceTurbotModRead(d *schema.ResourceData, meta interface{}) error {
	log.Println("resourceTurbotModRead")
	client := meta.(*apiclient.Client)
	id := d.Id()

	mod, err := client.ReadMod(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// mod was not found - clear id
			d.SetId("")
		}
		return err
	}
	// now load latest compatible version

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(d, meta); err != nil {
		return err
	}
	d.Set("parent", mod.Parent)
	d.Set("org", mod.Org)
	d.Set("mod", mod.Mod)
	d.Set("version", mod.Version)

	return nil
}

func resourceTurbotModUninstall(d *schema.ResourceData, meta interface{}) error {
	log.Println("resourceTurbotModUninstall")

	client := meta.(*apiclient.Client)
	id := d.Id()
	err := client.UninstallMod(id)
	if err != nil {
		return err
	}

	// clear the id to show we have deleted
	d.SetId("")

	return nil
}

func resourceTurbotModImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotModRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}

// the version in the config is a semver so may be a range. The version in the state file will be a specific version
// this will cause diffs to be identified
// supress diff if the latest compatible version is installed
func supressIfLatestCompatibleVersionInstalled(k, old, new string, d *schema.ResourceData) bool {
	return false
	//latestCompatibleVersion := d.Get("latest_compatible_version").(string)
	//return new == latestCompatibleVersion
}
