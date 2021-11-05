package pagerduty

import (
	"fmt"
	"log"
	"strings"
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
			State: resourcePagerDutyBusinessServiceSubscriberImport,
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
			// create subscriber assignment it as PagerDuty API does not return one
			assignmentID := createSubscriberID(businessServiceId, businessServiceSubscriber.Type, businessServiceSubscriber.ID)
			d.SetId(assignmentID)
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

func createSubscriberID(businessServiceId string, subscriberType string, subscriberID string) string {
	return fmt.Sprintf("%v.%v.%v", businessServiceId, subscriberType, subscriberID)
}

func resourcePagerDutyBusinessServiceSubscriberImport(d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	ids := strings.Split(d.Id(), ".")
	client, _ := meta.(*Config).Client()

	if len(ids) != 3 {
		return []*schema.ResourceData{}, fmt.Errorf("error importing pagerduty_business_service_subscriber. Expecting an importation ID formed as '<business_service_id>.<subscriber_type>.<subscriber_id>'")
	}

	businessServiceId, businessServiceSubscriberType, businessServiceSubscriberID := ids[0], ids[1], ids[2]
	subscriberResponse, _, err := client.BusinessServiceSubscribers.List(businessServiceId)
	if subscriberResponse != nil {
		// loop subscribers and find matching ID
		for _, subscriber := range subscriberResponse.BusinessServiceSubscribers {
			if subscriber.ID == businessServiceSubscriberID && subscriber.Type == businessServiceSubscriberType {
				// create subscriber assignment it as PagerDuty API does not return one
				assignmentID := createSubscriberID(businessServiceId, businessServiceSubscriberType, businessServiceSubscriberID)
				d.SetId(assignmentID)
				d.Set("business_service_id", businessServiceId)
				d.Set("subscriber_type", businessServiceSubscriberType)
				d.Set("subscriber_id", businessServiceSubscriberID)
				break
			}
		}
	}

	return []*schema.ResourceData{d}, err
}
