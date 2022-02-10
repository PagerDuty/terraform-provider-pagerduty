package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"log"
	"time"

	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/nordcloud/go-pagerduty/pagerduty"
)

func resourcePagerDutyWebhookSubscription() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyWebhookSubscriptionCreate,
		ReadContext:   resourcePagerDutyWebhookSubscriptionRead,
		UpdateContext: resourcePagerDutyWebhookSubscriptionUpdate,
		DeleteContext: resourcePagerDutyWebhookSubscriptionDelete,
		Importer: &schema.ResourceImporter{
			StateContext: schema.ImportStatePassthroughContext,
		},
		Schema: map[string]*schema.Schema{
			"delivery_method": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"temporarily_disabled": {
							Type:     schema.TypeBool,
							Optional: true,
							Computed: true,
						},
						"type": {
							Type:     schema.TypeString,
							Optional: true,
							Default:  "http_delivery_method",
							ValidateFunc: validateValueFunc([]string{
								"http_delivery_method",
							}),
						},
						"url": {
							Type:     schema.TypeString,
							Optional: true,
						},
					},
				},
			},
			"type": {
				Type:     schema.TypeString,
				Default:  "webhook_subscription",
				Optional: true,
				ValidateFunc: validateValueFunc([]string{
					"webhook_subscription",
				}),
			},
			"active": {
				Type:     schema.TypeBool,
				Optional: true,
				Computed: true,
			},
			"description": {
				Type:     schema.TypeString,
				Optional: true,
			},
			"events": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Schema{
					Type: schema.TypeString,
				},
			},
			"filter": {
				Type:     schema.TypeList,
				Required: true,
				Elem: &schema.Resource{
					Schema: map[string]*schema.Schema{
						"id": {
							Type:     schema.TypeString,
							Optional: true,
						},
						"type": {
							Type:     schema.TypeString,
							Required: true,
							ValidateFunc: validateValueFunc([]string{
								"service_reference",
								"team_reference",
								"account_reference",
							}),
						},
					},
				},
			},
		},
	}
}

func buildWebhookSubscriptionStruct(d *schema.ResourceData) *pagerduty.WebhookSubscription {
	webhook := pagerduty.WebhookSubscription{
		Type:           d.Get("type").(string),
		Active:         d.Get("active").(bool),
		Description:    d.Get("description").(string),
		DeliveryMethod: expandDeliveryMethod(d.Get("delivery_method").(interface{})),
		Events:         expandConfigList(d.Get("events").([]interface{})),
		Filter:         expandFilter(d.Get("filter").(interface{})),
	}
	return &webhook
}

func resourcePagerDutyWebhookSubscriptionCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	webhook := buildWebhookSubscriptionStruct(d)

	log.Printf("[INFO] Creating PagerDuty webhook subscription to be delivered to %s", webhook.DeliveryMethod.URL)

	retryErr := resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		if webhook, _, err := client.WebhookSubscriptions.Create(webhook); err != nil {
			if isErrCode(err, 400) || isErrCode(err, 429) {
				return resource.RetryableError(err)
			}

			return resource.NonRetryableError(err)
		} else if webhook != nil {
			d.SetId(webhook.ID)
		}
		return nil
	})

	if retryErr != nil {
		return diag.FromErr(retryErr)
	}

	return resourcePagerDutyWebhookSubscriptionRead(ctx, d, meta)

}

func resourcePagerDutyWebhookSubscriptionRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Reading PagerDuty webhook subscription %s", d.Id())

	return diag.FromErr(resource.RetryContext(ctx, 10*time.Minute, func() *resource.RetryError {
		webhook, _, err := client.WebhookSubscriptions.Get(d.Id())
		if checkErr := handleGenericErrors(err, d); checkErr.ShouldReturn {
			return checkErr.ReturnVal
		}

		if webhook != nil {
			setWebhookResourceData(d, webhook)
		}
		return nil
	}))
}
func resourcePagerDutyWebhookSubscriptionUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Updating PagerDuty webhook subscription %s", d.Id())
	whStruct := buildWebhookSubscriptionStruct(d)

	webhook, _, err := client.WebhookSubscriptions.Update(d.Id(), whStruct)
	if err != nil {
		return diag.FromErr(err)
	} else if webhook != nil {
		setWebhookResourceData(d, webhook)
	}

	return nil
}

func resourcePagerDutyWebhookSubscriptionDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	client, err := meta.(*Config).Client()
	if err != nil {
		return diag.FromErr(err)
	}

	log.Printf("[INFO] Deleting PagerDuty webhook subscription %s", d.Id())

	if _, err := client.WebhookSubscriptions.Delete(d.Id()); err != nil {
		return diag.FromErr(err)
	}

	d.SetId("")

	return nil
}

func setWebhookResourceData(d *schema.ResourceData, webhook *pagerduty.WebhookSubscription) {
	d.Set("type", webhook.Type)
	d.Set("active", webhook.Active)
	d.Set("description", webhook.Description)
	d.Set("events", flattenConfigList(webhook.Events))
	d.Set("delivery_method", flattenDeliveryMethod(webhook.DeliveryMethod))
	d.Set("filter", flattenFilter(webhook.Filter))
}

func expandDeliveryMethod(v interface{}) pagerduty.DeliveryMethod {
	dmMap := v.([]interface{})[0].(map[string]interface{})

	var method pagerduty.DeliveryMethod

	method = pagerduty.DeliveryMethod{
		TemporarilyDisabled: dmMap["temporarily_disabled"].(bool),
		Type:                dmMap["type"].(string),
		URL:                 dmMap["url"].(string),
	}
	return method
}
func expandFilter(v interface{}) pagerduty.Filter {
	filterMap := v.([]interface{})[0].(map[string]interface{})

	var filter pagerduty.Filter

	filter = pagerduty.Filter{
		ID:   filterMap["id"].(string),
		Type: filterMap["type"].(string),
	}
	return filter
}

func flattenDeliveryMethod(method pagerduty.DeliveryMethod) []map[string]interface{} {
	var methods []map[string]interface{}
	methodMap := map[string]interface{}{
		"temporarily_disabled": method.TemporarilyDisabled,
		"type":                 method.Type,
		"url":                  method.URL,
	}
	methods = append(methods, methodMap)
	return methods
}

func flattenFilter(filter pagerduty.Filter) []map[string]interface{} {
	var filters []map[string]interface{}
	filterMap := map[string]interface{}{
		"id":   filter.ID,
		"type": filter.Type,
	}
	filters = append(filters, filterMap)
	return filters
}
