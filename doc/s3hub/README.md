## s3hub - user-friendly s3 management tool
> [!IMPORTANT]  
> Not implemented yet.

The s3hub command provides following features:
- Create a bucket
- List buckets
- List contents of a bucket
- Copy files to a bucket
- Delete contents from a bucket
- Delete a bucket
  
## How to install
```shell
go install github.com/nao1215/rainbow/cmd/s3hub@latest
```

## How to use
S3hub operates without requiring the 's3://' protocol to be added to the bucket name.

### Create a bucket(s)

```shell
s3hub mb ${YOUR_BUCKET_NAME}
```

### List buckets
```shell
s3hub ls
```

### List contents of a bucket
```shell
s3hub ls ${YOUR_BUCKET_NAME}
```

### Copy files to a bucket
From local to S3:
```shell
s3hub cp ${YOUR_FILE_PATH} ${YOUR_BUCKET_NAME}
```

From S3 to local:
```shell
s3hub cp ${YOUR_BUCKET_NAME} ${YOUR_FILE_PATH}
```

### Delete contents from a bucket
If you want to delete a specific file(s), use the following command:
```shell
s3hub rm ${CONTENT_PATH_IN_BUCKET}
```

If you want to delete all contents in a bucket, use the wildcard:
```shell
s3hub rm ${YOUR_BUCKET_NAME}/*
```

### Delete a bucket(s)
```shell
s3hub rm --recursive ${YOUR_BUCKET_NAME}
```

### Interactive mode
You can use the interactive mode by omitting the arguments.
```shell
s3hub
```