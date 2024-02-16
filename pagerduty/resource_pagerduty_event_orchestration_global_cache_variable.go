package pagerduty

import (
	"context"
	"github.com/hashicorp/terraform-plugin-sdk/v2/diag"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
	"github.com/heimweh/go-pagerduty/pagerduty"
)

func resourcePagerDutyEventOrchestrationGlobalCacheVariable() *schema.Resource {
	return &schema.Resource{
		CreateContext: resourcePagerDutyEventOrchestrationGlobalCacheVariableCreate,
		ReadContext:   resourcePagerDutyEventOrchestrationGlobalCacheVariableRead,
		UpdateContext: resourcePagerDutyEventOrchestrationGlobalCacheVariableUpdate,
		DeleteContext: resourcePagerDutyEventOrchestrationGlobalCacheVariableDelete,
		Importer: &schema.ResourceImporter{
			StateContext: resourcePagerDutyEventOrchestrationGlobalCacheVariableImport,
		},
		CustomizeDiff: checkConfiguration,
		Schema: map[string]*schema.Schema{
			"event_orchestration": {
				Type:     schema.TypeString,
				Required: true,
				ForceNew: true,
			},
			"id": {
				Type:     schema.TypeString,
				Computed: true,
			},
			"name": {
				Type:     schema.TypeString,
				Required: true,
			},
			"condition": {
				Type:     schema.TypeList,
				Optional: true,
				Elem: &schema.Resource{
					Schema: resourceEventOrchestrationCacheVariableConditionSchema,
				},
			},
			"configuration": {
				Type:     schema.TypeList,
				Required: true,
				MaxItems: 1,
				Elem: &schema.Resource{
					Schema: resourceEventOrchestrationCacheVariableConfigurationSchema,
				},
			},
			"disabled": {
				Type:     schema.TypeBool,
				Optional: true,
			},
		},
	}
}

func resourcePagerDutyEventOrchestrationGlobalCacheVariableImport(ctx context.Context, d *schema.ResourceData, meta interface{}) ([]*schema.ResourceData, error) {
	return resourceEventOrchestrationCacheVariableImport(ctx, d, meta, pagerduty.CacheVariableTypeGlobal)
}

func resourcePagerDutyEventOrchestrationGlobalCacheVariableCreate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceEventOrchestrationCacheVariableCreate(ctx, d, meta, pagerduty.CacheVariableTypeGlobal)
}
func resourcePagerDutyEventOrchestrationGlobalCacheVariableRead(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceEventOrchestrationCacheVariableRead(ctx, d, meta, pagerduty.CacheVariableTypeGlobal)
}
func resourcePagerDutyEventOrchestrationGlobalCacheVariableUpdate(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceEventOrchestrationCacheVariableUpdate(ctx, d, meta, pagerduty.CacheVariableTypeGlobal)
}
func resourcePagerDutyEventOrchestrationGlobalCacheVariableDelete(ctx context.Context, d *schema.ResourceData, meta interface{}) diag.Diagnostics {
	return resourceEventOrchestrationCacheVariableDelete(ctx, d, meta, pagerduty.CacheVariableTypeGlobal)
}
