## Development Eviorment Setup
### Install Prerequisites
- [Install Golang](https://go.dev/doc/install)
- [Install make](https://www.gnu.org/software/make/)
- [Install Docker](https://docs.docker.com/get-docker/)
- [Install Docker Compose](https://docs.docker.com/compose/install/)
- Install development tools by running the following command:
    ```shell
    make tools
    ```

### Makefile Usage
If you want to print help information, you can run the following command:
```shell
$ make
build           Build binary
changelog       Generate changelog
clean           Clean project
docker          Start docker (localstack)
generate        Generate code from templates
test            Start unit test
tools           Install dependency tools 
```

### localstack Setup
The localstack is used to simulate the AWS environment. First, you configure a custom profile to use with LocalStack. Add the following profile to your AWS configuration file (by default, this file is at ~/.aws/config):

```shell
[profile localstack]
region=us-east-1
output=json
endpoint_url = http://localhost:4566
```

Add the following profile to your AWS credentials file (by default, this file is at ~/.aws/credentials):
```shell
[localstack]
aws_access_key_id=test
aws_secret_access_key=test
```

> [!NOTE]  
> Alternatively, you can also set the AWS_PROFILE=localstack environment variable, in which case the --profile localstack parameter can be omitted in the commands above.

### Run localstack
Run the following command to start localstack:
```shell
make docker
```

Check the status of localstack:
```shell
curl -s "http://127.0.0.1:4566/health" | jq .
{
  "services": {
    "acm": "available",
    "apigateway": "available",
    "cloudformation": "available",
    "cloudwatch": "available",
    "config": "available",
    "dynamodb": "available",
    "dynamodbstreams": "available",
    "ec2": "available",
    "es": "available",
    "events": "available",
    "firehose": "available",
    "iam": "available",
    "kinesis": "available",
    "kms": "available",
    "lambda": "available",
    "logs": "available",
    "opensearch": "available",
    "redshift": "available",
    "resource-groups": "available",
    "resourcegroupstaggingapi": "available",
    "route53": "available",
    "route53resolver": "available",
    "s3": "available",
    "s3control": "available",
    "secretsmanager": "available",
    "ses": "available",
    "sns": "available",
    "sqs": "available",
    "ssm": "available",
    "stepfunctions": "available",
    "sts": "available",
    "support": "available",
    "swf": "available",
    "transcribe": "available"
  },
  "version": "2.1.1.dev"
}
```

### Generate CHANGELOG
```shell
make changelog
```

### Generate coverage map
```shell
make coverage-tree
```

![coverage-image](cover.svg)