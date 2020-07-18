package main

import (
	"context"
	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"
)

func Handler(ctx context.Context, request events.APIGatewayProxyRequest) (*events.APIGatewayProxyResponse, error) {
	req := CreateRequest(request)
	userData := GetUserSession(req)

	return &events.APIGatewayProxyResponse{
		StatusCode: 200,
		Headers:    map[string]string{"Content-Type": "application/json"},
		Body:       userData,
	}, nil
}

func CreateRequest(request events.APIGatewayProxyRequest) *http.Request {
	endpoint, err := url.Parse("https://hackersandslackers.app/members/api/member/")
	var headers http.Header
	headers.Add("cookie", request.Body)
	if err != nil {
		log.Fatal(err)
	}
	req := &http.Request{URL: endpoint, Header: headers}
	return req
}

func GetUserSession(req *http.Request) string {
	// Create HTTP client with timeout
	client := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Request account information by session token.
	res, reqError := client.Do(req)
	if reqError != nil {
		log.Fatal(reqError)
	}
	if res.StatusCode != 200 {
		log.Fatal("status code error: %i", res.StatusCode)
	}
	defer res.Body.Close()

	// Parse response
	data, bodyErr := ioutil.ReadAll(res.Body)
	if bodyErr != nil {
		log.Fatal(bodyErr)
	}
	return string(data)
}

func main() {
	// Make the Handler available
	lambda.Start(Handler)
}
