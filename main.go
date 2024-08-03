package main

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/sirupsen/logrus"
)

func greet(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World! %s", time.Now())
}

func run(lambdaVersion string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Lambda Version: %s", lambdaVersion)

	}
}

func main() {
	logger := logrus.New()
	lambdaVersion := os.Getenv("lambda_version")
	if lambdaVersion == "" {
		logger.Fatal("failed to get lambda_version env var")
	}

	http.HandleFunc("/", greet)
	http.HandleFunc("/run", run(lambdaVersion))

	logger.WithField("port", "8080").Info("Server listening")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		logger.WithError(err).Error("HTTP Server error")
	}
}
