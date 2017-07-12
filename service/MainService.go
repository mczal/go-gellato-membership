package service

import (
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

var svc *dynamodb.DynamoDB // mczal: this variable will be visible to all services

// Only called once from Main.go
func GenerateDynamoDBSvc() {
	sess := session.Must(session.NewSession())
	// sess := session.Must(session.NewSession(&aws.Config{
	// 	Region: aws.String(endpoints.ApSoutheast1RegionID),
	// }))

	// creds := stscreds.NewCredentials(sess, "myRoleArn")
	svc = dynamodb.New(sess)
}
