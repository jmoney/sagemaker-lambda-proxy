package main

import (
	"flag"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sagemaker"
)

var (
	EndpointName       = flag.String("endpoint-name", "", "SageMaker endpoint name")
	EndpointConfigName = flag.String("endpoint-config", "", "SageMaker endpoint config name")
	Action             = flag.String("action", "", "Action to perform.  Options are delete or create")
	ilog               = log.New(os.Stdout, "[INFO] ", log.Ldate|log.Ltime|log.Lshortfile)
	wlog               = log.New(os.Stdout, "[WARN] ", log.Ldate|log.Ltime|log.Lshortfile)
	elog               = log.New(os.Stderr, "[ERROR] ", log.Ldate|log.Ltime|log.Lshortfile)
)

func main() {
	flag.Parse()

	switch *Action {
	case "delete":
		deleteEndpoint()
	case "create":
		createEndpoint()
	default:
		panic("Invalid action")
	}
}

func deleteEndpoint() {
	if *EndpointName == "" {
		panic("Endpoint name is required")
	}

	svc := sagemaker.New(session.Must(session.NewSession()))
	response, err := svc.DescribeEndpoint(&sagemaker.DescribeEndpointInput{
		EndpointName: EndpointName,
	})

	if err != nil {
		elog.Printf("%s\n", err.Error())
		return
	}

	ilog.Printf("Endpoint status: %s\n", *response.EndpointStatus)

	svc.DeleteEndpoint(&sagemaker.DeleteEndpointInput{
		EndpointName: EndpointName,
	})
}

func createEndpoint() {
	if *EndpointName == "" {
		panic("Endpoint name is required")
	}

	if *EndpointConfigName == "" {
		panic("Endpoint config name is required")
	}

	svc := sagemaker.New(session.Must(session.NewSession()))
	resp, err := svc.DescribeEndpoint(&sagemaker.DescribeEndpointInput{
		EndpointName: EndpointName,
	})

	if err == nil && *resp.EndpointStatus == "InService" {
		ilog.Printf("Endpoint %s already exists and is in service\n", *resp.EndpointName)
		return
	} else if err != nil {
		svc.CreateEndpoint(&sagemaker.CreateEndpointInput{
			EndpointName:       EndpointName,
			EndpointConfigName: EndpointConfigName,
		})
	}

	for {
		resp, err := svc.DescribeEndpoint(&sagemaker.DescribeEndpointInput{
			EndpointName: EndpointName,
		})

		if err != nil {
			wlog.Printf("Error describing endpoint: %s\n", err.Error())
		}

		ilog.Printf("Endpoint Name %s, status: %s\n", *resp.EndpointName, *resp.EndpointStatus)

		if *resp.EndpointStatus == "InService" {
			break
		}

		time.Sleep(2 * time.Second)
	}
}
