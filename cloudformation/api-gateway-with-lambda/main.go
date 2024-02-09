package main

import (
	"context"
	"encoding/json"
	"log"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

// HealthResponse struct defines the response structure
type HealthResponse struct {
	Status string `json:"status"`
}

// Handler is the Lambda function handler
func Handler(_ context.Context) (events.APIGatewayProxyResponse, error) {
	// Create a response
	responseBody, err := json.Marshal(HealthResponse{Status: "healthy"})
	if err != nil {
		log.Printf("Error marshaling JSON response: %v", err)
		return events.APIGatewayProxyResponse{StatusCode: 500, Body: `{"error":"Internal Server Error"}`}, nil
	}

	// Return API Gateway response
	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       string(responseBody),
	}, nil
}

func main() {
	// Start the Lambda handler
	lambda.Start(Handler)
}
