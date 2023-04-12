package pagerduty

import (
	"log"
)

// LicenseService handles the communication with license
// related methods of the PagerDuty API.
type LicenseService service

// License represents a License
type License struct {
	ID          string   `json:"id,omitempty"`
	Type        string   `json:"type,omitempty"`
	Name        string   `json:"name,omitempty"`
	Summary     string   `json:"summary,omitempty"`
	Description string   `json:"description,omitempty"`
	RoleGroup   string   `json:"role_group,omitempty"`
	ValidRoles  []string `json:"valid_roles,omitempty"`

	// The following values may be set to null or unset, so their types are
	// pointers to better translate these conditions rather than defaulting
	// to 0 or ""
	HTMLURL              *string `json:"html_url,omitempty"`
	Self                 *string `json:"self,omitempty"`
	AllocationsAvailable *int    `json:"allocations_available,omitempty"`
	CurrentValue         *int    `json:"current_value,omitempty"`
}

// LicenseAllocation represents a LicenseAllocation
type LicenseAllocation struct {
	License     *License       `json:"license,omitempty"`
	User        *UserReference `json:"user,omitempty"`
	AllocatedAt string         `json:"allocated_at,omitempty"`
}

// ListLicenseAllocationsOptions represents options when listing license_allocations.
type ListLicenseAllocationsOptions struct {
	Limit  int  `url:"limit,omitempty"`
	More   bool `url:"more,omitempty"`
	Offset int  `url:"offset,omitempty"`
	Total  int  `url:"total,omitempty"`
}

// ListLicenseAllocationsResponse represents a list response of license_allocations.
type ListLicenseAllocationsResponse struct {
	Limit              int                  `json:"limit,omitempty"`
	More               bool                 `json:"more,omitempty"`
	Offset             int                  `json:"offset,omitempty"`
	Total              int                  `json:"total,omitempty"`
	LicenseAllocations []*LicenseAllocation `json:"license_allocations,omitempty"`
}

// ListLicensesResponse represents a list response of licenses.
type ListResponse struct {
	Licenses []*License `json:"licenses,omitempty"`
}

// List lists existing Licenses.
func (s *LicenseService) List() ([]*License, *Response, error) {
	u := "/licenses"
	v := new(ListResponse)

	resp, err := s.client.newRequestDo("GET", u, nil, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v.Licenses, resp, nil
}

// ListAllocations lists existing LicenseAllocations.
func (s *LicenseService) ListAllocations(o *ListLicenseAllocationsOptions) (*ListLicenseAllocationsResponse, *Response, error) {
	u := "/license_allocations"
	v := new(ListLicenseAllocationsResponse)

	resp, err := s.client.newRequestDo("GET", u, o, nil, &v)
	if err != nil {
		return nil, nil, err
	}

	return v, resp, nil
}

// ListAllAllocations lists all existing LicenseAllocations for an Account.
func (s *LicenseService) ListAllAllocations(o *ListLicenseAllocationsOptions) ([]*LicenseAllocation, error) {
	o.More, o.Offset = true, 0
	var licenseAllocations = make([]*LicenseAllocation, 0, o.Limit)

	for o.More {
		log.Printf("==== Getting license_allocations at offset %d", o.Offset)
		v, _, err := s.ListAllocations(o)
		if err != nil {
			return licenseAllocations, err
		}
		licenseAllocations = append(licenseAllocations, v.LicenseAllocations...)
		o.More = v.More
		o.Offset = o.Offset + v.Limit
	}
	return licenseAllocations, nil
}
