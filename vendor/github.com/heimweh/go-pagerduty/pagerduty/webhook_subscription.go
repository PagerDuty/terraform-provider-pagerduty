package pagerduty

import "fmt"

// WebhookSubscriptionService handle v3 webhooks from PagerDuty.
type WebhookSubscriptionService service

// WebhookSubscription represents a webhook subscription.
type WebhookSubscription struct {
	ID             string         `json:"id,omitempty"`
	Type           string         `json:"type,omitempty"`
	Active         bool           `json:"active,omitempty"`
	Description    string         `json:"description,omitempty"`
	DeliveryMethod DeliveryMethod `json:"delivery_method,omitempty"`
	Events         []string       `json:"events,omitempty"`
	Filter         Filter         `json:"filter,omitempty"`
}

// DeliveryMethod represents a webhook delivery method
type DeliveryMethod struct {
	TemporarilyDisabled bool             `json:"temporarily_disabled,omitempty"`
	Type                string           `json:"type,omitempty"`
	URL                 string           `json:"url,omitempty"`
	CustomHeaders       []*CustomHeaders `json:"custom_headers"`
	Secret              string           `json:"secret,omitempty"`
}

type CustomHeaders struct {
	Name  string `json:"name,omitempty"`
	Value string `json:"value,omitempty"`
}

// Filter represents a webhook subscription filter
type Filter struct {
	ID   string `json:"id,omitempty"`
	Type string `json:"type,omitempty"`
}

// ListWebhookSubscriptionsResponse represents a list response of webhook subscriptions.
type ListWebhookSubscriptionsResponse struct {
	Total                int                    `json:"total,omitempty"`
	WebhookSubscriptions []*WebhookSubscription `json:"webhook_subscriptions,omitempty"`
	Offset               int                    `json:"offset,omitempty"`
	More                 bool                   `json:"more,omitempty"`
	Limit                int                    `json:"limit,omitempty"`
}

// WebhookSubscriptionPayload represents payload with a slack connect object
type WebhookSubscriptionPayload struct {
	WebhookSubscription *WebhookSubscription `json:"webhook_subscription,omitempty"`
}

// List lists existing webhook subscriptions.
func (s *WebhookSubscriptionService) List() (*ListWebhookSubscriptionsResponse, *Response, error) {
	u := "/webhook_subscriptions"
	v := new(ListWebhookSubscriptionsResponse)

	webhookSubscriptions := make([]*WebhookSubscription, 0)

	// Create a handler closure capable of parsing data from the webhook subscriptions endpoint
	// and appending resultant response plays to the return slice.
	responseHandler := func(response *Response) (ListResp, *Response, error) {
		var result ListWebhookSubscriptionsResponse

		if err := s.client.DecodeJSON(response, &result); err != nil {
			return ListResp{}, response, err
		}

		webhookSubscriptions = append(webhookSubscriptions, result.WebhookSubscriptions...)

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
	v.WebhookSubscriptions = webhookSubscriptions

	return v, nil, nil
}

// Create creates a new webhook subscription.
func (s *WebhookSubscriptionService) Create(sub *WebhookSubscription) (*WebhookSubscription, *Response, error) {
	u := "/webhook_subscriptions"
	v := new(WebhookSubscriptionPayload)
	p := &WebhookSubscriptionPayload{WebhookSubscription: sub}

	resp, err := s.client.newRequestDo("POST", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.WebhookSubscription, resp, nil
}

// Get gets a webhook subscription.
func (s *WebhookSubscriptionService) Get(ID string) (*WebhookSubscription, *Response, error) {
	u := fmt.Sprintf("/webhook_subscriptions/%s", ID)
	v := new(WebhookSubscriptionPayload)
	p := &WebhookSubscriptionPayload{}

	resp, err := s.client.newRequestDo("GET", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.WebhookSubscription, resp, nil
}

// Delete deletes a webhook subscription.
func (s *WebhookSubscriptionService) Delete(ID string) (*Response, error) {
	u := fmt.Sprintf("/webhook_subscriptions/%s", ID)
	return s.client.newRequestDo("DELETE", u, nil, nil, nil)
}

// Update updates a webhook subscription.
func (s *WebhookSubscriptionService) Update(ID string, sub *WebhookSubscription) (*WebhookSubscription, *Response, error) {
	u := fmt.Sprintf("/webhook_subscriptions/%s", ID)
	v := new(WebhookSubscriptionPayload)
	p := WebhookSubscriptionPayload{WebhookSubscription: sub}

	resp, err := s.client.newRequestDo("PUT", u, nil, p, v)
	if err != nil {
		return nil, nil, err
	}

	return v.WebhookSubscription, resp, nil
}
