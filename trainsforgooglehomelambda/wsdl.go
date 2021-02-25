package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"log"
	"net/http"
)

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

func getTrainsInformation(requestSoap requestSoapEnv) (*responseSoapEnv, error) {

	fmt.Println("Preparing XML Soap Request", requestSoap)

	payload, err := xml.MarshalIndent(requestSoap, "", "  ")
	// fmt.Println("Update")
	// fmt.Printf("%v", payload)
	fmt.Println("Executing SOAP Request")

	response, err := executeSOAPRequest(payload, "https://lite.realtime.nationalrail.co.uk/OpenLDBWS/ldb11.asmx")

	if err != nil {
		log.Fatal("Error on processing response. ", err.Error())

	}
	if response.StatusCode != 200 {
		fmt.Printf("Request failed with the error code %v and error %v", response.StatusCode, response.Status)
	}

	responseXMLObject := new(responseSoapEnv)
	err = xml.NewDecoder(response.Body).Decode(responseXMLObject)
	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
	}
	println("Response received from the LDBWS:", responseXMLObject)
	return responseXMLObject, err

}
