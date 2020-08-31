package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type appParams struct {
	LdbwsToken       string `json:"LdbwsToken"`
	LdbwsEndpoint    string `json:"Ldbwsendpoint"`
	DestPreposition  string `json:"DestPreposition"`
	SrcPreposition   string `json:"SrcPreposition"`
	DefaultTimeFrame string `json:"DefaultTimeFrame"`
}

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var awsRegion, secretName = os.Getenv("AWSRegion"), os.Getenv("secretName")

type soapenv struct {
	XMLName xml.Name
	Header  header
	Body    body
}

type header struct {
	AccessToken token `xml:"AccessToken"`
}

type token struct {
	TokenValue string `xml:"TokenValue"`
}

type body struct {
	Ldb ldb `xml:"GetDepBoardWithDetailsRequest"`
}

type ldb struct {
	NumRows    int    `xml:"numRows"`
	Crs        string `xml:"crs"`
	FilterCrs  string `xml:"filterCrs"`
	FilterType string `xml:"filterType"`
	TimeOffset int    `xml:"timeOffset"`
	TimeWindow int    `xml:"timeWindow"`
}

var payload = []byte(strings.TrimSpace(`
<soapenv:Envelope xmlns:soapenv="http://schemas.xmlsoap.org/soap/envelope/" 
xmlns:typ="http://thalesgroup.com/RTTI/2013-11-28/Token/types" 
xmlns:ldb="http://thalesgroup.com/RTTI/2017-10-01/ldb/">
<soapenv:Header>
	<typ:AccessToken>
		<typ:TokenValue>TokenWasHere</typ:TokenValue>
	</typ:AccessToken>
</soapenv:Header>
<soapenv:Body>
	<ldb:GetDepBoardWithDetailsRequest>
		<ldb:numRows>10</ldb:numRows>
		<ldb:crs>OKL</ldb:crs>
		<ldb:filterCrs>?</ldb:filterCrs>
		<ldb:filterType>to</ldb:filterType>
		<ldb:timeOffset>0</ldb:timeOffset>
		<ldb:timeWindow>60</ldb:timeWindow>
	</ldb:GetDepBoardWithDetailsRequest>
</soapenv:Body>
</soapenv:Envelope>`,
))

var response string

func router(req events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	switch req.HTTPMethod {
	case "GET":
		return processRequest(req)
		/* 	case "POST":
		return create(req) */
	default:
		return clientError(http.StatusMethodNotAllowed)
	}
}
func processRequest(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {

	//Get the Application Parameters
	secrets, err := getSecret()
	if err != nil {
		fmt.Printf("Coudn't retreive Secrets")
		return serverError(err)
	}
	if secrets == "" {
		fmt.Println("Secrets  are empty")
		return serverError(err)
	}
	ApplicationParameters := new(appParams)
	err = json.Unmarshal([]byte(secrets), &ApplicationParameters)
	if err != nil {
		fmt.Println("Couldn't unmarshal Application parameters")
		return clientError(http.StatusUnprocessableEntity)
	}
	fmt.Printf("Ldbws Endpoint:[ %s ] \n", ApplicationParameters.LdbwsEndpoint)
	// Prepare request

	fmt.Println("Preparing XML Soap Request")
	xmlObject := new(soapenv)
	fmt.Println("Unmarshaling Template")
	xml.Unmarshal([]byte(payload), &xmlObject)
	fmt.Println("Updating Token Value")
	xmlObject.Header.AccessToken.TokenValue = ApplicationParameters.LdbwsToken
	fmt.Println("Update")

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(response),
	}, nil
}

// Add a helper for handling errors. This logs any error to os.Stderr
// and returns a 500 Internal Server Error response that the AWS API
// Gateway understands.
func serverError(err error) (events.APIGatewayProxyResponse, error) {
	errorLogger.Println(err.Error())

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusInternalServerError,
		Body:       http.StatusText(http.StatusInternalServerError),
	}, nil
}

// Similarly add a helper for send responses relating to client errors.
func clientError(status int) (events.APIGatewayProxyResponse, error) {
	return events.APIGatewayProxyResponse{
		StatusCode: status,
		Body:       http.StatusText(status),
	}, nil
}

func main() {
	lambda.Start(router)
}
