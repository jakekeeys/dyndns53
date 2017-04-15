package main

import (
	"fmt"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53"
	"io/ioutil"
	"net/http"
	"strings"
)

const (
	MAX_ITEMS_1   = "1"
	ACTION_UPSERT = "UPSERT"
)

type Updater interface {
	GetCurrentAddress() (*string, error)
	GetRecordSet() (*route53.ResourceRecordSet, error)
	GetRecordedAddress(recordSet *route53.ResourceRecordSet) (*string, error)
	UpdateRecordedAddress(recordSet *route53.ResourceRecordSet, currentAddress string) error
}

type Route53Updater struct {
	IpServiceUrl string
	svc          *route53.Route53
	HostedZoneId string
	RecordName   string
	RecordType   string
}

func (route53Updater Route53Updater) GetCurrentAddress() (*string, error) {
	resp, err := http.Get(route53Updater.IpServiceUrl)
	if err != nil {
		fmt.Errorf("Error getting response: %v", err)
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("Error reading response body: %v", err)
	}

	currentAddress := strings.TrimSpace(string(body))
	return &currentAddress, nil
}

func (route53Updater Route53Updater) GetRecordSet() (*route53.ResourceRecordSet, error) {
	params := &route53.ListResourceRecordSetsInput{
		HostedZoneId:    aws.String(route53Updater.HostedZoneId),
		MaxItems:        aws.String(MAX_ITEMS_1),
		StartRecordName: aws.String(route53Updater.RecordName),
		StartRecordType: aws.String(route53Updater.RecordType),
	}

	resp, err := route53Updater.svc.ListResourceRecordSets(params)
	if err != nil {
		return nil, fmt.Errorf("Error listing records: %v", err)
	}
	if len(resp.ResourceRecordSets) < 1 {
		return nil, fmt.Errorf("No records sets for: %v", params)
	}

	return resp.ResourceRecordSets[0], nil
}

func (route53Updater Route53Updater) GetRecordedAddress(recordSet *route53.ResourceRecordSet) (*string, error) {
	if len(recordSet.ResourceRecords) < 1 {
		return nil, fmt.Errorf("No records for set %v", recordSet)
	}

	recordedAddress := strings.TrimSpace(*recordSet.ResourceRecords[0].Value)
	return &recordedAddress, nil
}

func (route53Updater Route53Updater) UpdateRecordedAddress(recordSet *route53.ResourceRecordSet, currentAddress string) error {
	recordSet.ResourceRecords[0].Value = &currentAddress

	params := &route53.ChangeResourceRecordSetsInput{
		ChangeBatch: &route53.ChangeBatch{
			Changes: []*route53.Change{
				{
					Action:            aws.String(ACTION_UPSERT),
					ResourceRecordSet: recordSet,
				},
			},
		},
		HostedZoneId: aws.String(route53Updater.HostedZoneId),
	}
	_, err := route53Updater.svc.ChangeResourceRecordSets(params)
	if err != nil {
		return fmt.Errorf("Error applying changes: %v", params)
	}

	return nil
}
