package model

import (
	"time"

	"github.com/fatih/color"
)

const (
	// CloudFormationRetryMaxAttempts is the maximum number of retries for CloudFormation.
	CloudFormationRetryMaxAttempts int = 2
	// CloudFormationWaitNanoSecTime is the time to wait for CloudFormation.
	// It is 1 hour.
	CloudFormationWaitNanoSecTime = time.Duration(6000000000000000)
)

// StackStatus is the status of a CloudFormation stack.
type StackStatus string

// String returns the string representation of StackStatus.
func (s StackStatus) String() string {
	return string(s)
}

// StringWithColor returns the string representation of StackStatus with color.
// Colors are chosen to represent the status intuitively: green for complete, red for failed, and yellow for in-progress or pending states.
func (s StackStatus) StringWithColor() string {
	switch s {
	case StackStatusCreateComplete:
		return color.GreenString("CREATE_COMPLETE")
	case StackStatusCreateFailed:
		return color.RedString("CREATE_FAILED")
	case StackStatusCreateInProgress:
		return color.YellowString("CREATE_IN_PROGRESS")
	case StackStatusDeleteComplete:
		return color.GreenString("DELETE_COMPLETE")
	case StackStatusDeleteFailed:
		return color.RedString("DELETE_FAILED")
	case StackStatusDeleteInProgress:
		return color.YellowString("DELETE_IN_PROGRESS")
	case StackStatusRollbackComplete:
		return color.GreenString("ROLLBACK_COMPLETE")
	case StackStatusRollbackFailed:
		return color.RedString("ROLLBACK_FAILED")
	case StackStatusRollbackInProgress:
		return color.YellowString("ROLLBACK_IN_PROGRESS")
	case StackStatusUpdateComplete:
		return color.GreenString("UPDATE_COMPLETE")
	case StackStatusUpdateCompleteCleanupInProgress:
		return color.YellowString("UPDATE_COMPLETE_CLEANUP_IN_PROGRESS")
	case StackStatusUpdateInProgress:
		return color.YellowString("UPDATE_IN_PROGRESS")
	case StackStatusUpdateRollbackComplete:
		return color.GreenString("UPDATE_ROLLBACK_COMPLETE")
	case StackStatusUpdateRollbackCompleteCleanupInProgress:
		return color.YellowString("UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS")
	case StackStatusUpdateRollbackFailed:
		return color.RedString("UPDATE_ROLLBACK_FAILED")
	case StackStatusUpdateFailed:
		return color.RedString("UPDATE_FAILED")
	case StackStatusUpdateRollbackInProgress:
		return color.YellowString("UPDATE_ROLLBACK_IN_PROGRESS")
	case StackStatusReviewInProgress:
		return color.YellowString("REVIEW_IN_PROGRESS")
	case StackStatusImportInProgress:
		return color.YellowString("IMPORT_IN_PROGRESS")
	case StackStatusImportComplete:
		return color.GreenString("IMPORT_COMPLETE")
	case StackStatusImportRollbackInProgress:
		return color.YellowString("IMPORT_ROLLBACK_IN_PROGRESS")
	case StackStatusImportRollbackFailed:
		return color.RedString("IMPORT_ROLLBACK_FAILED")
	case StackStatusImportRollbackComplete:
		return color.GreenString("IMPORT_ROLLBACK_COMPLETE")
	default:
		return color.RedString("UNKNOWN")
	}
}

