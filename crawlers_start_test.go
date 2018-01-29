package main

import (
	"log"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
	"github.com/stretchr/testify/assert"
)

func TestCrawlersGet(t *testing.T) {
	config := aws.NewConfig()
	config.WithCredentialsChainVerboseErrors(true)

	log.Println("##############", viper.GetString("AWS_ACCESS_KEY_ID"))
	log.Println("##############", viper.GetString("AWS_SECRET_ACCESS_KEY"))
	log.Println("##############", viper.GetString("AWS_REGION"))

	config.WithCredentials(
		credentials.NewStaticCredentials(
			viper.GetString("AWS_ACCESS_KEY_ID"),
			viper.GetString("AWS_SECRET_ACCESS_KEY"),
			"",
		),
	)

	config.WithRegion(viper.GetString("AWS_REGION"))

	svc := ec2.New(session.New(config))
	crawlers := CrawlersGet(CrawlerStartEvent{CrawlerId: "1912"}, svc)

	assert.Equal(t, 1, crawlers)
}
