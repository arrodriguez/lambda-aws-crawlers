package main

import (
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/spf13/viper"
)

func init() {
	viper.AutomaticEnv()
	viper.SetEnvPrefix("crawlers")
}

func main() {
	lambda.Start(CrawlersStop)
}
