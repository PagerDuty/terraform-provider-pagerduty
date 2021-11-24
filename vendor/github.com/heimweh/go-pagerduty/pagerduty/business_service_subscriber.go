package pagerduty

import (
	"errors"
	"fmt"
)

// BusinessServiceSubscriberService handles the communication with business service
// subscriber related methods of the PagerDuty API.
type BusinessServiceSubscriberService service

// BusinessService represents a business service.
type BusinessServiceSubscriber struct {
	ID               string `json:"subscriber_id,omitempty"`
	Type             string `json:"subscriber_type,omitempty"`
	SubscribableID   string `json:"subscribable_id,omitempty"`
	SubscribableType string `json:"subscribable_type,omitempty"`
	Result           string `json:"result,omitempty"`
}

// BusinessServiceSubscriberPayload represents payload with a business service subscriber object
type BusinessServiceSubscriberPayload struct {
	BusinessServiceSubscriber []*BusinessServiceSubscriber `json:"subscribers,omitempty"`
}

// CreateBusinessServiceSubscribersResponse represents a create response of business service subscription result.
type CreateBusinessServiceSubscribersResponse struct {
	BusinessServiceSubscriber []*BusinessServiceSubscriber `json:"subscriptions,omitempty"`
}

// ListBusinessServiceSubscribersResponse represents a list response of business service subscribers.
type ListBusinessServiceSubscribersResponse struct {
	Total                      int                          `json:"total,omitempty"`
	BusinessServiceSubscribers []*BusinessServiceSubscriber `json:"subscribers,omitempty"`
	Offset                     int                          `json:"offset,omitempty"`
	More                       bool                         `json:"more,omitempty"`
	Limit                      int                          `json:"limit,omitempty"`
}

// List lists existing business service subscribers.
func (s *BusinessServiceSubscriberService) List(businessServiceID string) (*ListBusinessServiceSubscribersResponse, *Response, error) {
	u := fmt.Sprintf("/business_services/%s/subscribers", businessServiceID)
	v := new(ListBusinessServiceSubscribersResponse)

	businessServiceSubscribers := make([]*BusinessServiceSubscriber, 0)

	// Create a handler closure capable of parsing data from the subscribers endpoint
	// and appending resultant response plays to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListBusinessServiceSubscribersResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		businessServiceSubscribers = append(businessServiceSubscribers, result.BusinessServiceSubscribers...)

		// Return stats on the current page. Caller can use this information to
		// adjust for requesting additional pages.
		return ListResp{
			More:   result.More,
			Offset: result.Offset,
			Limit:  result.Limit,
		}, response, nil
	}
	err := s.client.newRequestPagedGetDo(u, responseHandler)
	if err != nil {
		return nil, nil, err
	}
	v.BusinessServiceSubscribers = businessServiceSubscribers

	return v, nil, nil
}

// Create creates a new business service subscriber.
func (s *BusinessServiceSubscriberService) Create(businessServiceID string, subscriber *BusinessServiceSubscriber) (*Response, error) {
	u := fmt.Sprintf("/business_services/%s/subscribers", businessServiceID)
	v := new(BusinessServiceSubscriberPayload)
	subscriberArr := make([]*BusinessServiceSubscriber, 0)
	subscriberArr = append(subscriberArr, subscriber)
	p := &BusinessServiceSubscriberPayload{BusinessServiceSubscriber: subscriberArr}

	resp, err := s.client.newRequestDo("POST", u, nil, p, v)
	if err != nil {
		return nil, err
	}

	var result CreateBusinessServiceSubscribersResponse

	if err := s.client.DecodeJSON(resp, &result); err != nil {
		return nil, err
	}

	subscriptionResp := result.BusinessServiceSubscriber
	errorMessage := ""
	for _, subscription := range subscriptionResp {
		if subscription.Result != "success" {
			// append error message to message variable
			errorMessage = errorMessage + fmt.Sprintf("resulting status for subscription of %s %s to %s %s was: %s. ", subscription.Type, subscription.ID, subscription.SubscribableType, subscription.SubscribableID, subscription.Result)
		}
	}

	if errorMessage != "" {
		return nil, errors.New(errorMessage)
	}

	return resp, nil
}

// Delete deletes a business service subscriber.
func (s *BusinessServiceSubscriberService) Delete(businessServiceID string, subscriber *BusinessServiceSubscriber) (*Response, error) {
	u := fmt.Sprintf("/business_services/%s/unsubscribe", businessServiceID)
	v := new(BusinessServiceSubscriberPayload)
	subscriberArr := make([]*BusinessServiceSubscriber, 0)
	subscriberArr = append(subscriberArr, subscriber)
	p := &BusinessServiceSubscriberPayload{BusinessServiceSubscriber: subscriberArr}

	resp, err := s.client.newRequestDo("POST", u, nil, p, v)
	if err != nil {
		return nil, err
	}

	return resp, nil
}
