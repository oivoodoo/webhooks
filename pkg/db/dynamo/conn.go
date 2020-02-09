package dynamo

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"gitlab.com/oivoodoo/webhooks/pkg"
)

func Connect() *dynamodb.DynamoDB {
	config := pkg.App.Config

	awscfg := &aws.Config{}

	awscfg.WithRegion(config.AWS_REGION)

	sess := session.Must(session.NewSession(awscfg))
	svc := dynamodb.New(sess)

	return svc
}
