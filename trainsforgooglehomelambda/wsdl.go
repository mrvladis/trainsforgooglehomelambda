package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/aws/aws-xray-sdk-go/xray"
)

func executeSOAPRequest(ctx context.Context, payload []byte, url string) (*http.Response, error) {

	httpMethod := "POST"
	req, err := http.NewRequestWithContext(ctx, httpMethod, url, bytes.NewReader(payload))
	if err != nil {
		log.Fatal("Error on creating request object. ", err.Error())
		return nil, err
	}

	req.Header.Set("Content-type", "text/xml")

	client := xray.Client(&http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: true,
			},
		},
	})

	res, err := client.Do(req)
	if err != nil {
		log.Fatal("Error on dispatching request. ", err.Error())
		return nil, err
	}

	return res, nil
}

func getTrainsInformation(ctx context.Context, requestSoap requestSoapEnv) (*responseSoapEnv, error) {

	fmt.Println("Preparing XML Soap Request", requestSoap)

	payload, err := xml.MarshalIndent(requestSoap, "", "  ")
	fmt.Println("Payload to be send:")
	os.Stdout.Write(payload)
	fmt.Println()

	// fmt.Println("Update")
	// fmt.Printf("%v", payload)
	fmt.Println("Executing SOAP Request")

	response, err := executeSOAPRequest(ctx, payload, "https://lite.realtime.nationalrail.co.uk/OpenLDBWS/ldb11.asmx")

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
	return responseXMLObject, err

}
