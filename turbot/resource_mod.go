package turbot

import (
	"errors"
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
	"log"
	"time"
)

func resourceTurbotMod() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotModInstall,
		Read:   resourceTurbotModRead,
		Update: resourceTurbotModUpdate,
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
				ForceNew:         true,
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
				ForceNew: true,
			},
			"mod": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"version": {
				Type:     schema.TypeString,
				Optional: true,
				// default the version to any version
				Default:          "*",
				DiffSuppressFunc: supressIfLatestCompatibleVersionInstalled,
			},
			// store latest version which satisfies the version requirement
			"latest_compatible_version": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
	}
}

func resourceTurbotModExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotModInstall(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	org := d.Get("org").(string)
	modName := d.Get("mod").(string)
	modAka := buildModAka(org, modName)

	// install should only be called if the mod is not already installed
	mod, err := client.ReadResource(modAka, nil)
	if err != nil {
		return err
	}
	id := mod.Turbot.Id
	if id != "" {
		// TODO extract terraform name
		return fmt.Errorf("Cannot install mod %s as it is already installed. To manange this mod using Terraform, import the mod using command 'terraform import <resource_address> %s'", modAka, id)
	}

	return modInstall(d, meta)
}
func resourceTurbotModUpdate(d *schema.ResourceData, meta interface{}) error {
	return modInstall(d, meta)
}

// do tha eactual mode insatallation
func modInstall(d *schema.ResourceData, meta interface{}) error {

	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	org := d.Get("org").(string)
	modName := d.Get("mod").(string)
	version := d.Get("version").(string)

	// now determine latest compatible version
	targetVersion, err := getLatestCompatibleVersion(d, meta)
	if err != nil {
		return err
	}

	// install mod returns turbot resource metadata containing the id
	mod, err := client.InstallMod(parentAka, org, modName, version)
	if err != nil {
		log.Println("[ERROR] Turbot mod installation failed...", err)
		return err
	}

	modId := mod.Turbot.Id

	// now poll the mod resource to wait for the correct version
	if err = waitForInstallation(modId, targetVersion, client); err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(modId)
	d.Set("latest_compatible_version", targetVersion)

	return nil
}

func resourceTurbotModRead(d *schema.ResourceData, meta interface{}) error {
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
	// now determine latest compatible version
	targetVersion, err := getLatestCompatibleVersion(d, meta)
	if err != nil {
		return err
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(d, meta); err != nil {
		return err
	}
	d.Set("parent", mod.Parent)
	d.Set("org", mod.Org)
	d.Set("mod", mod.Mod)
	d.Set("version", mod.Version)
	d.Set("latest_compatible_version", targetVersion)

	return nil
}

func resourceTurbotModUninstall(d *schema.ResourceData, meta interface{}) error {
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

func buildModAka(org, mod string) string {
	return fmt.Sprintf("tmod:@%s/%s", org, mod)
}
func waitForInstallation(modId, targetVersion string, client *apiclient.Client) error {
	retryCount := 0
	// retry for 15 minutes
	maxRetries := 40
	sleep := 20 * time.Second
	log.Printf("[DEBUG] Wait for mod installation, targetVersion: %s", targetVersion)

	for retryCount < maxRetries {
		installedVersion, err := getInstalledModVersion(modId, client)
		if err != nil {
			return err
		}
		log.Println("[DEBUG] Installed version: ", installedVersion)
		if installedVersion == targetVersion {
			log.Printf("[DEBUG] Installed version: %s, target version: %s, mod is installed!", installedVersion, targetVersion)
			// success
			return nil
		}
		log.Printf("[DEBUG] Installed version: %s, target version: %s, retrying!", installedVersion, targetVersion)
		time.Sleep(sleep)
		retryCount++
	}

	return errors.New("Turbot mod installation timed out")
}

func getInstalledModVersion(modId string, client *apiclient.Client) (string, error) {
	properties := map[string]string{
		"version": "turbot.custom.installedVersion",
	}

	resource, err := client.ReadResource(modId, properties)
	if err != nil {
		return "", err
	}
	if resource.Data["version"] == nil {
		return "", nil
	}

	return resource.Data["version"].(string), nil
}

func getLatestCompatibleVersion(d *schema.ResourceData, meta interface{}) (string, error) {
	client := meta.(*apiclient.Client)
	org := d.Get("org").(string)
	modName := d.Get("mod").(string)
	version := d.Get("version").(string)
	modVersions, err := client.GetModVersions(org, modName)
	if err != nil {
		return "", err
	}

	c, err := semver.NewConstraint(version)
	if err != nil {
		return "", err
	}

	// now get latest version
	latestCompatibleVersion := ""
	for _, modVersion := range modVersions {
		if modVersion.Status == "available" {
			v, err := semver.NewVersion(modVersion.Version)
			if err != nil {
				return "", err
			}
			// does this version meet the requirement
			if c.Check(v) {
				latestCompatibleVersion = modVersion.Version
			}
		}
	}
	return latestCompatibleVersion, nil

}

// the version in the config is a semver so may be a range. The version in the state file will be a specific version
// this will cause diffs to be identified
// supress diff if the latest compatible version is installed
func supressIfLatestCompatibleVersionInstalled(k, old, new string, d *schema.ResourceData) bool {
	return false
	//latestCompatibleVersion := d.Get("latest_compatible_version").(string)
	//return new == latestCompatibleVersion
}
