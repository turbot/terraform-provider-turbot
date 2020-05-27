package turbot

import (
	"fmt"
	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/helper/schema"
	"github.com/terraform-providers/terraform-provider-turbot/apiClient"
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
		Timeouts: &schema.ResourceTimeout{
			Create: schema.DefaultTimeout(5 * time.Minute),
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

	var turbotResource *apiClient.Resource
	var err error
	errorCount := 0
	maxErrorRetries := 5
	err = resource.Retry(d.Timeout(schema.TimeoutCreate), func() *resource.RetryError {
		turbotResource, err = getResource(filter, resourceAka, client)
		// when we get NotFoundError, we retry for the timeout determined by the parameter TimeoutCreate, controlled by the config parameters timeouts.create (defaulting to 5 minutes). For other random/transient errors retry 5 times (maxErrorRetries)
		if err != nil {
			if apiClient.NotFoundError(err) {
				errorCount = 0
			} else {
				errorCount++
			}
			if errorCount == maxErrorRetries {
				return resource.NonRetryableError(err)
			}
			return resource.RetryableError(err)
		}
		return nil
	})
	if err != nil {
		return fmt.Errorf("turbot shadow resource creation failed: %s", err)
	}
	// assign the id
	d.SetId(turbotResource.Turbot.Id)
	return nil
}

func getResource(filter, resourceAka string, client *apiClient.Client) (*apiClient.Resource, error) {
	if resourceAka != "" {
		resource, err := client.ReadResource(resourceAka, nil)
		if err != nil {
			return nil, err
		}
		if resource.Turbot.Id != "" {
			return resource, nil
		}
		return nil, nil
	}
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
