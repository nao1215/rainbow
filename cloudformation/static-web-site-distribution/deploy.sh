#!/bin/bash
set -euxo pipefail

aws cloudformation create-stack \
    --endpoint-url "http://localhost:4566" \
    --stack-name "static-web-site-distribution" \
    --template-body "file://template.yml" \
    --parameters "file://parameters.json" \
    --capabilities CAPABILITY_NAMED_IAM
