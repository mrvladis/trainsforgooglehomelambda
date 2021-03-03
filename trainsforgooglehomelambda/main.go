package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//* Global Variables
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var awsRegion, secretName = os.Getenv("AWSRegion"), os.Getenv("secretName")
var awsAccessKeyID, awsSecretAccessKey = os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY")
var awsSessionToken = os.Getenv("AWS_SESSION_TOKEN")
var applicationParameters appParams

var requestTemplate = requestSoapEnv{
	XMLNsSoapEnv: "http://schemas.xmlsoap.org/soap/envelope/",
	XMLNsTyp:     "http://thalesgroup.com/RTTI/2013-11-28/Token/types",
	XMLNsLDB:     "http://thalesgroup.com/RTTI/2017-10-01/ldb/",
	Header: requestHeader{
		AccessToken: requestToken{
			TokenValue: "Future Token",
		},
	},
	Body: requestBody{
		Ldb: requestLdb{
			NumRows:    10,
			Crs:        "OKL",
			FilterCrs:  "?",
			FilterType: "to",
			TimeOffset: 0,
			TimeWindow: 60,
		},
	},
}

var response string
var err error

//var responseToGoogle responseGoogleHome
var requestFromGoogle requestGoogleHome

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
func processRequest(gRequest events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var APIGatewayProxyResponse events.APIGatewayProxyResponse

	fmt.Println("Getting the Application Parameters")
	fmt.Printf("Request Body: %v", gRequest.Body)
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

	err = json.Unmarshal([]byte(secrets), &applicationParameters)
	if err != nil {
		fmt.Println("Couldn't unmarshal Application parameters")
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("UnmarshalTypeError %v - %v - %v\n", ute.Value, ute.Type, ute.Offset)
		} else {
			fmt.Println("Other error:", err)
		}
		return clientError(http.StatusUnprocessableEntity)
	}
	fmt.Printf("Ldbws Endpoint:[ %s ] \n", applicationParameters.LdbwsEndpoint)

	// Check the incoming Google request
	//Unmarshall Googl request
	err = json.Unmarshal([]byte(gRequest.Body), &requestFromGoogle)
	if err != nil {
		fmt.Println("Couldn't unmarshal Google Request")
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("UnmarshalTypeError %v - %v - %v\n", ute.Value, ute.Type, ute.Offset)
		} else {
			fmt.Println("Other error:", err)
		}
		return clientError(http.StatusUnprocessableEntity)

	}
	//Analysing the initial query
	intentNameValue := *requestFromGoogle.Intent.Name
	switch intentNameValue {
	case "Initial_Train_Check":
		fmt.Println("Performing intial train information request")
		APIGatewayProxyResponse, err = initialTrainCheck(requestFromGoogle)
	default:
		return clientError(http.StatusMethodNotAllowed)
	}

	return APIGatewayProxyResponse, err

}

func initialTrainCheck(requestFromGoogle requestGoogleHome) (events.APIGatewayProxyResponse, error) {
	var buffer bytes.Buffer

	requestSoap, err := prepareRequestToNationalRail(requestFromGoogle)

	if err != nil {
		log.Fatal("Failed to process google Request ", err.Error())
		return serverError(err)

	}
	//
	// Prepare request
	responseXMLObject, err := getTrainsInformation(requestSoap)

	if err != nil {
		log.Fatal("Error on processing response. ", err.Error())
		return serverError(err)
	}
	fmt.Println("LDBWS response received:")
	output, err := xml.MarshalIndent(responseXMLObject, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
	fmt.Println()

	fmt.Println("Preparing Reponse to Google")
	responseToGoogle, err := prepareGoogleResponse(responseXMLObject)
	fmt.Println("Updating Google Response Session ID")
	responseToGoogle.Session.ID = requestFromGoogle.Session.ID
	fmt.Println("Encoding Google Response into JSON")
	json.NewEncoder(&buffer).Encode(&responseToGoogle)
	reponseToGoogleBody, err := json.Marshal(responseToGoogle)
	if err != nil {
		fmt.Println(err)
		serverError(err)
	}
	fmt.Println("Marshaled Json:", string(reponseToGoogleBody))
	fmt.Println("Encoded Json:", buffer.String())
	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       buffer.String(),
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
	lambda.Start(processRequest)
}
