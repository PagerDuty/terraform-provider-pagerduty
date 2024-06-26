package apiutil

import (
	"context"
	"time"

	"github.com/PagerDuty/terraform-provider-pagerduty/util"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/retry"
)

// AllFunc is a signature to use with function `All`, it receives the current
// number of items already listed, it returns a boolean signaling whether the
// system should keep requesting more items, and an error if any occured.
type AllFunc = func(offset int) (bool, error)

// Limit is the maximum amount of items a single request to PagerDuty's API
// should response
const Limit = 100

// All provides a boilerplate to request all pages from a list of a resource
// from PagerDuty's API
func All(ctx context.Context, requestFn AllFunc) error {
	offset := 0
	keepSearching := true

	for keepSearching {
		err := retry.RetryContext(ctx, 2*time.Minute, func() *retry.RetryError {
			more, err := requestFn(offset)

			if err != nil {
				if util.IsBadRequestError(err) {
					return retry.NonRetryableError(err)
				}
				return retry.RetryableError(err)
			}

			offset += Limit
			keepSearching = more
			return nil
		})

		if err != nil {
			return err
		}
	}

	return nil
}
