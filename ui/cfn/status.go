package cfn

// status is the status of the cfn operation.
type status uint

const (
	// statusNone is the status when the cfn operation is not executed.
	statusNone status = iota
	// statusRegionSelecting is the status when the cfn operation is executed and the region is being selected.
	statusRegionSelecting
	// statusStacksFetching is the status when the cfn operation is executed and the stacks are being fetched.
	statusStacksFetching
	// statusStacksFetched is the status when the cfn operation is executed and the stacks are fetched.
	statusStacksFetched
	// statusStacksListed is the status when the cfn operation is executed and the stacks are listed.
	statusStacksListed
	// statusReturnToTop is the status when the cfn operation is executed and the user wants to return to the top.
	statusReturnToTop
	// statusQuit is the status when the cfn operation is executed and the user wants to quit.
	statusQuit
)
