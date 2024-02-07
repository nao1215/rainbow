.DEFAULT_GOAL := help

.PHONY: help deploy test-deploy
help:  ## Show this help message  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'

test-deploy:  ## Deploy CloudFormation Template to localstack
	@echo "Deploying S3 and CloudFront Template"
	aws cloudformation create-stack --endpoint-url "http://localhost:4566" --stack-name "static-web-site-distribution" \
    --template-body "file://template.yml" --parameters "file://parameters.json" --capabilities CAPABILITY_NAMED_IAM

deploy:  ## Deploy CloudFormation Template
	@echo "Deploying S3 and CloudFront Template"
	aws cloudformation create-stack --stack-name "static-web-site-distribution" \
    --template-body "file://template.yml" --parameters "file://parameters.json" --capabilities CAPABILITY_NAMED_IAM