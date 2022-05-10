package pagerduty

import (
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"
)

var PagerDutyEventOrchestrationPathParent = map[string]*schema.Schema{
	"id": {
		Type:     schema.TypeString,
		Required: true,
	},
	"type": {
		Type:     schema.TypeString,
		Required: true,
	},
	"self": {
		Type:     schema.TypeString,
		Required: true,
	},
}

var PagerDutyEventOrchestrationPathConditions = map[string]*schema.Schema{
	"expression": {
		Type:     schema.TypeString,
		Required: true,
	},
}
