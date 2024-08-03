package main

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/lambda"
	"github.com/sirupsen/logrus"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func run(logger *logrus.Logger, lambdaClient *lambda.Client, lambdaFuncName string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		ctx := r.Context()
		fmt.Fprintf(w, "Lambda func name: %s", lambdaFuncName)
		resp, err := lambdaClient.Invoke(ctx, &lambda.InvokeInput{
			FunctionName: &lambdaFuncName,
		})
		if err != nil {
			logger.WithError(err).Error("Failed to Invoke lambda")
			fmt.Fprintf(w, "failed: %s", err.Error())
			return
		}

		fmt.Fprintf(w, "response payload: %v", resp.Payload)
	}
}

func main() {
	logger := logrus.New()
	ctx := context.Background()

	awsConfig, err := config.LoadDefaultConfig(ctx)
	if err != nil {
		logger.WithError(err).Fatal("Loading default AWS config")
	}

	lambdaFuncName := os.Getenv("LAMBDA_FUNC_NAME")
	if lambdaFuncName == "" {
		logger.Fatal("lambda func name is empty")
	}

	lambdaClient := lambda.NewFromConfig(awsConfig)

	http.HandleFunc("/", greet)
	http.HandleFunc("/run", run(logger, lambdaClient, lambdaFuncName))

	logger.WithField("port", "8080").Info("Server listening")
	err = http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.WithError(err).Error("HTTP Server error")
	}
}