// CloudFormation stack status constants
// Ref. https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-describing-stacks.html
const (
	// StackStatusCreateInProgress is ongoing creation of one or more stacks.
	StackStatusCreateInProgress StackStatus = "CREATE_IN_PROGRESS"
	// StackStatusCreateComplete is successful creation of one or more stacks.
	StackStatusCreateComplete StackStatus = "CREATE_COMPLETE"
	// StackStatusCreateFailed is unsuccessful creation of one or more stacks.
	// View the stack events to see any associated error messages. Possible reasons
	// for a failed creation include insufficient permissions to work with all
	// resources in the stack, parameter values rejected by an AWS service, or a
	// timeout during resource creation.
	StackStatusCreateFailed StackStatus = "CREATE_FAILED"

	// StackStatusRollbackInProgress is ongoing removal of one or more stacks after a failed
	// stack creation or after an explicitly canceled stack creation.
	StackStatusRollbackInProgress StackStatus = "ROLLBACK_IN_PROGRESS"
	// StackStatusRollbackComplete is successful removal of one or more stacks after a failed
	// stack creation or after an explicitly canceled stack creation. The stack returns to
	// the previous working state. Any resources that were created during the create stack
	// operation are deleted.
	// This status exists only after a failed stack creation. It signifies that all operations
	// from the partially created stack have been appropriately cleaned up. When in this state,
	// only a delete operation can be performed.
	StackStatusRollbackComplete StackStatus = "ROLLBACK_COMPLETE"
	// StackStatusRollbackFailed is unsuccessful removal of one or more stacks after a failed
	// stack creation or after an explicitly canceled stack creation. Delete the stack or view
	// the stack events to see any associated error messages.
	StackStatusRollbackFailed StackStatus = "ROLLBACK_FAILED"

	// StackStatusDeleteInProgress is ongoing removal of one or more stacks.
	StackStatusDeleteInProgress StackStatus = "DELETE_IN_PROGRESS"
	// StackStatusDeleteComplete is successful deletion of one or more stacks.
	// Deleted stacks are retained and viewable for 90 days.
	StackStatusDeleteComplete StackStatus = "DELETE_COMPLETE"
	// StackStatusDeleteFailed is unsuccessful deletion of one or more stacks.
	// Because the delete failed, you might have some resources that are still
	// running; however, you can't work with or update the stack. Delete the stack
	// again or view the stack events to see any associated error messages.
	StackStatusDeleteFailed StackStatus = "DELETE_FAILED"

	// StackStatusUpdateInProgress is ongoing creation of one or more stacks with
	// an expected StackId but without any templates or resources.
	// A stack with this status code counts against the maximum possible number of stacks.
	StackStatusReviewInProgress StackStatus = "REVIEW_IN_PROGRESS"
	// StackStatusUpdateInProgress is ongoing update of one or more stacks.
	StackStatusUpdateInProgress StackStatus = "UPDATE_IN_PROGRESS"
	// StackStatusUpdateCompleteCleanupInProgress is ongoing removal of old resources for
	// one or more stacks after a successful stack update. For stack updates that require
	// resources to be replaced, CloudFormation creates the new resources first and then
	// deletes the old resources to help reduce any interruptions with your stack. In this
	// state, the stack has been updated and is usable, but CloudFormation is still deleting
	// the old resources.
	StackStatusUpdateCompleteCleanupInProgress StackStatus = "UPDATE_COMPLETE_CLEANUP_IN_PROGRESS"
	// StackStatusUpdateComplete is successful update of one or more stacks.
	StackStatusUpdateComplete StackStatus = "UPDATE_COMPLETE"
	// StackStatusUpdateFailed is unsuccessful update of one or more stacks. View the stack events
	// to see any associated error messages.
	StackStatusUpdateFailed StackStatus = "UPDATE_FAILED"
	// StackStatusUpdateRollbackComplete is successful return of one or more stacks to a previous
	// working state after a failed stack update.
	StackStatusUpdateRollbackComplete StackStatus = "UPDATE_ROLLBACK_COMPLETE"
	// StackStatusUpdateRollbackCompleteCleanupInProgress is ongoing removal of new resources
	// for one or more stacks after a failed stack update. In this state, the stack has been
	// rolled back to its previous working state and is usable, but CloudFormation is still
	// deleting any new resources it created during the stack update.
	StackStatusUpdateRollbackCompleteCleanupInProgress StackStatus = "UPDATE_ROLLBACK_COMPLETE_CLEANUP_IN_PROGRESS"
	// StackStatusUpdateRollbackFailed is unsuccessful return of one or more stacks to a
	// previous working state after a failed stack update. When in this state, you can
	// delete the stack or continue rollback. You might need to fix errors before your
	// stack can return to a working state.
	StackStatusUpdateRollbackFailed StackStatus = "UPDATE_ROLLBACK_FAILED"
	// StackStatusUpdateRollbackInProgress is ongoing return of one or more stacks
	// to the previous working state after failed stack update.
	StackStatusUpdateRollbackInProgress StackStatus = "UPDATE_ROLLBACK_IN_PROGRESS"

	// StackStatusImportInProgress is the import operation is currently in progress.
	StackStatusImportInProgress StackStatus = "IMPORT_IN_PROGRESS"
	// StackStatusImportComplete is the import operation successfully completed for
	// all resources in the stack that support resource import.
	StackStatusImportComplete StackStatus = "IMPORT_COMPLETE"
	// StackStatusImportRollbackInProgress is import will roll back to the previous
	// template configuration.
	StackStatusImportRollbackInProgress StackStatus = "IMPORT_ROLLBACK_IN_PROGRESS"
	// StackStatusImportRollbackComplete is import successfully rolled back to the previous template configuration.
	StackStatusImportRollbackComplete StackStatus = "IMPORT_ROLLBACK_COMPLETE"
	// StackStatusImportRollbackFailed is the import rollback operation failed for at
	// least one resource in the stack. Results will be available for the resources
	// CloudFormation successfully imported.
	StackStatusImportRollbackFailed StackStatus = "IMPORT_ROLLBACK_FAILED"
)

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

// ResourceStatus is the status of a CloudFormation stack resource.
type ResourceStatus string

