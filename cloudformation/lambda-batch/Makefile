.DEFAULT_GOAL := help

BINARY_NAME = bootstrap
AWS_STACK_NAME = lambda-batch
TEMPLATE_FILE = template.yml

.PHONY: help clean dependency_check build deploy
help:  ## Show this help message  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'

clean: ## Clean project
	-rm -rf $(BINARY_NAME)

dependency_check:  ## Check dependencies
	@command -v sam > /dev/null 2>&1 || { echo "Error: sam is not installed. See https://github.com/aws/aws-sam-cli"; exit 1; }

build:  ## Build binary
	GOOS=linux GOARCH=arm64 go build -o $(BINARY_NAME) main.go

deploy: dependency_check build ## Deploy CloudFormation Template
	sam deploy --stack-name $(AWS_STACK_NAME) --template-file $(TEMPLATE_FILE) --capabilities CAPABILITY_IAM --resolve-s3 --region ap-northeast-1