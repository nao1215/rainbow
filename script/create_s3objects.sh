#!/bin/bash
set -oeu pipefail

ROOT_DIR=$(git rev-parse --show-toplevel)
SCRIPT_DIR="${ROOT_DIR}/script"
THIS_SCRIPT="${SCRIPT_DIR}/create_s3objects.sh"

export AWS_PROFILE="localstack"

echo "creating 'test-bucket-on-localstack'(s3 bucket)"
aws s3 mb s3://test-bucket-on-localstack --endpoint-url=http://localhost:4566

echo "copying 10000 files to 'test-bucket-on-localstack'(s3 bucket)"
for i in {1..10000}; do
    aws s3 cp "${THIS_SCRIPT}" "s3://test-bucket-on-localstack/object${i}.txt" --endpoint-url=http://localhost:4566 --quiet

    if [ $((i % 100)) -eq 0 ]; then
        echo "Processed $i objects"
    fi
done

echo "done"
