package db

import (
	"os"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
)

func TestDynamoDBStore(t *testing.T) {
	accessKey := os.Getenv("SB_AWS_ACCESS_KEY")
	if accessKey == "" {
		t.Skip("Skipping test, aws access key is not set")
	}

	secretKey := os.Getenv("SB_AWS_SECRET_KEY")
	if secretKey == "" {
		t.Skip("Skipping test, aws secret key is not set")
	}

	region := os.Getenv("SB_AWS_REGION")
	if region == "" {
		region = "us-west-2"
	}

	table := os.Getenv("SB_DYNAMODB_TEST_TABLE")
	if table == "" {
		t.Skip("Skipping test, dynamodb test table is not set")
	}

	config := &aws.Config{
		Credentials: credentials.NewStaticCredentials(accessKey, secretKey, ""),
		Region:      aws.String(region),
	}

	testStore(t, NewDynamoDBStore(session.New(config), table))
}
