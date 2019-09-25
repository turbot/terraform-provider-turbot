package turbot

import (
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiclient"
)

// properties which must be passed to a create/update call
var shadowResourceProperties = []string{"filter"}

func resourceTurbotShadowResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotShadowResourceCreate,
		Read:   resourceTurbotShadowResourceRead,
		Delete: resourceTurbotShadowResourceDelete,
		Exists: resourceTurbotShadowResourceExists,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotShadowResourceImport,
		},
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Required: true,
			},
		},
	}
}

func resourceTurbotShadowResourceExists(d *schema.ResourceData, meta interface{}) (b bool, e error) {
	client := meta.(*apiclient.Client)
	id := d.Id()
	return client.ResourceExists(id)
}

func resourceTurbotShadowResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	filter := d.Get("filter").(string)
	// create folder returns turbot resource metadata containing the id
	resourceList, err := waitForResource(filter)
	if err != nil {
		log.Println("[ERROR] Turbot shadow resource creation failed...", err)
		return err
	}

	// assign the id
	d.SetId(resourceList.Id)

	return nil
}

func waitForResource(filter string, client *apiclient.Client) (string, error) {
	retryCount := 0
	// retry for 15 minutes
	maxRetries := 40
	sleep := 20 * time.Second
	log.Printf("Wait for the resource with filter: %s", filter)

	for retryCount < maxRetries {
		resourceList, err := getResource(filter, client)
		if err != nil {
			return "", err
		}
		if resourceList {
			log.Printf("resource with filter: %s", filter)
			// success
			return resourceList, nil
		}
		}
		log.Printf("no resource with filter: %s, retrying!", filter)
		time.Sleep(sleep)
		retryCount++
	}
	return "", errors.New("Fetching resource timed out")
}

func getResource(filter string, client *apiclient.Client) ( Resource, error) {

	list, err := client.ReadResourceList(filter, nil)
	if err != nil {
		return "", err
	}
	if len(list) == 1 {
		log.Printf("resource with filter: %s", filter)
		// success
		return resourceList, nil
	}
	if len(list) > 1 {
		return nil, errors.New("Try a better filter")

	}
	id := list.Data["id"]
	return
}

func resourceShadowResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiclient.Client)
	id := d.Id()

	folder, err := client.ReadShadowResource(id)
	if err != nil {
		if apiclient.NotFoundError(err) {
			// folder was not found - clear id
			d.SetId("")
		}
		return err
	}

	// assign results back into ResourceData

	d.Set("filter", ShadowResource.Filter)

	return nil
}

func resourceTurbotShadowResourceDelete(d *schema.ResourceData, meta interface{}) error {
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

func resourceTurbotShadowResourceImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	if err := resourceTurbotShadowResourceRead(d, meta); err != nil {
		return nil, err
	}
	return []*schema.ResourceData{d}, nil
}