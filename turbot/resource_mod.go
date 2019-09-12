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
			"uri": {
				Type:     schema.TypeString,
				Computed: true,
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
			// store the version or version range specified in the config
			// this is necessary as if the config specifies a version range, the "version" field will contain the actual
			// installed version
			// we need to store the range in case a new version is released which satisfies the requirement
			"version_range": {
				Type:     schema.TypeString,
				Computed: true,
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
		return fmt.Errorf("Mod %s is already installed ( id: %s ). To manage this mod using Terraform, import the mod using command 'terraform import <resource_address> <id>'", modAka, id)
	}

	return modInstall(d, meta)
}

func resourceTurbotModUpdate(d *schema.ResourceData, meta interface{}) error {
	return modInstall(d, meta)
}

// do the actual mode installation
func modInstall(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	parentAka := d.Get("parent").(string)
	org := d.Get("org").(string)
	modName := d.Get("mod").(string)
	versionRange := d.Get("version").(string)

	// install mod returns turbot resource metadata containing the id
	mod, err := client.InstallMod(parentAka, org, modName, versionRange)
	if err != nil {
		log.Println("[ERROR] Turbot mod installation failed...", err)
		return err
	}

	modId := mod.Turbot.Id

	// now poll the mod resource to wait for the correct version
	installedVersion, err := waitForInstallation(modId, mod.Build, client)
	if err != nil {
		return err
	}

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(mod.Turbot.ParentId, d, meta); err != nil {
		return err
	}
	// assign the id
	d.SetId(modId)
	d.Set("latest_compatible_version", installedVersion)
	d.Set("version_range", versionRange)
	d.Set("uri", mod.Turbot.Akas[0])
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
	var targetVersion string
	// if 'version' is set in resourceData, fetch the latest version which satisfies this requirement
	if d.Get("version").(string) != "" {

		targetVersion, err = getLatestCompatibleVersion(d, meta)
		log.Printf("resourceTurbotModRead config version %s mod version %s latest %s", d.Get("version").(string), mod.Version, targetVersion)
		if err != nil {
			return err
		}
	} else {
		log.Printf("resourceTurbotModRead no version in resource data mod.Version %s", mod.Version)
		// if version is NOT set (e.g. for an import), just use the actual mod version as target version
		targetVersion = mod.Version
	}

	// assign results back into ResourceData

	// set parent_akas property by loading parent resource and fetching the akas
	if err = setParentAkas(mod.Parent, d, meta); err != nil {
		return err
	}
	d.Set("parent", mod.Parent)
	d.Set("org", mod.Org)
	d.Set("mod", mod.Mod)
	d.Set("version", mod.Version)
	d.Set("latest_compatible_version", targetVersion)
	d.Set("uri", mod.Uri)

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

func waitForInstallation(modId, targetBuild string, client *apiclient.Client) (string, error) {
	retryCount := 0
	// retry for 15 minutes
	maxRetries := 40
	sleep := 20 * time.Second
	log.Printf("Wait for mod installation, targetBuild: %s", targetBuild)

	for retryCount < maxRetries {
		installedVersion, installedBuild, err := getInstalledModVersion(modId, client)
		if err != nil {
			return "", err
		}
		if installedBuild == targetBuild {
			log.Printf("installed build: %s, target build: %s, mod is installed!", installedBuild, targetBuild)
			// success
			return installedVersion, nil
		}
		log.Printf("installed build: %s, target build: %s, retrying!", installedBuild, targetBuild)
		time.Sleep(sleep)
		retryCount++
	}
	return "", errors.New("Turbot mod installation timed out")
}

func getInstalledModVersion(modId string, client *apiclient.Client) (version, build string, err error) {
	properties := map[string]string{
		"version": "version",
		"build":   "build",
	}

	resource, err := client.ReadResource(modId, properties)
	if err != nil {
		return "", "", err
	}
	versionData := resource.Data["version"]
	buildData := resource.Data["build"]
	if (versionData == nil) || (buildData == nil) {
		return "", "", nil
	}
	version = versionData.(string)
	build = buildData.(string)
	return
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
// this may cause incorrect diffs to be identified
//
// supress diff if:
// - the current version ('old') is the latest compatible version,
// 	AND
// - the latest compatible version satisfies the new version requirements
//
// where
//   'old' is the currently installed version
//   'new' is the version field in the config
//   'latestCompatibleVersion' is the latest version which satisfies the _old_ config version  requirement
func supressIfLatestCompatibleVersionInstalled(_, old, new string, d *schema.ResourceData) bool {

	// latest compatible version is based on the old version parameter - if the new version has changed it is not valie
	latestCompatibleVersion := d.Get("latest_compatible_version").(string)
	log.Printf("supressIfLatestCompatibleVersionInstalled old: %s, new: %s, latestCompatibleVersion: %s", old, new, latestCompatibleVersion)
	c, err := semver.NewConstraint(new)
	if err != nil {
		return false
	}
	v, err := semver.NewVersion(latestCompatibleVersion)
	if err != nil {
		return false
	}

	return old == latestCompatibleVersion && c.Check(v)
}
