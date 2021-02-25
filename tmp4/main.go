package main

import (
	"encoding/xml"
	"fmt"
	"log"
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

	fmt.Println("Preparing to process StationFrom Parameter, if available.")

	fmt.Println("Preparing to process Time Parameter, if available.")

	requestSoap.Header.AccessToken.TokenValue = "a35d9469-3417-4279-9812-aaa8d0a17db9"
	requestSoap.Body.Ldb.TimeWindow = 60

	fmt.Println("Produced Saop request body", requestSoap)
}
