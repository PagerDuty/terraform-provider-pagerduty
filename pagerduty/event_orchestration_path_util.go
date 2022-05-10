package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var PagerDutyEventOrchestrationPathConditions = map[string]*schema.Schema{
	"expression": {
		Type:     schema.TypeString,
		Required: true,
	},
}
