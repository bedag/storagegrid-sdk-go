package models

// S3ObjectLock represents the grid-wide S3 Object Lock (compliance-global) settings.
// Fields are pointers so callers can submit partial PUT bodies and so GET responses
// with omitted fields decode cleanly.
type S3ObjectLock struct {
	ComplianceEnabled             *bool `json:"complianceEnabled,omitempty"`             // Whether S3 Object Lock is enabled grid-wide.
	LegacyComplianceEnabled       *bool `json:"legacyComplianceEnabled,omitempty"`       // Whether the deprecated legacy Compliance feature is enabled.
	CreateLegacyComplianceBuckets *bool `json:"createLegacyComplianceBuckets,omitempty"` // Whether tenants are allowed to create new legacy Compliance buckets.
}
