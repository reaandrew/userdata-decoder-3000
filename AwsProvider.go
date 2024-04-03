package main

import (
	"context"
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

func NewAwsProvider(client ec2ClientAPI) *AWSProvider {
	return &AWSProvider{client: client}
}

func NewDefaultAwsProvider() *AWSProvider {
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatalf("Error loading config: %v", err)
	}

	return &AWSProvider{
		client: ec2.NewFromConfig(cfg),
	}
}

func (p *AWSProvider) FetchData() ([]DataOutputPair, error) {
	instancesOutput, err := p.client.DescribeInstances(context.TODO(), &ec2.DescribeInstancesInput{})
	if err != nil {
		return nil, fmt.Errorf("error describing instances: %w", err)
	}

	tasks := make(chan types.Instance, 100) // Channel for tasks
	var outputPairs []DataOutputPair
	var wg sync.WaitGroup
	var mu sync.Mutex // Protects outputPairs

	// Start a fixed number of workers
	for w := 0; w < 10; w++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for instance := range tasks {
				instanceID := *instance.InstanceId
				attributeOutput, err := p.client.DescribeInstanceAttribute(context.TODO(), &ec2.DescribeInstanceAttributeInput{
					Attribute:  types.InstanceAttributeNameUserData,
					InstanceId: aws.String(instanceID),
				})

				if err != nil {
					log.Printf("Error describing instance attributes for %s: %v", instanceID, err)
					continue // Skip this instance on error
				}

				if attributeOutput != nil && attributeOutput.UserData != nil && attributeOutput.UserData.Value != nil {
					userData := *attributeOutput.UserData.Value

					mu.Lock()
					outputPairs = append(outputPairs, DataOutputPair{
						Data:      []byte(userData),
						OutputDir: instanceID,
					})
					fmt.Println("User data collected")
					mu.Unlock()
				} else {
					fmt.Println("User data was empty")
				}
			}
		}()
	}

	// Enqueue tasks
	for _, reservation := range instancesOutput.Reservations {
		for _, instance := range reservation.Instances {
			tasks <- instance
		}
	}
	close(tasks) // Close channel to signal workers to stop

	wg.Wait() // Wait for all workers to finish processing

	return outputPairs, nil
}
