
.DEFAULT_GOAL := help
help: ## Show this help  
	@grep -E '^[0-9a-zA-Z_-]+[[:blank:]]*:.*?## .*$$' $(MAKEFILE_LIST) | sort \
	| awk 'BEGIN {FS = ":.*?## "}; {printf "\033[1;32m%-15s\033[0m %s\n", $$1, $$2}'

deploy: ## Deploy CloudFormation Template
	cd static-web-site-distribution && make test-deploy || { echo "Deployment of static-web-site-distribution failed"; exit 1; }
	cd lambda-batch && make test-deploy || { echo "Deployment of lambda-batch failed"; exit 1; }
	cd lambda-with-api-gw && make test-deploy || { echo "Deployment of lambda-with-api-gw failed"; exit 1; }
	cd cloudwatch-rum && make test-deploy || { echo "Deployment of cloudwatch-rum failed"; exit 1; }
	cd daily-cost-notification && EMAIL_ADDRESS=dummy@example.com make test-deploy || { echo "Deployment of daily-cost-notification failed"; exit 1; }