const (
	// ResourceStatusCreateInProgress is the resource is being created.
	ResourceStatusCreateInProgress ResourceStatus = "CREATE_IN_PROGRESS"
	// ResourceStatusCreateFailed is the resource creation failed.
	ResourceStatusCreateFailed ResourceStatus = "CREATE_FAILED"
	// ResourceStatusCreateComplete is the resource has been created.
	ResourceStatusCreateComplete ResourceStatus = "CREATE_COMPLETE"
	// ResourceStatusDeleteInProgress is the resource is being deleted.
	ResourceStatusDeleteInProgress ResourceStatus = "DELETE_IN_PROGRESS"
	// ResourceStatusDeleteFailed is the resource deletion failed.
	ResourceStatusDeleteFailed ResourceStatus = "DELETE_FAILED"
	// ResourceStatusDeleteComplete is the resource has been deleted.
	ResourceStatusDeleteComplete ResourceStatus = "DELETE_COMPLETE"
	// ResourceStatusDeleteSkipped is the resource was not successfully deleted. It might still be
	ResourceStatusDeleteSkipped ResourceStatus = "DELETE_SKIPPED"
	// ResourceStatusUpdateInProgress is the resource is being updated.
	ResourceStatusUpdateInProgress ResourceStatus = "UPDATE_IN_PROGRESS"
	// ResourceStatusUpdateFailed is the resource update failed.
	ResourceStatusUpdateFailed ResourceStatus = "UPDATE_FAILED"
	// ResourceStatusUpdateComplete is the resource has been updated.
	ResourceStatusUpdateComplete ResourceStatus = "UPDATE_COMPLETE"
	// ResourceStatusImportFailed is the resource import failed.
	ResourceStatusImportFailed ResourceStatus = "IMPORT_FAILED"
	// ResourceStatusImportComplete is the resource has been imported.
	ResourceStatusImportComplete ResourceStatus = "IMPORT_COMPLETE"
	// ResourceStatusImportInProgress is the resource is being imported into a stack.
	ResourceStatusImportInProgress ResourceStatus = "IMPORT_IN_PROGRESS"
	// ResourceStatusImportRollbackInProgress is the resource is being rolled back to its previous
	ResourceStatusImportRollbackInProgress ResourceStatus = "IMPORT_ROLLBACK_IN_PROGRESS"
	// ResourceStatusImportRollbackFailed is the resource import failed and the resource is
	ResourceStatusImportRollbackFailed ResourceStatus = "IMPORT_ROLLBACK_FAILED"
	// ResourceStatusImportRollbackComplete is the resource was rolled back to its previous
	ResourceStatusImportRollbackComplete ResourceStatus = "IMPORT_ROLLBACK_COMPLETE"
	// ResourceStatusUpdateRollbackInProgress is the resource is being rolled back as part of a
	ResourceStatusUpdateRollbackInProgress ResourceStatus = "UPDATE_ROLLBACK_IN_PROGRESS"
	// ResourceStatusUpdateRollbackComplete is the resource was rolled back to its previous
	ResourceStatusUpdateRollbackComplete ResourceStatus = "UPDATE_ROLLBACK_COMPLETE"
	// ResourceStatusUpdateRollbackFailed is the resource update failed and the resource is being rolled back to its previous configuration.
	ResourceStatusUpdateRollbackFailed ResourceStatus = "UPDATE_ROLLBACK_FAILED"
	// ResourceStatusRollbackInProgress is the resource is being rolled back.
	ResourceStatusRollbackInProgress ResourceStatus = "ROLLBACK_IN_PROGRESS"
	// ResourceStatusRollbackComplete is the resource was rolled back.
	ResourceStatusRollbackComplete ResourceStatus = "ROLLBACK_COMPLETE"
	// ResourceStatusRollbackFailed is the resource rollback failed.
	ResourceStatusRollbackFailed ResourceStatus = "ROLLBACK_FAILED"
)

// Values returns all known values for ResourceStatus. Note that this can be
// expanded in the future, and so it is only as up to date as the client. The
// ordering of this slice is not guaranteed to be stable across updates.
func (ResourceStatus) Values() []ResourceStatus {
	return []ResourceStatus{
		ResourceStatusCreateInProgress,
		ResourceStatusCreateFailed,
		ResourceStatusCreateComplete,
		ResourceStatusDeleteInProgress,
		ResourceStatusDeleteFailed,
		ResourceStatusDeleteComplete,
		ResourceStatusDeleteSkipped,
		ResourceStatusUpdateInProgress,
		ResourceStatusUpdateFailed,
		ResourceStatusUpdateComplete,
		ResourceStatusImportFailed,
		ResourceStatusImportComplete,
		ResourceStatusImportInProgress,
		ResourceStatusImportRollbackInProgress,
		ResourceStatusImportRollbackFailed,
		ResourceStatusImportRollbackComplete,
		ResourceStatusUpdateRollbackInProgress,
		ResourceStatusUpdateRollbackComplete,
		ResourceStatusUpdateRollbackFailed,
		ResourceStatusRollbackInProgress,
		ResourceStatusRollbackComplete,
		ResourceStatusRollbackFailed,
	}
}

