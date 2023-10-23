package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"log"
)

type AWSProvider struct {
	client ec2ClientAPI
}

type ec2ClientAPI interface {
	DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error)
	DescribeInstanceAttribute(ctx context.Context, params *ec2.DescribeInstanceAttributeInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceAttributeOutput, error)
}

func NewAwsProvider(client ec2ClientAPI) AWSProvider {
	return AWSProvider{client: client}
}

func NewDefaultAwsProvider() AWSProvider {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return AWSProvider{
		client: ec2.NewFromConfig(cfg),
	}
}

func (p AWSProvider) FetchData() ([]DataOutputPair, error) {

	instancesOutput, err := p.client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatalf("Error describing instances: %v", err)
	}

	var outputPairs []DataOutputPair

	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instanceID := *instance.InstanceId

			attributeOutput, err := p.client.DescribeInstanceAttribute(context.TODO(), &ec2.DescribeInstanceAttributeInput{
				Attribute:  "userData",
				InstanceId: aws.String(instanceID),
			})
			if err != nil {
				log.Fatalf("Error describing instance attributes for %s: %v", instanceID, err)
			}

			if attributeOutput.UserData != nil {
				userData := *attributeOutput.UserData.Value
				outputPairs = append(outputPairs, DataOutputPair{
					Data:      []byte(userData),
					OutputDir: instanceID,
				})
			}
		}
	}

	return outputPairs, nil
}
