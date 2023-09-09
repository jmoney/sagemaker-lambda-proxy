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
)

func main() {
	lambda.Start(handler)
}
func handler(_ context.Context, event events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	if Nonce != "" && strings.Compare(event.Headers["nonce"], Nonce) != 0 {
		return events.APIGatewayProxyResponse{
			StatusCode: 403,
			Body:       "FORBIDDEN",
		}, nil
	}

	svc := sagemakerruntime.New(session.Must(session.NewSession()))
	response, err := svc.InvokeEndpoint(&sagemakerruntime.InvokeEndpointInput{
		EndpointName:     &EndpointName,
		ContentType:      aws.String("application/json"),
		Body:             []byte(event.Body),
		CustomAttributes: aws.String("accept_eula=true"),
	})

	if err != nil {
		return events.APIGatewayProxyResponse{
			StatusCode: 500,
			Body:       err.Error(),
		}, nil
	}

	return events.APIGatewayProxyResponse{
		StatusCode: 200,
		Body:       string(response.Body),
	}, nil
}
