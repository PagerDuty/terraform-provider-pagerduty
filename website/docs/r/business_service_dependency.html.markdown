---
layout: "pagerduty"
page_title: "PagerDuty: pagerduty_business_service_dependency"
sidebar_current: "docs-pagerduty-resource-business-service-dependency"
description: |-
  Creates and manages a business service dependency in PagerDuty.
---

# pagerduty\_business\_service\_dependency

A [business service dependency](https://api-reference.pagerduty.com/#!/Business_Services/get_service_dependencies_business_services_id) is a relationship between a business service and technical and business services that this service uses, or that are used by this service, and are critical for successful operation.


## Example Usage

```hcl
resource "pagerduty_business_service_dependency" "foo" {
	relationship {
		supporting_service {
			id = pagerduty_business_service.foo.id
			type = "business_service"
		}
		dependent_service {
			id = pagerduty_service.foo.id
			type = "service"
		}
	}
	relationship {
		supporting_service {
			id = pagerduty_business_service.foo.id
			type = "business_service"
		}
		dependent_service {
			id = pagerduty_service.two.id
			type = "service"
		}
	}
}
```

## Argument Reference

The following arguments are supported:

  * `relationship` - (Required) The relationship between the `supporting_service` and `dependent_service`.
  * `supporting_service` - (Required) The service that supports  the  dependent service.
  * `dependent_service` - (Required) The service that id dependent on the supporting service.

