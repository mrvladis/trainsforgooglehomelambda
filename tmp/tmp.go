package main

import (
	"encoding/xml"
	"fmt"
	"strings"
)

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

func main() {
	xmlObject := new(soapenv)
	xml.Unmarshal([]byte(payload), &xmlObject)
	fmt.Println("Start")
	fmt.Printf("%v", xmlObject.Header.AccessToken.TokenValue)
	xmlObject.Header.AccessToken.TokenValue = "MyNewToken"
	fmt.Println("Update")
	fmt.Printf("%v", xmlObject)

}
