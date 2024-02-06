package main

import (
	"context"
	"fmt"
	"os"

	"github.com/aws/aws-lambda-go/lambda"
)

// handler is the Lambda function handler
func handler(_ context.Context) error {
	if env, exists := os.LookupEnv("APP_ENV"); exists {
		fmt.Printf("Hello, world from %s!!\n", env)
	} else {
		fmt.Println("Hello, world!!")
	}
	return nil
}

func main() {
	lambda.Start(handler)
}
