package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

// handler is the Lambda function handler
func handler(_ context.Context) error {
	fmt.Printf("Hello, world from %s!!\n", os.Getenv("APP_ENV"))
	return nil
}

func main() {
	lambda.Start(handler)
}
