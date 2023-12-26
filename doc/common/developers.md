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
```shell
make docker
```

