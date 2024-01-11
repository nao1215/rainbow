package model

import "time"

const (
	// CloudFormationRetryMaxAttempts is the maximum number of retries for CloudFormation.
	CloudFormationRetryMaxAttempts int = 2
)

// StackStatus is the status of a CloudFormation stack.
type StackStatus string

// StackDriftInformationSummary contains information about whether the stack's
// actual configuration differs, or has drifted, from its expected configuration,
// as defined in the stack template and any values specified as template parameters.
// A stack is considered to have drifted if one or more of its resources have drifted.
type StackDriftInformationSummary struct {
	// StackDriftStatus is status of the stack's actual configuration compared to its expected template
	// configuration.
	StackDriftStatus StackDriftStatus
	// LastCheckTimestamp is most recent time when a drift detection operation was
	// initiated on the stack, or any of its individual resources that support drift detection.
	LastCheckTimestamp *time.Time
}

// StackDriftStatus is the status of a stack's actual configuration compared to
// its expected template configuration.
type StackDriftStatus string

const (
	// StackDriftStatusDrifted is the stack differs from its expected template configuration.
	// A stack is considered to have drifted if one or more of its resources have drifted.
	StackDriftStatusDrifted StackDriftStatus = "DRIFTED"
	// StackDriftStatusInSync is the stack's actual configuration matches its expected template
	// configuration.
	StackDriftStatusInSync StackDriftStatus = "IN_SYNC"
	// StackDriftStatusNotChecked is CloudFormation hasn't checked if the stack differs from its
	// expected template configuration.
	StackDriftStatusNotChecked StackDriftStatus = "NOT_CHECKED"
	// StackDriftStatusUnknown is this value is reserved for future use.
	StackDriftStatusUnknown StackDriftStatus = "UNKNOWN"
)

// Values returns all known values for StackDriftStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (StackDriftStatus) Values() []StackDriftStatus {
	return []StackDriftStatus{
		"DRIFTED",
		"IN_SYNC",
		"UNKNOWN",
		"NOT_CHECKED",
	}
}

// Stack is a CloudFormation stack. It is same as types.StackSummary.
type Stack struct {
	// CreationTime is the time the stack was created.
	CreationTime *time.Time
	// StackName is the name associated with the stack.
	StackName *string
	// StackStatus is the current status of the stack.
	StackStatus StackStatus
	// DeletionTime is the time the stack was deleted.
	DeletionTime *time.Time
	// DriftInformation is summarizes information about whether a stack's actual
	// configuration differs, or has drifted, from its expected configuration,
	// as defined in the stack template and any values specified as template parameters.
	DriftInformation *StackDriftInformationSummary
	// LastUpdatedTime is the time the stack was last updated.
	LastUpdatedTime *time.Time
	// ParentID is used for nested stacks --stacks created as resources for
	// another stack-- the stack ID of the direct parent of this stack.
	// For the first level of nested stacks, the root stack is also the parent stack.
	ParentID *string
	// RootID id used for nested stacks --stacks created as resources for
	// another stack--the stack ID of the top-level stack to which the nested stack
	// ultimately belongs.
	RootID *string
	// StackID is unique stack identifier.
	StackID *string
	// StackStatusReason is Success/Failure message associated with the stack status.
	StackStatusReason *string
	// TemplateDescription is the template description of the template used to create the stack.
	TemplateDescription *string
}
