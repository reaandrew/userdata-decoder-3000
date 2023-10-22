package main

import (
	"context"
	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/ec2"
	"github.com/aws/aws-sdk-go-v2/service/ec2/types"
	"github.com/stretchr/testify/assert"
	"testing"
)

type MockEC2Client struct {
	ec2ClientAPI
}

type mockEC2Client struct{}

func (m mockEC2Client) DescribeInstances(ctx context.Context, params *ec2.DescribeInstancesInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstancesOutput, error) {
	return &ec2.DescribeInstancesOutput{
		Reservations: []types.Reservation{
			{
				Instances: []types.Instance{
					{
						InstanceId: aws.String("i-1234567890"),
					},
				},
			},
		},
	}, nil
}

func (m mockEC2Client) DescribeInstanceAttribute(ctx context.Context, params *ec2.DescribeInstanceAttributeInput, optFns ...func(*ec2.Options)) (*ec2.DescribeInstanceAttributeOutput, error) {
	return &ec2.DescribeInstanceAttributeOutput{
		UserData: &types.AttributeValue{Value: aws.String("something")},
	}, nil
}

func TestFetchData(t *testing.T) {
	mockClient := &mockEC2Client{}

	p := NewAwsProvider(mockClient)

	outputPairs, err := p.FetchData()
	assert.NoError(t, err)
	assert.NotNil(t, outputPairs)

	assert.Equal(t, 1, len(outputPairs))
	if len(outputPairs) > 0 {
		assert.Equal(t, "i-1234567890", outputPairs[0].OutputDir)
		assert.Equal(t, "something", string(outputPairs[0].Data))
	}
}
