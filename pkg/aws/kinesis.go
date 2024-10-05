package aws

import (
	"context"
	"fmt"
	"time"

	amzn "github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/kinesis"
)

type Kinesis struct {
	client *kinesis.Client
}

func NewKinesis(cfg *amzn.Config) *Kinesis {
	return &Kinesis{
		client: kinesis.NewFromConfig(*cfg),
	}
}

func (k *Kinesis) PublishMessage(stream, msg string) error {
	partationKey := fmt.Sprint(time.Now().Unix())

	record := &kinesis.PutRecordInput{
		StreamName:   &stream,
		Data:         []byte(msg),
		PartitionKey: &partationKey,
	}

	_, err := k.client.PutRecord(context.TODO(), record)
	return err
}
