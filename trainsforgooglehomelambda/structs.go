package main

import "encoding/xml"

// Request Google Structure
type requestGoogleHome struct {
	Handler gHandler `json:"handler"`
	Intent  gIntent  `json:"intent"`
	Scene   gScene   `json:"scene"`
	Session gSession `json:"session"`
	User    gUser    `json:"user"`
	Home    gHome    `json:"home"`
	Device  gDevice  `json:"device"`
	Context gContext `json:"context"`
}
type gHandler struct {
	Name string `json:"name"`
}

type gIntent struct {
	Name   string        `json:"name"`
	Params gIntentParams `json:"params"`
	Query  string        `json:"query"`
}
type gIntentParams struct {
	Direction gIntentParameterValue `json:"direction"`
	Station   gIntentParameterValue `json:"station"`
	Transport gIntentParameterValue `json:"transport"`
}

type gIntentParameterValue struct {
	Original string `json:"original"`
	Resolved string `json:"resolved"`
}

type gScene struct {
	Name              string     `json:"name"`
	SlotFillingStatus string     `json:"slotFillingStatus"`
	Slots             gSlot      `json:"slots"`
	Next              gNextScene `json:"next"`
}

// Need to be validated. Leaving with fake "name" property for now.
// Details: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Scene
type gSlot struct {
	Name string `json:"name"`
}

// Need to be validated. Leaving with fake "name" property for now.
// Details: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Scene
type gNextScene struct {
	Name string `json:"name"`
}

type gSession struct {
	ID            string         `json:"id"`
	Params        gSessionParams `json:"params"`
	TypeOverrides gTypeOverride  `json:"typeOverrides"`
	LanguageCode  string         `json:"languageCode"`
}

// Need to be validated. Leaving with fake "name" property for now.
// Details https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Session
type gSessionParams struct {
	Name string `json:"name"`
}

type gTypeOverride struct {
	Name    string `json:"name"`
	Mode    string `json:"mode"`
	Synonym string `json:"synonym"` // Need to be Updated. https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#SynonymType
}

type gUser struct {
	Locale               string                   `json:"locale"`
	Params               gUserParams              `json:"params"`
	AccountLinkingStatus string                   `json:"accountLinkingStatus"`
	VerificationStatus   string                   `json:"verificationStatus"`
	LastSeenTime         string                   `json:"lastSeenTime"`
	Engagement           gUserEngagement          `json:"engagement"`
	PackageEntitlements  gUserPackageEntitlements `json:"packageEntitlements"`
}

// Need to be validated. Leaving with fake "name" property for now.
// Details https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#User
type gUserParams struct {
	Name string `json:"name"`
}
type gUserEngagement struct {
	PushNotificationIntents string `json:"pushNotificationIntents"` // Need to be Updated. https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Engagement
	DailyUpdateIntents      string `json:"dailyUpdateIntents"`      // Need to be Updated. https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Engagement
}

type gUserPackageEntitlements struct {
	string `json:""`
}

type gHome struct {
	Params gHomeParams `json:"params"`
}
type gHomeParams struct {
	Name string `json:"name"` // Need to be updated. https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Home
}
type gDevice struct {
	Capabilities []string `json:"capabilities"`
}
type gContext struct {
	Media gMediaContext `json:"media"`
}
type gMediaContext struct {
	Progress string `json:"progress"`
}

// Response Google Structure

type responseGoogleHome struct {
	Prompt   gPrompt   `json:"prompt"`
	Scene    gScene    `json:"scene"`
	Session  gSession  `json:"session"`
	User     gUser     `json:"user"`
	Home     gHome     `json:"home"`
	Device   gDevice   `json:"device"`
	Expected gExpected `json:"expected"`
}

type gPrompt struct {
	Override    bool           `json:"override"`
	FirstSimple gSimple        `json:"firstSimple"`
	Content     gContent       `json:"content"`
	LastSimple  gSimple        `json:"lastSimple"`
	Suggestions []gSuggestions `json:"suggestions"`
	Link        gLink          `json:"link"`
	Canvas      gCanvas        `json:"canvas"`
	OrderUpdate gOrderUpdate   `json:"orderUpdate"`
}

type gSimple struct {
	Speech string `json:"speech"`
	Text   string `json:"text"`
}

type gContent struct {
	Card       gCard       `json:"card"`
	Image      gImage      `json:"image"`
	Table      gTable      `json:"table"`
	Media      gMedia      `json:"media"`
	Collection gCollection `json:"collection"`
	List       gList       `json:"list"`
}

type gCard struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Text      string `json:"text"`
	Image     gImage `json:"image"`
	ImageFill string `json:"imageFill"`
	Button    string `json:"button"` // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Card
}
type gImage struct {
	URL    string `json:"url"`
	Alt    string `json:"alt"`
	Height int    `json:"height"`
	Width  int    `json:"width"`
}

type gTable struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Image     gImage `json:"image"`
	ImageFill string `json:"imageFill"`
	Columns   string `json:"columns"` // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Table
	Rows      string `json:"rows"`    // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Table
	Button    string `json:"button"`  // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Card
}

type gMedia struct {
	MediaType             string   `json:"mediaType"`
	StartOffset           string   `json:"startOffset"`
	OptionalMediaControls []string `json:"optionalMediaControls"`
	MediaObjects          string   `json:"mediaObjects"` // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Media
}
type gCollection struct {
	Title     string `json:"title"`
	Subtitle  string `json:"subtitle"`
	Items     string `json:"items"` // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Collection
	ImageFill string `json:"imageFill"`
}

type gList struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
	Items    string `json:"items"` // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#List
}

type gLink struct {
	Name string   `json:"name"`
	Open gOpenURL `json:"open"`
}

type gOpenURL struct {
	URL  string `json:"url"`
	Hint string `json:"hint"`
}
type gCanvas struct {
	URL         string   `json:"url"`
	Data        []string `json:"data"`
	SuppressMic bool     `json:"suppressMic"`
}
type gOrderUpdate struct { //https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#OrderUpdate
	Type             string `json:"type"`
	Order            string `json:"order"` // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#OrderUpdate
	UpdateMask       string `json:"updateMask"`
	UserNotification string `json:"userNotification"` // Need Updating: https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#OrderUpdate
	Reason           string `json:"reason"`
}

type gExpected struct { //https://developers.google.com/assistant/conversational/reference/rest/v1/TopLevel/fulfill#Expected
	Speech       []string `json:"speech"`
	LanguageCode string   `json:"languageCode"`
}

// Request SOAP Structure
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
type gSuggestions struct {
	Title string `json:"title"`
}

// Response SOAP structure
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

// Generic variables

type appParams struct {
	LdbwsToken       string `json:"LdbwsToken"`
	LdbwsEndpoint    string `json:"Ldbwsendpoint"`
	DestPreposition  string `json:"DestPreposition"`
	SrcPreposition   string `json:"SrcPreposition"`
	DefaultTimeFrame string `json:"DefaultTimeFrame"`
}
