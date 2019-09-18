package main

import (
	"context"
	"fmt"
	"net/http"
	"net/url"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-xray-sdk-go/xray"
)

const responseTemplate = `
<html>
<body>
	<img src="http://meetup.com/mu_static/en-US/logo--script.004ada05.svg" alt="Meetup logo" height="44px"><br>
	<br>
	<a href="https://aws.amazon.com/what-is-cloud-computing"><img src="https://d0.awsstatic.com/logos/powered-by-aws.png" alt="Powered by AWS Cloud Computing"></a></br>
	<br>
	<b>%s</b>
</body>
</html>`

const tableName = "Fort_Wayne_AWS_User_Group_Guestbook"

var dynamo *dynamodb.DynamoDB

func handler(ctx context.Context, request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	name := getName(request)

	found, err := findInGuestbook(ctx, dynamo, name)
	if err != nil {
		return createResponse(http.StatusInternalServerError, err.Error()), err
	} else if found {
		return createResponse(http.StatusOK, fmt.Sprintf("You already registred %s! Stop cheating!!", name)), nil
	}

	err = addToGuestbook(ctx, dynamo, name)
	if err != nil {
		return createResponse(http.StatusInternalServerError, err.Error()), err
	}

	return createResponse(http.StatusOK, fmt.Sprintf("Welcome %s! I put your name in the guestbook!!", name)), nil
}

func getName(request events.APIGatewayProxyRequest) string {
	name, ok := request.PathParameters["name"]
	if !ok {
		name = "you"
	}
	name, err := url.PathUnescape(name)
	if err != nil {
		name = "you"
	}
	return name
}

func findInGuestbook(ctx context.Context, dynamo *dynamodb.DynamoDB, name string) (bool, error) {
	input := &dynamodb.QueryInput{
		ExpressionAttributeNames: map[string]*string{
			"#name": aws.String("name"),
		},
		ExpressionAttributeValues: map[string]*dynamodb.AttributeValue{
			":name": {
				S: aws.String(name),
			},
		},
		KeyConditionExpression: aws.String("#name = :name"),
		TableName:              aws.String(tableName),
	}
	output, err := dynamo.QueryWithContext(ctx, input)
	if err != nil {
		return false, fmt.Errorf("error querying guestbook for [%s]: %w", name, err)
	}
	return aws.Int64Value(output.Count) > 0, nil
}

func addToGuestbook(ctx context.Context, dynamo *dynamodb.DynamoDB, name string) error {
	input := &dynamodb.PutItemInput{
		Item: map[string]*dynamodb.AttributeValue{
			"name": {
				S: aws.String(name),
			},
			"timestamp": {
				S: aws.String(fmt.Sprintf("%s", time.Now().Format(time.RFC3339Nano))),
			},
		},
		TableName: aws.String(tableName),
	}
	_, err := dynamo.PutItemWithContext(ctx, input)
	if err != nil {
		return fmt.Errorf("error adding [%s] to guestbook: %w", name, err)
	}
	return nil
}

func createResponse(status int, msg string) events.APIGatewayProxyResponse {
	return events.APIGatewayProxyResponse{
		Body:       fmt.Sprintf(responseTemplate, msg),
		Headers:    map[string]string{"Content-Type": "text/html; charset=utf-8"},
		StatusCode: status,
	}
}

func main() {
	sess := session.New()
	config := aws.NewConfig()
	if os.Getenv("AWS_SAM_LOCAL") == "true" {
		config = config.WithEndpoint("http://docker.for.mac.localhost:8000/")
	}
	dynamo = dynamodb.New(sess, config)
	xray.AWS(dynamo.Client)
	lambda.Start(handler)
}
