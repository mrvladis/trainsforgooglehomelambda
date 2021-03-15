package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
)

//* Global Variables
var errorLogger = log.New(os.Stderr, "ERROR ", log.Llongfile)

var response, googleHomeMessage, message string
var responseToGoogle responseGoogleHome
var requestFromGoogle requestGoogleHome
var simpleR gSimple
var promptR gPrompt

func main() {
	gRequest := strings.TrimSpace(`{
		"handler": {
			"name": "TrainsCheck"
		},
		"intent": {
			"name": "Initial_Train_Check",
			"params": {
				"transport": {
					"original": "trains",
					"resolved": "train"
				},
				"directionFrom": {
					"original": "from",
					"resolved": "from"
				},
				"stationFrom": {
					"original": "Oakleigh Park",
					"resolved": "Oakleigh Park"
				},
				"stationTo": {
					"original": "Finsbury Park",
					"resolved": "Finsbury Park"
				},
				"directionTo": {
					"original": "to",
					"resolved": "to"
				},
				"time": {
					"original": "Wednesday at 3pm.",
					"resolved": {
						"month": 11,
						"hours": 15,
						"seconds": 0,
						"year": 2020,
						"nanos": 0,
						"minutes": 0,
						"day": 18
					}
				}
			},
			"query": "Tell me about trains from Oakleigh Park to Finsbury Park on Wednesday at 3pm."
		},
		"scene": {
			"name": "actions.scene.START_CONVERSATION",
			"slotFillingStatus": "UNSPECIFIED",
			"slots": {}
		},
		"session": {
			"id": "ABwppHG68SM0v1C4YxBwNSNG4j64PfVhK4EdFd69p4p6dRVN5cy_N1bejgOqZiFGQNPNAJ3yGTT8psYI83L49VvN7OWwzw",
			"params": {},
			"typeOverrides": [],
			"languageCode": ""
		},
		"user": {
			"locale": "en-US",
			"params": {},
			"accountLinkingStatus": "NOT_LINKED",
			"verificationStatus": "VERIFIED",
			"packageEntitlements": [],
			"lastSeenTime": "2020-11-15T20:56:33Z"
		},
		"home": {
			"params": {}
		},
		"device": {
			"capabilities": [
				"SPEECH",
				"RICH_RESPONSE",
				"LONG_FORM_AUDIO"
			]
		}
	}`,
	)

	/* gResponse := strings.TrimSpace(`{
		"session": {
		  "id": "example_session_id",
		  "params": {}
		},
		"prompt": {
		  "override": false,
		  "content": {
			"card": {
			  "title": "Card Title",
			  "subtitle": "Card Subtitle",
			  "text": "Card Content",
			  "image": {
				"alt": "Google Assistant logo",
				"height": 0,
				"url": "https://developers.google.com/assistant/assistant_96.png",
				"width": 0
			  }
			}
		  },
		  "firstSimple": {
			"speech": "This is a card rich response.",
			"text": ""
		  }
		}
	  }
	`,
	) */

	fmt.Println("Getting the Application Parameters")
	fmt.Printf("Request Body: %v", gRequest)
	//Get the Application Parameters

	err := json.Unmarshal([]byte(gRequest), &requestFromGoogle)

	if err != nil {
		fmt.Println("Couldn't unmarshal Google Request")
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("UnmarshalTypeError %v - %v - %v\n", ute.Value, ute.Type, ute.Offset)
		} else {
			fmt.Println("Other error:", err)
		}

	}

	attibutes, err := dynamodbattribute.MarshalMap(requestFromGoogle)
	fmt.Println("att:", attibutes)
	fmt.Println("Request Sucessfully unmarshalled")
	fmt.Println("Unmarshalled object:", requestFromGoogle)
	responseToGoogle.Session.ID = requestFromGoogle.Session.ID

	if requestStationTo := requestFromGoogle.Intent.Params.StationTo.Resolved; requestStationTo != "" {

	}
	if requestStationFrom := requestFromGoogle.Intent.Params.StationFrom.Resolved; requestStationFrom != "" {

	}
	/* 	if requestTime := requestFromGoogle.Intent.Params.Time.Resolved; requestTime.Day != nil {

	   	} */

	googleHomeMessage = fmt.Sprintln("There are currently services scheduled from  within the next  minutes:")

	simpleR.Speech = &googleHomeMessage
	promptR.FirstSimple = &simpleR
	responseToGoogle.Prompt = &promptR

	var buffer bytes.Buffer
	json.NewEncoder(&buffer).Encode(&responseToGoogle)
	reponseToGoogleBody, err := json.Marshal(responseToGoogle)
	if err != nil {
		fmt.Println("Couldn't marshal Google Response")
		if ute, ok := err.(*json.UnmarshalTypeError); ok {
			fmt.Printf("MarshalTypeError %v - %v - %v\n", ute.Value, ute.Type, ute.Offset)
		} else {
			fmt.Println("Other error:", err)
		}

	}
	fmt.Println("Marshaled Json:", string(reponseToGoogleBody))
	fmt.Println("Encoded Json:", buffer.String())

}
