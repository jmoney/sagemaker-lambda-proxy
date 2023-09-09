package main

import (
	"context"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemakerruntime"
)

var (
	EndpointName = os.Getenv("ENDPOINT_NAME")
	Nonce        = os.Getenv("NONCE")
	Sagemaker    = sagemakerruntime.New(session.Must(session.NewSession()))
)

func nonce(value string) *events.APIGatewayProxyResponse {
	if Nonce != "" && strings.Compare(value, Nonce) != 0 {
		return &events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "FORBIDDEN",
		}
	}
	return nil
}

func invokeEndpoint(body string) *events.APIGatewayProxyResponse {
	response, err := Sagemaker.InvokeEndpoint(&sagemakerruntime.InvokeEndpointInput{
		EndpointName:     &EndpointName,
		ContentType:      aws.String("application/json"),
		Body:             []byte(body),
		CustomAttributes: aws.String("accept_eula=true"),
	})

	if err != nil {
		return &events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}
	}

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response.Body),
	}
}

func handler(_ context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if response := nonce(event.Headers["x-nonce"]); response != nil {
		return *response, nil
	}
	return *invokeEndpoint(event.Body), nil
}

func main() {
	lambda.Start(handler)
}
