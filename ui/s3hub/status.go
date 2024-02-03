package s3hub

// status is the status of the s3hub operation.
type status uint

const (
	// statusNone is the status when the s3hub operation is not executed.
	statusNone status = iota
	// statusBucketFetching is the status when the s3hub operation is executed and the bucket is being fetched.
	statusBucketFetching
	// statusBucketFetched is the status when the s3hub operation is executed and the bucket is fetched.
	statusBucketFetched
	// statusS3ObjectFetching is the status when the s3hub operation is executed and the S3 object is being fetched.
	statusS3ObjectFetching
	// statusS3ObjectFetched is the status when the s3hub operation is executed and the S3 object is fetched.
	statusS3ObjectFetched
	// statusBucketListed is the status when the s3hub operation is executed and the bucket is listed.
	statusBucketListed
	// statusS3ObjectListed is the status when the s3hub operation is executed and the S3 object is listed.
	statusS3ObjectListed
	// statusBucketCreating is the status when the s3hub operation is executed and the bucket is being created.
	statusBucketCreating
	// statusBucketCreated is the status when the s3hub operation is executed and the bucket is created.
	statusBucketCreated
	// statusDownloading is the status when the s3hub operation is executed and the object is being downloaded.
	statusDownloading
	// statusDownloaded is the status when the s3hub operation is executed and the object is downloaded.
	statusDownloaded
	// statusBucketDeleting is the status when the s3hub operation is executed and the bucket is being deleted.
	statusBucketDeleting
	// statusBucketDeleted is the status when the s3hub operation is executed and the bucket is deleted.
	statusBucketDeleted
	// statusS3ObjectDeleting is the status when the s3hub operation is executed and the S3 object is being deleted.
	statusS3ObjectDeleting
	// statusS3ObjectDeleted is the status when the s3hub operation is executed and the S3 object is deleted.
	statusS3ObjectDeleted
	// statusReturnToTop is the status when the s3hub operation is executed and the user wants to return to the top.
	statusReturnToTop
	// statusQuit is the status when the s3hub operation is executed and the user wants to quit.
	statusQuit
)
