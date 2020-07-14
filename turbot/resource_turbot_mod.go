package turbot

import (
	"fmt"
	"github.com/Masterminds/semver"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"log"
	"strings"
	"time"
)

var modInputProperties = []interface{}{"parent", "org", "mod", "version"}

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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(15 * time.Minute),
		},
		Schema: map[string]*schema.Schema{
			// aka of the parent resource
			"parent": {
				Type:     schema.TypeString,
				Optional: true,
				// when doing a diff, the state file will contain the id of the parent but the config contains the aka,
				// so we need custom diff code
				DiffSuppressFunc: suppressIfAkaMatches("parent_akas"),
				ForceNew:         true,
				Default:          "tmod:@turbot/turbot#/",
			},
			// when doing a read, fetch the parent akas to use in suppressIfAkaMatches
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
				Default: "*",
			},
			// store the version currently installed (as the 'version' property may be a range)
			"version_current": {
				Type:     schema.TypeString,
				Computed: true,
			},
			// store latest version which satisfies the version requirement
			"version_latest": {
				Type:     schema.TypeString,
				Computed: true,
			},
		},
		CustomizeDiff: resourceTurbotModCustomizeDiff,
	}
}

func resourceTurbotModCustomizeDiff(d *schema.ResourceDiff, meta interface{}) error {

	versionCurrent := d.Get("version_current").(string)
	var versionLatest string
	// if the version has changed, re-fetch the latest compatible version to detect if we need to change the installed version
	if d.HasChange("version") {
		var err error
		org := d.Get("org").(string)
		modName := d.Get("mod").(string)
		version := d.Get("version").(string)
		versionLatest, err = getLatestCompatibleVersion(org, modName, version, meta)
		if err != nil {
			return err
		}
	} else {
		// otherwise if version has not changed, use the saved value of version_latest
		versionLatest = d.Get("version_latest").(string)
	}
	// if the current version is not the latest which satisfied the version requirements, raise a diff
	if versionCurrent != versionLatest {
		if err := d.SetNew("version_current", versionLatest); err != nil {
			return err
		}
	}
	return nil
}

func resourceTurbotModExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiClient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotModInstall(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	org := d.Get("org").(string)
	modName := d.Get("mod").(string)
	modAka := buildModAka(org, modName)

	// install should only be called if the mod is not already installed
	mod, err := client.ReadResource(modAka, nil)
	if err == nil {
		// if there is no error, the mod is already installed
		id := mod.Turbot.Id
		return fmt.Errorf("mod %s is already installed ( id: %s ). To manage this mod using Terraform, import the mod using command 'terraform import <resource_address> <id>'", modAka, id)
	}
	if !apiClient.NotFoundError(err) {
		// if the error is not a 'not found' error, the mod is already installed
		return err
	}

	return modInstall(d, meta)
}

func resourceTurbotModUpdate(d *schema.ResourceData, meta interface{}) error {
	return modInstall(d, meta)
}

// do the actual mode installation
func modInstall(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	// install mod returns turbot resource metadata containing the id
	input := mapFromResourceData(d, modInputProperties)
	mod, err := client.InstallMod(input)
	if err != nil {
		log.Println("[ERROR] Turbot mod installation failed...", err)
		return err
	}

	modId := mod.Turbot.Id
	// now poll the mod resource to wait for the correct version
	targetBuild := mod.Build
	log.Printf("Wait for mod installation, targetBuild: %s", targetBuild)
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		installedVersion, installedBuild, err := getInstalledModVersion(modId, client)
		if installedBuild == targetBuild {
			log.Printf("installed version: %s, installed build: %s, target build: %s, mod is installed!", installedVersion, installedBuild, targetBuild)
			// success
			return nil
		}
		if err == nil {
			return resource.RetryableError(err)
		}
		if err != nil {
			return resource.NonRetryableError(err)
		}
		return resource.NonRetryableError(err)
	})
	if err != nil {
		return err
	}

	// assign the id
	d.SetId(modId)
	return resourceTurbotModRead(d, meta)
}

func resourceTurbotModRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	mod, err := client.ReadMod(id)
	if err != nil {
		if apiClient.NotFoundError(err) {
			// mod was not found - clear id
			d.SetId("")
		}
		return err
	}
	// now determine latest compatible version
	var targetVersion string
	// if 'version' is set in resourceData, fetch the latest version which satisfies this requirement
	if version := d.Get("version").(string); version != "" {
		org := d.Get("org").(string)
		modName := d.Get("mod").(string)
		targetVersion, err = getLatestCompatibleVersion(org, modName, version, meta)
		log.Printf("resourceTurbotModRead config version %s installed version %s latest version%s", version, mod.Version, targetVersion)
		if err != nil {
			return err
		}
	} else {
		log.Printf("resourceTurbotModRead no version in resource data mod.Version %s", mod.Version)
		// if version is NOT set in resource data (e.g. for an import), just use the actual mod version and targetVersion
		targetVersion = mod.Version
	}

	// assign results back into ResourceData

	d.Set("parent", mod.Parent)
	d.Set("org", mod.Org)
	d.Set("mod", mod.Mod)
	d.Set("version_current", mod.Version)
	d.Set("version_latest", targetVersion)
	d.Set("uri", mod.Uri)

	// set parent_akas property by loading resource and fetching the akas
	return storeAkas(mod.Parent, "parent_akas", d, meta)
}

func resourceTurbotModUninstall(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
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

func getInstalledModVersion(modId string, client *apiClient.Client) (version, build string, err error) {
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

func getLatestCompatibleVersion(org, modName, version string, meta interface{}) (string, error) {
	client := meta.(*apiClient.Client)
	modVersions, err := client.GetModVersions(org, modName)
	if err != nil {
		return "", err
	}

	// create semver constraint from required version range
	c, err := semver.NewConstraint(version)
	if err != nil {
		return "", err
	}

	// now get latest version
	var latestVersion *semver.Version
	for _, modVersion := range modVersions {
		modStatus := strings.ToLower(modVersion.Status)
		if modStatus == "available" || modStatus == "recommended" {
			// create semver version from this version
			v, err := semver.NewVersion(modVersion.Version)
			if err != nil {
				return "", err
			}
			// does this version meet the requirement
			if c.Check(v) && (latestVersion == nil || v.GreaterThan(latestVersion)) {
				latestVersion = v
			}
		}
	}
	latestVersionString := ""
	if latestVersion != nil {
		latestVersionString = latestVersion.String()
	}
	return latestVersionString, nil

}
