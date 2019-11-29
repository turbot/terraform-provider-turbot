package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
	"log"
	"time"
)

// properties which must be passed to a create/update call
var shadowResourceProperties = []interface{}{"filter", "resource"}

func resourceTurbotShadowResource() *schema.Resource {
	return &schema.Resource{
		Create: resourceTurbotShadowResourceCreate,
		Read:   resourceTurbotShadowResourceRead,
		Delete: resourceTurbotShadowResourceDelete,
		Importer: &schema.ResourceImporter{
			State: resourceTurbotShadowResourceImport,
		},
		Schema: map[string]*schema.Schema{
			"filter": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
			"resource": {
				Type:     schema.TypeString,
				Optional: true,
				ForceNew: true,
			},
		},
	}
}

func resourceTurbotShadowResourceCreate(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)

	filter := d.Get("filter").(string)
	resourceAka := d.Get("resource").(string)
	if filter == "" && resourceAka == "" {
		return fmt.Errorf("one of resource or filter should be specified")
	}

	if filter != "" && resourceAka != "" {
		return fmt.Errorf("resource and filter must not both be specified")
	}

	if resourceAka != "" {
		filter = fmt.Sprintf("resource:%s level: self", resourceAka)
	}
	// create folder returns turbot resource metadata containing the id
	resource, err := waitForResource(filter, client)
	if err != nil {
		log.Println("[ERROR] Turbot shadow resource creation failed...", err)
		return err
	}

	// assign the id
	d.SetId(resource.Turbot.Id)
	return nil
}

func waitForResource(filter string, client *apiClient.Client) (*apiClient.Resource, error) {
	retryCount := 0
	// retry for 5 minutes
	timeoutMins := 5
	retryIntervalSecs := 5
	maxRetries := (timeoutMins * 60) / retryIntervalSecs
	sleep := time.Duration(retryIntervalSecs) * time.Second
	for retryCount < maxRetries {
		resource, err := getResource(filter, client)
		if err != nil {
			return nil, err
		}
		if resource != nil {
			log.Printf("found resource")
			// success
			return resource, nil
		}
		time.Sleep(sleep)
		retryCount++
	}
	return nil, fmt.Errorf("fetching resource with filter timed out after %d minutes", timeoutMins)
}

func getResource(filter string, client *apiClient.Client) (*apiClient.Resource, error) {
	resourceList, err := client.ReadResourceList(filter, nil)
	if err != nil {
		return nil, err
	}
	if len(resourceList) == 1 {
		// success
		return &resourceList[0], nil
	}
	if len(resourceList) > 1 {
		return nil, fmt.Errorf("filter \"%s\" returned %d items. Specify a filter returning a single item", filter, len(resourceList))
	}
	return nil, nil
}

func resourceTurbotShadowResourceRead(d *schema.ResourceData, meta interface{}) error {
	client := meta.(*apiClient.Client)
	id := d.Id()
	exists, err := client.ResourceExists(id)
	if err != nil {
		return err
	}
	if !exists {
		d.SetId("")
		return nil
	}
	d.Set("filter", d.Get("filter"))
	return nil
}

func resourceTurbotShadowResourceDelete(d *schema.ResourceData, meta interface{}) error {
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
