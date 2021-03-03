package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"
)

//* Global Variables
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

var response, googleHomeMessage, message string
var responseToGoogle responseGoogleHome
var requestFromGoogle requestGoogleHome
var simpleR gSimple
var promptR gPrompt

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
var err error

func main() {

	requestSoap := requestTemplate

	output, err := xml.MarshalIndent(requestSoap, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}

	os.Stdout.Write(output)
	fmt.Print("Initial Saop request body", requestSoap)
	fmt.Println("Preparing to process StationTo Parameter, if available.")

	destinationStation := "OKL"
	requestSoap.Body.Ldb.FilterCrs = destinationStation
	fmt.Println("Target Station code identified as:", destinationStation)
	fmt.Println("Updated Saop request body", requestSoap)

	sourceStation := "KGX"
	requestSoap.Body.Ldb.Crs = sourceStation
	fmt.Println("Preparing to process StationFrom Parameter, if available.")

	fmt.Println("Preparing to process Time Parameter, if available.")

	requestSoap.Header.AccessToken.TokenValue = "a35d9469-3417-4279-9812-aaa8d0a17db9"
	requestSoap.Body.Ldb.TimeWindow = 60

	fmt.Println("Produced Saop request body", requestSoap)

	responseXMLObject, err := getTrainsInformation(requestSoap)

	output, err = xml.MarshalIndent(responseXMLObject, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
	fmt.Println()

}

func getTrainsInformation(requestSoap requestSoapEnv) (*responseSoapEnv, error) {

	fmt.Println("Preparing XML Soap Request", requestSoap)

	payload, err := xml.MarshalIndent(requestSoap, "", "  ")
	fmt.Println("Payload to be send:")
	os.Stdout.Write(payload)
	fmt.Println()

	// fmt.Println("Update")
	// fmt.Printf("%v", payload)
	fmt.Println("Executing SOAP Request")

	response, err := executeSOAPRequest(payload, "https://lite.realtime.nationalrail.co.uk/OpenLDBWS/ldb11.asmx")

	if err != nil {
		fmt.Printf("Request failed with the error %v", err.Error())
		log.Fatal("Error on processing response. ", err.Error())

	}
	if response.StatusCode != 200 {
		fmt.Printf("HTTP Request returned the following error code %v and error %v", response.StatusCode, response.Status)
	}

	responseXMLObject := new(responseSoapEnv)
	err = xml.NewDecoder(response.Body).Decode(responseXMLObject)
	if err != nil {
		fmt.Printf("Unmarshaling xml failed with the error %v", err.Error())
		log.Fatal("Error on unmarshaling xml. ", err.Error())
	}
	println("Response received from the LDBWS:")
	output, err := xml.MarshalIndent(responseXMLObject, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
	fmt.Println()
	return responseXMLObject, err

}

func executeSOAPRequest(payload []byte, url string) (*http.Response, error) {

	httpMethod := "POST"
	req, err := http.NewRequest(httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return nil, err
	}

	req.Header.Set("Content-type", "text/xml")

	client := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	}

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return nil, err
	}

	return res, nil
}
