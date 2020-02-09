package webhooks

import "github.com/aws/aws-sdk-go/service/dynamodb"

type Repo interface {
	BatchInsert(webhooks []*Webhook) error
}

type DynamoRepo struct {
	Conn *dynamodb.DynamoDB
}

func (DynamoRepo) BatchInsert(webhooks []*Webhook) error {
	return nil
}

type MockRepo struct {
}

func (MockRepo) BatchInsert(webhooks []*Webhook) error {
	return nil
}
