package main

import (
	"encoding/json"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-lambda-go/events"
	"github.com/aws/aws-lambda-go/lambda"
)

type requestSoapEnv struct {
	XMLName      xml.Name      `xml:"soapenv:Envelope"`
	XMLNsSoapEnv string        `xml:"xmlns:soapenv,attr"`
	XMLNsTyp     string        `xml:"xmlns:typ,attr"`
	XMLNsLDB     string        `xml:"xmlns:ldb,attr"`
	Header       requestHeader `xml:"soapenv:Header"`
	Body         requestBody   `xml:"soapenv:Body"`
}

type requestHeader struct {
	AccessToken requestToken `xml:"typ:AccessToken"`
}

type requestToken struct {
	TokenValue string `xml:"typ:TokenValue"`
}

type requestBody struct {
	Ldb requestLdb `xml:"ldb:GetDepBoardWithDetailsRequest"`
}

type requestLdb struct {
	NumRows    int    `xml:"ldb:numRows"`
	Crs        string `xml:"ldb:crs"`
	FilterCrs  string `xml:"ldb:filterCrs"`
	FilterType string `xml:"ldb:filterType"`
	TimeOffset int    `xml:"ldb:timeOffset"`
	TimeWindow int    `xml:"ldb:timeWindow"`
}

type appParams struct {
	LdbwsToken       string `json:"LdbwsToken"`
	LdbwsEndpoint    string `json:"Ldbwsendpoint"`
	DestPreposition  string `json:"DestPreposition"`
	SrcPreposition   string `json:"SrcPreposition"`
	DefaultTimeFrame string `json:"DefaultTimeFrame"`
}

// Response structure
type responseSoapEnv struct {
	XMLName xml.Name
	Body    responseBody `xml:"Body"`
}

type responseBody struct {
	XMLName                        xml.Name
	GetDepBoardWithDetailsResponse responseBoardWithDetailsResponse `xml:"GetDepBoardWithDetailsResponse"`
}

type responseBoardWithDetailsResponse struct {
	XMLName               xml.Name
	GetStationBoardResult responseStationBoardResult `xml:"GetStationBoardResult"`
}

type responseStationBoardResult struct {
	XMLName            xml.Name
	GeneratedAt        string                `xml:"generatedAt"`
	LocationName       string                `xml:"locationName"`
	Crs                string                `xml:"crs"`
	FilterLocationName string                `xml:"filterLocationName"`
	Filtercrs          string                `xml:"filtercrs"`
	PlatformAvailable  bool                  `xml:"platformAvailable"`
	TrainServices      responseTrainServices `xml:"trainServices"`
}

type responseTrainServices struct {
	XMLName xml.Name
	Service []responseService `xml:"service"`
}
type responseService struct {
	XMLName                 xml.Name
	Std                     string                    `xml:"std"`
	Etd                     string                    `xml:"etd"`
	Platform                string                    `xml:"platform"`
	Operator                string                    `xml:"operator"`
	OperatorCode            string                    `xml:"operatorCode"`
	ServiceType             string                    `xml:"serviceType"`
	Length                  int                       `xml:"length"`
	ServiceID               string                    `xml:"serviceID"`
	Rsid                    string                    `xml:"rsid"`
	Origin                  responseOrigin            `xml:"origin"`
	Destination             responseOrigin            `xml:"destination"`
	SubsequentCallingPoints responseCallingPointsList `xml:"subsequentCallingPoints"`
	Other                   interface{}               `xml:",any"`
}

type responseCallingPointsList struct {
	XMLName          xml.Name
	CallingPointList responseCallingPoint `xml:"callingPointList"`
}

type responseCallingPoint struct {
	XMLName      xml.Name
	CallingPoint []responsePoint `xml:"callingPoint"`
}

type responsePoint struct {
	XMLName      xml.Name
	LocationName string `xml:"locationName"`
	Crs          string `xml:"crs"`
	St           string `xml:"st"`
	Et           string `xml:"et"`
	Length       int    `xml:"length"`
}

type responseOrigin struct {
	Location responseLocation `xml:"location"`
}

type responseLocation struct {
	XMLName      xml.Name
	LocationName string `xml:"locationName"`
	Crs          string `xml:"crs"`
}

var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)
var awsRegion, secretName, AWS_ACCESS_KEY_ID, AWS_SECRET_ACCESS_KEY, AWS_SESSION_TOKEN = os.Getenv("AWSRegion"), os.Getenv("secretName"), os.Getenv("AWS_ACCESS_KEY_ID"), os.Getenv("AWS_SECRET_ACCESS_KEY"), os.Getenv("AWS_SESSION_TOKEN")

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
	fmt.Println("Getting the Application Parameters")

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

	requestTemplate.Header.AccessToken.TokenValue = ApplicationParameters.LdbwsToken
	requestTemplate.Body.Ldb.FilterCrs = "FPK"

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
	fmt.Println("Preparing Result")
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return serverError(err)
	}

	fmt.Println("Result:")
	fmt.Println(string(bodyBytes))

	/* fmt.Println("Processing XML Soap Response")
	xmlObject := new(responseSoapEnv)
	xmlObject1 := new(responseSoapEnv)
	fmt.Println("Unmarshaling Template")
	err = xml.NewDecoder(response.Body).Decode(xmlObject)
	xml.Unmarshal(bodyBytes, &xmlObject1)

	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	} */

	//xmlObject.Header.AccessToken.TokenValue = ApplicationParameters.LdbwsToken

	return events.APIGatewayProxyResponse{
		StatusCode: http.StatusOK,
		Body:       string(string(bodyBytes)),
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
