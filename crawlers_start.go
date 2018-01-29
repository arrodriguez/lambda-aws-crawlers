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

type CrawlerStartEvent struct {
	CrawlerId string `json:"crawler_id"`
	ImageName string `json:"image_name"`
	ImageType string `json:"image_type"`
}

type CrawlerResponseEvent struct {
	CrawlerId  string `json:"crawler_id"`
	InstanceId string `json:"instance_id"`
}

var encryptedKey string
var encryptedSecretKey string
var encryptedRegion string

func CrawlersStart(event CrawlerStartEvent) (CrawlerResponseEvent, error) {
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
		return CrawlerResponseEvent{}, err
	}

	fmt.Println("#################", string(b))

	svc := ec2.New(session.New(config))

	crawlers := CrawlersGet(event, svc)

	fmt.Println("#####################", "Paso el Crawlers Get")

	networkElements := []*ec2.InstanceNetworkInterfaceSpecification{
		&ec2.InstanceNetworkInterfaceSpecification{
			AssociatePublicIpAddress: aws.Bool(true),
			DeleteOnTermination:      aws.Bool(true),
			DeviceIndex:              aws.Int64(0),
			Groups:                   []*string{aws.String("sg-7aa2ab1f")},
			SubnetId:                 aws.String("subnet-b0247cf9"),
		},
	}

	fmt.Println("#####################", "Paso el Network Elements")

	blockMapping := []*ec2.BlockDeviceMapping{
		&ec2.BlockDeviceMapping{
			DeviceName: aws.String("/dev/sda1"),
			Ebs: &ec2.EbsBlockDevice{
				DeleteOnTermination: aws.Bool(true),
				VolumeType:          aws.String("gp2"),
				VolumeSize:          aws.Int64(40),
			},
		},
	}

	fmt.Println("#####################", "Paso el Block Mapping")
	// Specify the details of the instance that you want to create.
	runResult, err := svc.RunInstances(&ec2.RunInstancesInput{
		// An Amazon Linux AMI ID for t2.micro instances in the us-west-2 region
		ImageId:      aws.String(event.ImageName),
		InstanceType: aws.String(event.ImageType),
		MinCount:     aws.Int64(1),
		MaxCount:     aws.Int64(1),
		InstanceInitiatedShutdownBehavior: aws.String("terminate"),
		NetworkInterfaces:                 networkElements,
		BlockDeviceMappings:               blockMapping,
		TagSpecifications: []*ec2.TagSpecification{
			&ec2.TagSpecification{
				ResourceType: aws.String("instance"),
				Tags: []*ec2.Tag{
					&ec2.Tag{
						Key:   aws.String("Name"),
						Value: aws.String(fmt.Sprintf("Crawler%sManaged%d", event.CrawlerId, crawlers)),
					},
				},
			},
		},
	})

	fmt.Println("#####################", "Paso el run instances")

	if err != nil {
		log.Println("Could not create intance", err)
		return CrawlerResponseEvent{}, err
	}

	log.Println("Created instance", *runResult.Instances[0].InstanceId)

	return CrawlerResponseEvent{
		CrawlerId:  event.CrawlerId,
		InstanceId: *runResult.Instances[0].InstanceId,
	}, nil

}

func CrawlersGet(event CrawlerStartEvent, svc *ec2.EC2) int {
	var instancesIterator = 0
	resp, err := svc.DescribeInstances(&ec2.DescribeInstancesInput{
		Filters: []*ec2.Filter{
			&ec2.Filter{
				Name:   aws.String("tag:Name"),
				Values: []*string{aws.String(fmt.Sprintf("Crawler%sManaged%s", event.CrawlerId, "*"))},
			},
		},
	})

	if err == nil {
		if resp.Reservations != nil {
			for _, reservation := range resp.Reservations {
				if reservation.Instances != nil {
					for _, instance := range reservation.Instances {
						if instance.State != nil {
							if instance.State.Code != nil && *(instance.State.Code) == 16 {
								// Aqui hay que extraer la ip privada y pegarle a un servicio rest para saber si tambien esta levantado el serv.
								// Si no lo esta, habria que avisar q no esta running, Delegarlo en un servicio
								instancesIterator += 1
							}
						}
					}
				}
			}
			return instancesIterator
		}
	}

	return 0
}
