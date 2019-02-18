package main

import (
	"encoding/json"
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

type Request events.APIGatewayProxyRequest

type Item struct {
	NoteId string `json:"noteId"`
	UserId string `json:"userId"`
}

// Handler is our lambda handler invoked by the `lambda.Start` function call
func Handler(request events.APIGatewayProxyRequest) (Response, error) {
	fmt.Printf("Received request: %+v\n", request)
	fmt.Printf("body: %s", request.Body)

	sess, _ := session.NewSession(&aws.Config{
		Region: aws.String("ap-southeast-2"),
	})
	service := dynamodb.New(sess)
	fmt.Printf("dynamodb session started:\n%+v", service)

	rawIn := json.RawMessage(request.Body)
	bytes, err := rawIn.MarshalJSON()
	if err != nil {
		panic(err)
	}

	var item Item
	err = json.Unmarshal(bytes, &item)
	if err != nil {
		panic(err)
	}

	fmt.Printf("item after unmarshalling: %+v", item)

	headers := map[string]string{
		"Access-Control-Allow-Origin": "*",
		"Access-Control-Allow-Credentials": "true",
	}

	av, err := dynamodbattribute.MarshalMap(item)
	if err != nil {
		errorMessage := fmt.Sprintf("Error marshalling item: %+v", err)
		return Response{Headers: headers, Body: errorMessage, StatusCode: 404}, nil
	}

	input := &dynamodb.PutItemInput{
		Item: av,
		TableName: aws.String("notes"),
	}

	_ , err = service.PutItem(input)

	if err != nil {
		errorMessage := fmt.Sprintf("Error adding item to table: %+v", err)
		fmt.Printf(errorMessage)
		return Response{Headers: headers, Body: errorMessage, StatusCode: 400}, nil
	}


	return Response{Headers: headers, Body: request.Body, StatusCode: 200}, nil

}

func main() {
	lambda.Start(Handler)
}
