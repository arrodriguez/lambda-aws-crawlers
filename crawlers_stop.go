package main

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/spf13/viper"
)

type CrawlerStopEvent struct {
	InstanceId string `json:"instance_id"`
}

var encryptedKey string
var encryptedSecretKey string
var encryptedRegion string

func CrawlersStop(event CrawlerStopEvent) error {
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

	b, err := json.Marshal(event)

	if err != nil {
		return err
	}

	fmt.Println("#################", string(b))

	svc := ec2.New(session.New(config))

	fmt.Println("#####################", "Paso el Block Mapping")
	// Specify the details of the instance that you want to create.
	_, err = svc.StopInstances(&ec2.StopInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		Force: aws.Bool(true),
		InstanceIds: []*string{
			aws.String(event.InstanceId),
		},
	})

	fmt.Println("#####################", "Paso el stop instances")

	if err != nil {
		log.Println("Could not stop intance", err)
		return err
	}

	log.Println("Stopped instance", event)

	return nil
}
