package main

import (
	"fmt"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

// Response is of type APIGatewayProxyResponse since we're leveraging the
// AWS Lambda Proxy Request functionality (default behavior)
//
// https://serverless.com/framework/docs/providers/aws/events/apigateway/#lambda-proxy-integration
type Response events.APIGatewayProxyResponse

type Item struct {
	NoteId string `json:"noteId"`
	UserId string `json:"userId"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Println("Received body: ", request.Body)
	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	service := dynamodb.New(sess)

	item := request.Body

	av, err := dynamodbattribute.MarshalMap(item)

	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling item: %+v", err)
		return Response{Body: errorMessage, StatusCode: 404}, nil
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String("notes"),
	}

	_ , err = service.PutItem(input)

	if err != nil {
		errorMessage := fmt.Sprintf("Error adding item ot table: %+v", err)
		fmt.Printf(errorMessage)
		return Response{ Body: errorMessage, StatusCode: 400}, nil
	}

	return Response{Body: request.Body, StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}
