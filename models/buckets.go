package models

import "time"

type Bucket struct {
	// if true, object versioning will be enabled for the bucket.
	EnableVersioning *bool `json:"enableVersioning,omitempty"`
	// Bucket name. Must be unique across the grid and DNS compliant. See the instructions for using S3 for details.
	Name string `json:"name"`
	// the region for this bucket, which must already be defined (defaults to us-east-1 if not specified)
	Region       string                      `json:"region,omitempty"`
	S3ObjectLock *BucketS3ObjectLockSettings `json:"s3ObjectLock,omitempty"`
	// the creation time of the bucket
	CreationTime time.Time `json:"creationTime,omitempty"`
	// compliance settings for the bucket
	Compliance *BucketComplianceSettings `json:"compliance,omitempty"`
	// status of the object deletion
	DeleteObjectStatus *BucketDeleteObjectStatus `json:"deleteObjectStatus,omitempty"`
}

// BucketS3ObjectLockSettings Settings for S3 Object Lock. Cannot be used with legacy Compliance.
type BucketS3ObjectLockSettings struct {
	// whether the bucket has S3 Object Lock enabled
	Enabled                 *bool                                       `json:"enabled"`
	DefaultRetentionSetting *BucketS3ObjectLockDefaultRetentionSettings `json:"defaultRetentionSetting,omitempty"`
}

// BucketS3ObjectLockDefaultRetentionSettings Default retention settings for S3 Object Lock.
type BucketS3ObjectLockDefaultRetentionSettings struct {
	// The retention mode used for new objects added to this bucket. Must be compliance, which means that an object version cannot be overwritten or deleted by any user.
	Mode string `json:"mode"`
	// The length of the default retention period for new objects added to this bucket, in days. If provided, must be paired with retentionMode. Does not affect existing bucket objects or objects with their own retain-until-date settings.
	Days int32 `json:"days,omitempty"`
	// The length of the default retention period for new objects added to this bucket, in years. If provided, must be paired with retentionMode. Does not affect existing bucket objects or objects with their own retain-until-date settings.
	Years int32 `json:"years,omitempty"`
}

type BucketComplianceSettings struct {
	// Wether the objects are autoDeleted
	AutoDelete *bool `json:"autoDelete"`
	// time to legally hold the objects
	LegalHold *bool `json:"legalHold"`
	// amount of minuts for retention
	RetentionPeriodMinutes *int32 `json:"retentionPeriodMinutes,omitempty"`
}

type BucketDeleteObjectStatus struct {
	// are the objects being deleted
	IsDeletingObjects *bool `json:"isDeletingObjects"`
	// initial Object count before operation
	InitialObjectCount *int32 `json:"initialObjectCount,omitempty"`
	// initial Object Bytes before the operation
	InitialObjectBytes *int64 `json:"initialObjectBytes,omitempty"`
}

// BucketUsage represents the per-bucket usage metrics returned by
// GET /org/containers/{bucketName}/usage.
type BucketUsage struct {
	// number of objects in the bucket
	ObjectCount *int64 `json:"objectCount,omitempty"`
	// logical size in bytes of all objects in the bucket
	DataBytes *int64 `json:"dataBytes,omitempty"`
}

// BucketCorsConfiguration wraps the CORS XML configuration for an S3 bucket.
// A nil Cors field disables CORS on the bucket.
type BucketCorsConfiguration struct {
	// XML for configuring CORS, or null to disable CORS
	Cors *string `json:"cors"`
}

// BucketNotificationConfiguration wraps the notification XML configuration for an S3 bucket.
// A nil Notification field disables notifications on the bucket.
type BucketNotificationConfiguration struct {
	// notification configuration XML, or null to disable notifications
	Notification *string `json:"notification"`
}

// BucketPolicyConfiguration wraps the bucket policy document for an S3 bucket.
// A nil Policy field disables the bucket policy.
type BucketPolicyConfiguration struct {
	Policy *BucketPolicy `json:"policy"`
}

// BucketPolicy represents an S3 bucket policy document.
type BucketPolicy struct {
	// Optional policy identifier.
	Id string `json:"Id,omitempty"`
	// Policy language version (e.g. "2012-10-17" or "2015-09-08").
	Version string `json:"Version,omitempty"`
	// One or more policy statements.
	Statement []BucketPolicyStatement `json:"Statement"`
}

// BucketPolicyStatement represents a single statement within a BucketPolicy.
//
// Several fields use interface{} because the S3 policy language allows either a
// single string or a list of strings (and Principal/NotPrincipal additionally
// allow either the string "*" or an object such as {"AWS": "..."}). Callers may
// pass a string, []string, map[string]interface{}, or json.RawMessage as
// appropriate.
type BucketPolicyStatement struct {
	// Optional statement identifier.
	Sid string `json:"Sid,omitempty"`
	// "Allow" or "Deny".
	Effect string `json:"Effect"`
	// String or []string. Mutually exclusive with NotAction.
	Action interface{} `json:"Action,omitempty"`
	// String or []string. Mutually exclusive with Action.
	NotAction interface{} `json:"NotAction,omitempty"`
	// String or []string. Mutually exclusive with NotResource.
	Resource interface{} `json:"Resource,omitempty"`
	// String or []string. Mutually exclusive with Resource.
	NotResource interface{} `json:"NotResource,omitempty"`
	// Condition block keyed by condition type, then condition key.
	Condition map[string]map[string]interface{} `json:"Condition,omitempty"`
	// "*" for anonymous, or an object such as {"AWS": "<account-or-arn>"}.
	// Mutually exclusive with NotPrincipal.
	Principal interface{} `json:"Principal,omitempty"`
	// "*" or an object. Mutually exclusive with Principal.
	NotPrincipal interface{} `json:"NotPrincipal,omitempty"`
}
