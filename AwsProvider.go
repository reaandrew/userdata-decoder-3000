package main

import (
	"context"
	"fmt"
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
	// Initialize EC2 client
	return AWSProvider{
		client: ec2.NewFromConfig(cfg),
	}
}

func (p AWSProvider) FetchData() ([]DataOutputPair, error) {
	// List instances (you might want to filter them)
	instancesOutput, err := p.client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		log.Fatalf("Error describing instances: %v", err)
	}

	var outputPairs []DataOutputPair
	// Loop through reservations and instances
	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			instanceID := *instance.InstanceId

			// Describe instance attributes to get user data
			attributeOutput, err := p.client.DescribeInstanceAttribute(context.TODO(), &ec2.DescribeInstanceAttributeInput{
				Attribute:  "userData",
				InstanceId: aws.String(instanceID),
			})
			if err != nil {
				log.Fatalf("Error describing instance attributes for %s: %v", instanceID, err)
			}

			// Display user data (base64-encoded)
			if attributeOutput.UserData != nil {
				userData := *attributeOutput.UserData.Value
				fmt.Printf("Instance ID: %s, User Data: %s\n", instanceID, userData)
				outputPairs = append(outputPairs, DataOutputPair{
					Data:      []byte(userData),
					OutputDir: instanceID,
				})
			} else {
				fmt.Printf("Instance ID: %s, User Data: None\n", instanceID)
			}
		}
	}

	// Fetch data from AWS and return it as a slice of DataOutputPair
	return outputPairs, nil
}
