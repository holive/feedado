package gocloud

import (
	"context"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/sns"
	"github.com/aws/aws-sdk-go/service/sqs"
	"github.com/pkg/errors"
	"gocloud.dev/pubsub"

	// PubSub Packages Providers
	_ "gocloud.dev/pubsub/awssnssqs"
)

const (
	// SQS - Amazon Simple Queue Service
	SQS = "sqs"

	// SNS - Amazon Simple Notification Service
	SNS = "sns"
)

// Client to interact with a PubSub
type Client struct {
	service      string
	region       string
	resourceName string
	url          string
	subscription *pubsub.Subscription
	topic        *pubsub.Topic
}

// Publish content to Topic.
func (ps *Client) Publish(ctx context.Context, message []byte) error {
	msg := &pubsub.Message{
		Body: message,
	}

	err := ps.sendWithRetry(ctx, msg, 5)
	if err != nil {
		return errors.Wrap(err, "Cloud not publish message to topic")
	}

	return nil
}

func (ps *Client) Receive(ctx context.Context) (*pubsub.Message, error) {
	return ps.subscription.Receive(ctx)
}

func (ps *Client) sendWithRetry(ctx context.Context, msg *pubsub.Message, maxCount int) error {
	if maxCount <= 1 {
		return ps.topic.Send(ctx, msg)
	}

	err := ps.topic.Send(ctx, msg)
	if err != nil {
		maxCount--
		<-time.After(300 * time.Millisecond)
		return ps.sendWithRetry(ctx, msg, maxCount)
	}

	return nil
}

// Shutdown Topic
func (ps *Client) Shutdown(ctx context.Context) error {
	return ps.topic.Shutdown(ctx)
}

func (ps *Client) start() error {
	var err error

	err = ps.getURL()
	if err != nil {
		return errors.Wrap(err, "Could not get PubSub URL")
	}

	ps.topic, err = pubsub.OpenTopic(context.Background(), ps.url)
	if err != nil {
		return errors.Wrap(err, "Could not open PubSub topic")
	}

	ps.subscription, err = pubsub.OpenSubscription(context.Background(), ps.url)
	if err != nil {
		return errors.Wrap(err, "Could not open PubSub subscription")
	}

	return nil
}

func (ps *Client) getURL() error {
	switch ps.service {
	case SQS:
		return ps.getSQSURL()
	case SNS:
		return ps.getSNSTopicARN()
	}

	return nil
}

func (ps *Client) getSQSURL() error {
	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ps.region),
	})
	if err != nil {
		return errors.Wrap(err, "Could not init AWS Session")
	}

	result, err := sqs.New(sess).GetQueueUrl(&sqs.GetQueueUrlInput{
		QueueName: aws.String(ps.resourceName),
	})
	if err != nil {
		return errors.Wrap(err, "Could not get SQS URL")
	}

	ps.url = strings.Replace(*result.QueueUrl, "https://", "awssqs://", -1) + "?region=" + ps.region

	return nil
}

func (ps *Client) getSNSTopicARN() error {
	var nextToken *string

	sess, err := session.NewSession(&aws.Config{
		Region: aws.String(ps.region),
	})
	if err != nil {
		return errors.Wrap(err, "Could not init AWS Session")
	}

	snsClient := sns.New(sess)

begin:
	result, err := snsClient.ListTopics(&sns.ListTopicsInput{NextToken: nextToken})
	if err != nil {
		return err
	}

	for _, topic := range result.Topics {
		params := &sns.GetTopicAttributesInput{
			TopicArn: topic.TopicArn,
		}

		resp, err := snsClient.GetTopicAttributes(params)
		if err != nil {
			return err
		}

		arn, ok := resp.Attributes["TopicArn"]
		if !ok {
			return errors.New("missing topic arn")
		}

		if strings.HasSuffix(*arn, ps.resourceName) {
			ps.url = "awssns:///" + *arn + "?region=" + ps.region
			return nil
		}
	}

	if result.NextToken != nil {
		nextToken = result.NextToken
		goto begin
	}

	return errors.New("could not find the SNS topic")
}

// NewClient returns a new client to interact with PubSub.
func NewClient(options ...func(*Client) error) (*Client, error) {
	client := &Client{}

	for _, option := range options {
		if err := option(client); err != nil {
			return nil, errors.Wrap(err, "error during initialization")
		}
	}

	if client.resourceName == "" {
		return nil, errors.New("PubSub resource name cant be empty")
	}

	if client.service == "" {
		return nil, errors.New("PubSub service cant be empty")
	}

	if client.region == "" {
		return nil, errors.New("PubSub region cant be empty")
	}

	if err := client.start(); err != nil {
		return nil, errors.Wrap(err, "Could not start PubSub Client")
	}

	return client, nil
}

// SetResourceName set the logger on PubSub.
func SetResourceName(name string) func(*Client) error {
	return func(client *Client) error {
		client.resourceName = name
		return nil
	}
}

// SetService set the logger on PubSub.
func SetService(service string) func(*Client) error {
	return func(client *Client) error {
		client.service = service
		return nil
	}
}

// SetRegion set the logger on PubSub.
func SetRegion(region string) func(*Client) error {
	return func(client *Client) error {
		client.region = region
		return nil
	}
}
