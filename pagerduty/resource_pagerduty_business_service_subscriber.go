package pagerduty

import (
	"log"
	"time"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyBusinessServiceSubscriber() *schema.Resource {
	return &schema.Resource{
		Create: resourcePagerDutyBusinessServiceSubscriberCreate,
		Read:   resourcePagerDutyBusinessServiceSubscriberRead,
		Delete: resourcePagerDutyBusinessServiceSubscriberDelete,
		Importer: &schema.ResourceImporter{
			State: schema.ImportStatePassthrough,
		},
		Schema: map[string]*schema.Schema{
			"subscriber_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"subscriber_type": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
				ValidateFunc: validateValueFunc([]string{
					"team",
					"user",
				}),
			},
			"business_service_id": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
		},
	}
}

func buildBusinessServiceSubscriberStruct(d *schema.ResourceData) (*pagerduty.BusinessServiceSubscriber, error) {
	subscriber := pagerduty.BusinessServiceSubscriber{
		ID:   d.Get("subscriber_id").(string),
		Type: d.Get("subscriber_type").(string),
	}

	return &subscriber, nil
}

func resourcePagerDutyBusinessServiceSubscriberCreate(d *schema.ResourceData, meta interface{}) error {
	client, _ := meta.(*Config).Client()
	businessServiceId := d.Get("business_service_id").(string)

	retryErr := resource.Retry(2*time.Minute, func() *resource.RetryError {

		businessServiceSubscriber, err := buildBusinessServiceSubscriberStruct(d)
		if err != nil {
			return resource.NonRetryableError(err)
		}

		log.Printf("[INFO] Creating PagerDuty business service %s subscriber %s type %s", businessServiceId, businessServiceSubscriber.ID, businessServiceSubscriber.Type)
		if _, err = client.BusinessServiceSubscribers.Create(businessServiceId, businessServiceSubscriber); err != nil {
			return resource.RetryableError(err)
		} else if businessServiceSubscriber != nil {
			d.SetId(businessServiceSubscriber.ID)
		}
		return nil
	})
	if retryErr != nil {
		time.Sleep(2 * time.Second)
		return retryErr
	}

	return resourcePagerDutyBusinessServiceSubscriberRead(d, meta)
}

func resourcePagerDutyBusinessServiceSubscriberRead(d *schema.ResourceData, meta interface{}) error {
	client, _ := meta.(*Config).Client()

	businessServiceId := d.Get("business_service_id").(string)
	businessServiceSubscriber, _ := buildBusinessServiceSubscriberStruct(d)

	log.Printf("[INFO] Reading PagerDuty business service %s subscriber %s type %s", businessServiceId, businessServiceSubscriber.ID, businessServiceSubscriber.Type)

	return resource.Retry(30*time.Second, func() *resource.RetryError {
		if subscriberResponse, _, err := client.BusinessServiceSubscribers.List(businessServiceId); err != nil {
			time.Sleep(2 * time.Second)
			return resource.RetryableError(err)
		} else if subscriberResponse != nil {
			var foundSubscriber *pagerduty.BusinessServiceSubscriber

			// loop subscribers and find matching ID
			for _, subscriber := range subscriberResponse.BusinessServiceSubscribers {
				if subscriber.ID == businessServiceSubscriber.ID && subscriber.Type == businessServiceSubscriber.Type {
					foundSubscriber = subscriber
					break
				}
			}
			if foundSubscriber == nil {
				d.SetId("")
				return nil
			}
		}
		return nil
	})
}

func resourcePagerDutyBusinessServiceSubscriberDelete(d *schema.ResourceData, meta interface{}) error {
	client, _ := meta.(*Config).Client()

	businessServiceId := d.Get("business_service_id").(string)
	businessServiceSubscriber, _ := buildBusinessServiceSubscriberStruct(d)

	log.Printf("[INFO] Deleting PagerDuty business service %s subscriber %s type %s", businessServiceId, businessServiceSubscriber.ID, businessServiceSubscriber.Type)

	if _, err := client.BusinessServiceSubscribers.Delete(businessServiceId, businessServiceSubscriber); err != nil {
		return err
	}

	d.SetId("")

	return nil
}
