package main

import (
	"bytes"
	"crypto/tls"
	"encoding/xml"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
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

// Response structure
type responseSoapEnv struct {
	XMLName   xml.Name
	XMLNSsoap string       `xml:"xmlns:soap,attr"`
	XMLNSxsi  string       `xml:"xmlns:xsi,attr"`
	XMLNSxsd  string       `xml:"xmlns:xsd,attr"`
	Body      responseBody `xml:"Body"`
	//	Other     interface{}  `xml:",any"`
}

type responseBody struct {
	XMLName                        xml.Name
	GetDepBoardWithDetailsResponse responseBoardWithDetailsResponse `xml:"GetDepBoardWithDetailsRequest,attr"`
	//	Other                          interface{}                      `xml:",any"`
}

type responseBoardWithDetailsResponse struct {
	XMLName               xml.Name
	GetStationBoardResult responseStationBoardResult `xml:"GetDepBoardWithDetailsRequest,attr"`
	//	Other                 interface{}                `xml:",any"`
}

type responseStationBoardResult struct {
	XMLName            xml.Name
	GeneratedAt        string                  `xml:"generatedAt"`
	LocationName       string                  `xml:"locationName"`
	Crs                string                  `xml:"crs"`
	FilterLocationName string                  `xml:"filterLocationName"`
	Filtercrs          string                  `xml:"filtercrs"`
	PlatformAvailable  bool                    `xml:"platformAvailable"`
	TrainServices      []responseTrainServices `xml:"trainServices"`
	//	Other              interface{}             `xml:",any"`
}

type responseTrainServices struct {
	XMLName                 xml.Name
	Std                     string                          `xml:"std"`
	Etd                     string                          `xml:"etd"`
	Platform                string                          `xml:"platform"`
	Operator                string                          `xml:"operator"`
	OperatorCode            string                          `xml:"operatorCode"`
	ServiceType             string                          `xml:"serviceType"`
	Length                  int                             `xml:"length"`
	ServiceID               string                          `xml:"serviceID"`
	Rsid                    string                          `xml:"rsid"`
	Origin                  responseOrigin                  `xml:"origin"`
	Destination             responseOrigin                  `xml:"destination"`
	SubsequentCallingPoints responseSubsequentCallingPoints `xml:"subsequentCallingPoints"`
	// Other                   interface{}                     `xml:",any"`
}

type responseOrigin struct {
	XMLName  xml.Name
	Location responseLocation `xml:"location"`
	//	Other    interface{}      `xml:",any"`
}

type responseLocation struct {
	XMLName      xml.Name
	LocationName string `xml:"locationName"`
	Crs          string `xml:"crs"`
	//	Other        interface{} `xml:",any"`
}

type responseSubsequentCallingPoints struct {
	XMLName          xml.Name
	CallingPointList []responseCallingPointList `xml:"callingPointList"`
	//Other            interface{}                `xml:",any"`
}

type responseCallingPointList struct {
	XMLName      xml.Name
	LocationName string `xml:"locationName"`
	Crs          string `xml:"crs"`
	St           string `xml:"st"`
	Et           string `xml:"et"`
	Length       int    `xml:"length"`
	//	Other        interface{} `xml:",any"`
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

func main() {
	fmt.Println("Start")
	requestTemplate.Header.AccessToken.TokenValue = "a35d9469-3417-4279-9812-aaa8d0a17db9"
	requestTemplate.Body.Ldb.FilterCrs = "FPK"

	payload, err := xml.MarshalIndent(requestTemplate, "", "  ")
	fmt.Println("Update")
	fmt.Printf("%v", payload)
	fmt.Println("Executing SOAP Request")

	response, err := executeSOAPRequest(payload, "https://lite.realtime.nationalrail.co.uk/OpenLDBWS/ldb11.asmx")
	if err != nil {
		log.Fatal("Error on processing response. ", err.Error())
		return

	}
	if response.StatusCode != 200 {
		fmt.Printf("Request failed with the error code %v and error %v", response.StatusCode, response.Status)
		return
	}
	fmt.Println("Preparing Result")
	bodyBytes, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return
	}

	fmt.Println("Result:")
	fmt.Println(string(bodyBytes))

	fmt.Println("Processing XML Soap Response")
	xmlObject := new(responseSoapEnv)
	xmlObject1 := new(responseSoapEnv)
	fmt.Println("Unmarshaling Template")
	err = xml.NewDecoder(response.Body).Decode(xmlObject)
	xml.Unmarshal(bodyBytes, &xmlObject1)

	if err != nil {
		log.Fatal("Error on unmarshaling xml. ", err.Error())
		return
	}

	//xmlObject.Header.AccessToken.TokenValue = ApplicationParameters.LdbwsToken
	fmt.Printf("%v", xmlObject)
	fmt.Println("finished")
	defer response.Body.Close()
}
