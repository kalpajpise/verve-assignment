package aws

import (
	"context"
	"fmt"

	amzn "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/credentials"
)

// Uses AWS access key and secret key to create an access config
func ConfigWithSecretKey(access, secret, region string) (*amzn.Config, error) {

	creds := amzn.NewCredentialsCache(credentials.NewStaticCredentialsProvider(access, secret, ""))
	cfg, err := config.LoadDefaultConfig(context.TODO(),
		config.WithRegion(region),
		config.WithCredentialsProvider(creds),
	)

	if err != nil {
		return nil, fmt.Errorf("failed to initiliaze aws config with given access and secret key")
	}

	return &cfg, nil
}