// StackResourceDriftStatus is status of the resource's actual configuration
type StackResourceDriftStatus string

const (
	// StackResourceDriftStatusInSync is the resource's actual configuration matches its expected
	StackResourceDriftStatusInSync StackResourceDriftStatus = "IN_SYNC"
	// StackResourceDriftStatusModified is the resource differs from its expected configuration.
	StackResourceDriftStatusModified StackResourceDriftStatus = "MODIFIED"
	// StackResourceDriftStatusDeleted is the resource differs from its expected configuration in that it has been deleted.
	StackResourceDriftStatusDeleted StackResourceDriftStatus = "DELETED"
	// StackResourceDriftStatusNotChecked is CloudFormation hasn't checked if the resource differs from its expected configuration.
	StackResourceDriftStatusNotChecked StackResourceDriftStatus = "NOT_CHECKED"
)

// Values returns all known values for StackResourceDriftStatus. Note that this
// can be expanded in the future, and so it is only as up to date as the client.
// The ordering of this slice is not guaranteed to be stable across updates.
func (StackResourceDriftStatus) Values() []StackResourceDriftStatus {
	return []StackResourceDriftStatus{
		StackResourceDriftStatusInSync,
		StackResourceDriftStatusModified,
		StackResourceDriftStatusDeleted,
		StackResourceDriftStatusNotChecked,
	}
}

// StackResourceDriftInformationSummary is summarizes information about whether
// the resource's actual configuration differs, or has drifted, from its expected configuration.
type StackResourceDriftInformationSummary struct {
	// StackResourceDriftStatus is status of the resource's actual
	// configuration compared to its expected configuration.
	StackResourceDriftStatus StackResourceDriftStatus
	// LastCheckTimestamp is when CloudFormation last checked if the
	// resource had drifted from its expected configuration.
	LastCheckTimestamp *time.Time
}

// ModuleInfo is contains information about the module from which the resource
// was created, if the resource was created from a module included in the stack
// template. For more information about modules, see Using modules to encapsulate
// and reuse resource configurations (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/modules.html)
// in the CloudFormation User Guide.
type ModuleInfo struct {
	// LogicalIDHierarchy is a concatenated list of the logical IDs of the module
	// or modules containing the resource. Modules are listed starting with the
	// inner-most nested module, and separated by / .
	// In the following example, the resource was created from a module, moduleA,
	// that's nested inside a parent module, moduleB . moduleA/moduleB For more
	// information, see Referencing resources in a module
	// (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/modules.html#module-ref-resources)
	// in the CloudFormation User Guide.
	LogicalIDHierarchy *string
	// TypeHierarchy is a concatenated list of the module type or types containing
	// the resource. Module types are listed starting with the inner-most nested
	// module, and separated by /.
	// In the following example, the resource was created from a module of type
	// AWS::First::Example::MODULE , that's nested inside a parent module of type
	// AWS::Second::Example::MODULE .
	// AWS::First::Example::MODULE/AWS::Second::Example::MODULE
	TypeHierarchy *string
}

// StackResource is a CloudFormation stack resource. It is same as types.StackResourceSummary.
type StackResource struct {
	// LastUpdatedTimestamp is time the status was updated.
	LastUpdatedTimestamp *time.Time
	// LogicalResourceID is the logical name of the resource specified in the template.
	LogicalResourceID *string
	// ResourceStatus is current status of the resource.
	ResourceStatus ResourceStatus
	// ResourceType is type of resource. For more information, go to Amazon Web Services Resource
	// Types Reference (https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/aws-template-resource-type-ref.html)
	// in the CloudFormation User Guide.
	ResourceType *string
	// DriftInformation is information about whether the resource's actual
	// configuration differs, or has drifted, from its expected configuration,
	// as defined in the stack template and any values specified as template
	// parameters. For more information, see Detecting Unregulated Configuration
	// Changes to Stacks and Resources
	// https://docs.aws.amazon.com/AWSCloudFormation/latest/UserGuide/using-cfn-stack-drift.html
	DriftInformation *StackResourceDriftInformationSummary
	// ModuleInfo is contains information about the module from which the resource was created, if
	// the resource was created from a module included in the stack template.
	ModuleInfo *ModuleInfo
	// PhysicalResourceID is the name or unique identifier that corresponds to a
	// physical instance ID of the resource.
	PhysicalResourceID *string
	// ResourceStatusReason is success/failure message associated with the resource.
	ResourceStatusReason *string
}
