package main

import (
	"github.com/aws/aws-lambda-go/lambda"
)

func MakeRoutes() {
	lambda.Start(CrawlersStart)
}
