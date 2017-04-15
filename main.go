package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/route53"
	"github.com/jawher/mow.cli"
	"os"
	"time"
)

const (
	APP_NAME        = "dyndns53"
	APP_DESCRIPTION = ""
)

var (
	gitHash string
)

func main() {
	app := cli.App(APP_NAME, APP_DESCRIPTION)
	pollInterval := *app.Int(cli.IntOpt{
		Name:   "poll-interval-seconds",
		Desc:   "Time in seconds between each poll",
		EnvVar: "POLL_INTERVAL_SECONDS",
	})
	ipServiceUrl := *app.String(cli.StringOpt{
		Name:   "ip-service-url",
		Desc:   "URL of service for evaluating current address",
		EnvVar: "IP_SERVICE_URL",
	})

	awsRegion := *app.String(cli.StringOpt{
		Name:   "aws-region",
		Desc:   "AWS region",
		EnvVar: "AWS_REGION",
	})
	awsAccessKey := *app.String(cli.StringOpt{
		Name:   "aws-access-key",
		Desc:   "AWS access key",
		EnvVar: "AWS_ACCESS_KEY",
	})
	awsSecretKey := *app.String(cli.StringOpt{
		Name:   "aws-secret-key",
		Desc:   "AWS secret key",
		EnvVar: "AWS_SECRET_KEY",
	})
	awsHostedZoneId := *app.String(cli.StringOpt{
		Name:   "aws-hosted-zone-id",
		Desc:   "ID of the hosted the record to be updated resides in",
		EnvVar: "AWS_HOSTED_ZONE_ID",
	})
	awsRecordName := *app.String(cli.StringOpt{
		Name:   "aws-record-name",
		Desc:   "The name of the record to be updated",
		EnvVar: "AWS_RECORD_NAME",
	})
	awsRecordType := *app.String(cli.StringOpt{
		Name:   "aws-record-type",
		Desc:   "The type of the record to be updated",
		EnvVar: "AWS_RECORD_TYPE",
	})

	app.Action = func() {
		svc := route53.New(session.New(), &aws.Config{
			Region:      aws.String(awsRegion),
			Credentials: credentials.NewStaticCredentials(awsAccessKey, awsSecretKey, ""),
		})

		updater := Route53Updater{
			IpServiceUrl: ipServiceUrl,
			svc:          svc,
			HostedZoneId: awsHostedZoneId,
			RecordName:   awsRecordName,
			RecordType:   awsRecordType,
		}

		for {
			poll(updater)
			<-time.After(time.Duration(pollInterval) * time.Second)
		}
	}

	app.Run(os.Args)
}

func poll(updater Updater) {
	currentAddress, err := updater.GetCurrentAddress()
	if err != nil {
		fmt.Printf("Error getting current address: %v\r\n", err)
		return
	}

	recordSet, err := updater.GetRecordSet()
	if err != nil {
		fmt.Printf("Error getting record set: %v\r\n", err)
		return
	}

	recordedAddress, err := updater.GetRecordedAddress(recordSet)
	if err != nil {
		fmt.Printf("Error getting recorded address %v\r\n", err)
		return
	}

	if *currentAddress == *recordedAddress {
		fmt.Printf("Address up to date %s\r\n", *currentAddress)
		return
	}

	err = updater.UpdateRecordedAddress(recordSet, *currentAddress)
	if err != nil {
		fmt.Printf("Error updating recorded address: %v\r\n", err)
	}

	fmt.Printf("Address updated from: %s to: %s\r\n", *recordedAddress, *currentAddress)
}
