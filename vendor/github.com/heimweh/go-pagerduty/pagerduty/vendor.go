package pagerduty

import (
	"fmt"
)

// VendorService handles the communication with vendor
// related methods of the PagerDuty API.
type VendorService service

// Vendor represents a vendor.
type Vendor struct {
	Description         string `json:"description,omitempty"`
	GenericServiceType  string `json:"generic_service_type,omitempty"`
	HTMLURL             string `json:"html_url,omitempty"`
	ID                  string `json:"id,omitempty"`
	IntegrationGuideURL string `json:"integration_guide_url,omitempty"`
	LogoURL             string `json:"logo_url,omitempty"`
	Name                string `json:"name,omitempty"`
	Self                string `json:"self,omitempty"`
	Summary             string `json:"summary,omitempty"`
	ThumbnailURL        string `json:"thumbnail_url,omitempty"`
	Type                string `json:"type,omitempty"`
	WebsiteURL          string `json:"website_url,omitempty"`
}

// ListVendorsOptions represents options when listing vendors.
type ListVendorsOptions struct {
	Limit  int    `url:"limit,omitempty"`
	More   bool   `url:"more,omitempty"`
	Offset int    `url:"offset,omitempty"`
	Total  int    `url:"total,omitempty"`
	Query  string `url:"query,omitempty"`
}

// ListVendorsResponse represents a list response of vendors.
type ListVendorsResponse struct {
	Limit   int       `json:"limit,omitempty"`
	More    bool      `json:"more,omitempty"`
	Offset  int       `json:"offset,omitempty"`
	Total   int       `json:"total,omitempty"`
	Vendors []*Vendor `json:"vendors,omitempty"`
}

// VendorPayload represents a vendor.
type VendorPayload struct {
	Vendor *Vendor `json:"vendor,omitempty"`
}

// List lists existing vendors.
func (s *VendorService) List(o *ListVendorsOptions) (*ListVendorsResponse, *Response, error) {
	u := "/vendors"
	v := new(ListVendorsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// Get retrieves information about a vendor.
func (s *VendorService) Get(id string) (*Vendor, *Response, error) {
	u := fmt.Sprintf("/vendors/%s", id)
	v := new(VendorPayload)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Vendor, resp, nil
}
