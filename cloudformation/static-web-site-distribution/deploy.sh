#!/bin/bash
set -euxo pipefail

aws cloudformation create-stack \
    --stack-name static-web-site-distribution \
    --template-body template.yml \
    --parameters parameters.json \
    --capabilities CAPABILITY_NAMED_IAM
