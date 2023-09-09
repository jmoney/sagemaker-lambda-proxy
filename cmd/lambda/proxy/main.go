package main

import (
	"context"
	"log"
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
	ilog         = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
)

func checkNonce(value string) *events.APIGatewayProxyResponse {
	if strings.Compare(value, Nonce) != 0 {
		return &events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "FORBIDDEN",
		}
	}
	return nil
}

func invokeEndpoint(mediaType, body string) *events.APIGatewayProxyResponse {
	ilog.Printf("Invoking endpoint: %s\n", EndpointName)
	ilog.Printf("Body: %s\n", body)
	response, err := Sagemaker.InvokeEndpoint(&sagemakerruntime.InvokeEndpointInput{
		EndpointName:     &EndpointName,
		ContentType:      aws.String(mediaType),
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
	nonce := event.Headers["x-nonce"]
	mediaType := event.Headers["content-type"]

	if response := checkNonce(nonce); response != nil {
		return *response, nil
	}

	if strings.Compare(mediaType, "application/json") != 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 415,
			Body:       "UNSUPPORTED_MEDIA_TYPE",
		}, nil
	}
	return *invokeEndpoint(mediaType, event.Body), nil
}

func main() {
	lambda.Start(handler)
}
