package main

import (
	"context"
	"encoding/base64"
	"fmt"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"log"
	"sync"
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

func (p *AWSProvider) FetchData() ([]DataOutputPair, error) {
	instancesOutput, err := p.client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %w", err)
	}

	var (
		outputPairs []DataOutputPair
		wg          sync.WaitGroup
		mu          sync.Mutex // Protects outputPairs
	)

	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			wg.Add(1)
			go func(instance types.Instance) {
				defer wg.Done()

				instanceID := *instance.InstanceId
				attributeOutput, err := p.client.DescribeInstanceAttribute(context.TODO(), &ec2.DescribeInstanceAttributeInput{
					Attribute:  types.InstanceAttributeNameUserData,
					InstanceId: aws.String(instanceID),
				})

				if err != nil {
					// Log the error; adjust error handling as necessary for your application.
					log.Printf("Error describing instance attributes for %s: %v", instanceID, err)
					return // Proceed with the next iteration.
				}

				if attributeOutput != nil && attributeOutput.UserData != nil && attributeOutput.UserData.Value != nil {
					userData, decodeErr := base64.StdEncoding.DecodeString(*attributeOutput.UserData.Value)
					if decodeErr != nil {
						log.Printf("Error decoding user data for instance %s: %v", instanceID, decodeErr)
						return
					}

					mu.Lock()
					outputPairs = append(outputPairs, DataOutputPair{
						Data:      userData,
						OutputDir: instanceID,
					})

					fmt.Println("User data collected")
					mu.Unlock()
				} else {
					fmt.Println("User Data is Empty")
				}
			}(instance)
		}
	}

	wg.Wait() // Wait for all goroutines to complete
	return outputPairs, nil
}
