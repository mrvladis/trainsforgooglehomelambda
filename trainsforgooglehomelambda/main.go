package main

import (
	"bytes"
	"encoding/json"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

//* Global Variables
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var awsRegion, secretName = os.Getenv("AWSRegion"), os.Getenv("secretName")
var awsAccessKeyID, awsSecretAccessKey = os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY")
var awsSessionToken = os.Getenv("AWS_SESSION_TOKEN")

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

var response, googleHomeMessage, message string
var responseToGoogle responseGoogleHome
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
	ApplicationParameters := new(appParams)
	err = json.Unmarshal([]byte(secrets), &ApplicationParameters)
	if err != nil {
		fmt.Println("Couldn't unmarshal Application parameters")
		return clientError(http.StatusUnprocessableEntity)
	}
	fmt.Printf("Ldbws Endpoint:[ %s ] \n", ApplicationParameters.LdbwsEndpoint)
	// Check the incoming Google request

	err = json.Unmarshal([]byte(gRequest.Body), &requestFromGoogle)
	if err != nil {
		fmt.Println("Couldn't unmarshal Application parameters")
		return clientError(http.StatusUnprocessableEntity)
	}
	responseToGoogle.Session.ID = requestFromGoogle.Session.ID

	// Prepare request
	fmt.Println("Preparing XML Soap Request")

	requestTemplate.Header.AccessToken.TokenValue = ApplicationParameters.LdbwsToken
	requestTemplate.Body.Ldb.FilterCrs = "FPK"
	requestTemplate.Body.Ldb.TimeWindow, err = strconv.Atoi(ApplicationParameters.DefaultTimeFrame)
	if err != nil {
		log.Fatal("Failed to convert ApplicationParameters.DefaultTimeFrame value to integer ", err.Error())
		return serverError(err)

	}
	payload, err := xml.MarshalIndent(requestTemplate, "", "  ")
	fmt.Println("Update")
	fmt.Printf("%v", payload)
	fmt.Println("Executing SOAP Request")

	response, err := executeSOAPRequest(payload, "https://lite.realtime.nationalrail.co.uk/OpenLDBWS/ldb11.asmx")
	if err != nil {
		log.Fatal("Error on processing response. ", err.Error())
		return serverError(err)

	}
	if response.StatusCode != 200 {
		fmt.Printf("Request failed with the error code %v and error %v", response.StatusCode, response.Status)
		return serverError(err)
	}

	responseXMLObject := new(responseSoapEnv)
	err = xml.NewDecoder(response.Body).Decode(responseXMLObject)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return serverError(err)
	}

	fmt.Println("Preparing Result")
	trainsCount := len(responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.TrainServices.Service)
	fmt.Printf("There are %v trains", trainsCount)
	currentServices := responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.TrainServices.Service
	fmt.Printf("There are %v trains", trainsCount)
	fmt.Printf("Processing Trains Information")
	if trainsCount > 0 {
		googleHomeMessage = fmt.Sprintf("There are currently %v services scheduled from %v within the next %v minutes:\n", trainsCount, responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.LocationName, ApplicationParameters.DefaultTimeFrame)
		for _, trainService := range currentServices {
			if strings.EqualFold(trainService.Etd, "Cancelled") {
				message = fmt.Sprintf("%v %v %v service has been %v . \n", trainService.Std, trainService.Operator, trainService.Destination.Location.LocationName, trainService.Etd)
			} else {
				message = fmt.Sprintf("%v %v %v service running %v formed of %v coaches. \n", trainService.Std, trainService.Operator, trainService.Destination.Location.LocationName, trainService.Etd, trainService.Length)
			}
			googleHomeMessage += message
		}
	}

	responseToGoogle.Prompt.FirstSimple.Speech = googleHomeMessage
	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(&responseToGoogle)
	reponseToGoogleBody, err := json.Marshal(responseToGoogle)
	fmt.Printf("Marshaled Json: %v", string(reponseToGoogleBody))
	fmt.Printf("Encoded Json: %v", buffer.String())
	if err != nil {
		fmt.Println(err)
		serverError(err)
	}
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
