package pagerduty

type resourceReference struct {
	HTMLURL string `json:"html_url,omitempty"`
	ID      string `json:"id,omitempty"`
	Self    string `json:"self,omitempty"`
	Summary string `json:"summary,omitempty"`
	Type    string `json:"type,omitempty"`
}

// UserReferenceWrapper is a wrapper around UserReference
type UserReferenceWrapper struct {
	User *UserReference `json:"user,omitempty"`
}

// UserReference represents a reference to a user.
type UserReference resourceReference

// EscalationPolicyReference represents a reference to an escalation policy.
type EscalationPolicyReference resourceReference

// ResponsePlayReference represents a reference to a response play.
type ResponsePlayReference resourceReference

// ScheduleReference represents a reference to a schedule.
type ScheduleReference resourceReference

// TeamReference represents a reference to a team.
type TeamReference resourceReference

// LicenseReference represents a reference to a team.
type LicenseReference resourceReference

// ContactMethodReference represents a reference to a contact method.
type ContactMethodReference resourceReference

// AddonReference represents a reference to an add-on.
type AddonReference resourceReference

// ServiceReference represents a reference to a service.
type ServiceReference resourceReference

// IntegrationReference represents a reference to an integration.
type IntegrationReference resourceReference

// EscalationTargetReference represents a reference to an escalation target
type EscalationTargetReference resourceReference

// VendorReference represents a reference to a vendor
type VendorReference resourceReference

// ExtensionSchemaReference represents a reference to an extension schema
type ExtensionSchemaReference resourceReference

// RulesetReference represents a reference to a ruleset schema
type RulesetReference resourceReference

// SubscriberReference represents a reference to a subscriber schema
type SubscriberReference resourceReference

// IncidentAttributeReference represents a reference to a Incident
// Attribute schema
type IncidentAttributeReference resourceReference

// CustomFieldSchemaReference represents a reference to a Custom
// Field schema
type CustomFieldSchemaReference resourceReference
