package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"strings"
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
			"name": "Trains_Check",
			"params": {
				"transport": {
					"original": "trains",
					"resolved": "train"
				},
				"direction": {
					"original": "to",
					"resolved": "to"
				},
				"station": {
					"original": "Kings Cross",
					"resolved": "London Kings Cross"
				}
			},
			"query": "What about trains to Kings Cross?"
		},
		"scene": {
			"name": "actions.scene.START_CONVERSATION",
			"slotFillingStatus": "UNSPECIFIED",
			"slots": {}
		},
		"session": {
			"id": "ABwppHEu07ND12dZraO-vgejbGTQ1SIXAoEh_gXi0g1TLn2ZXsaK6E5k31f5E21AE0nzd-SJ-ca0-3v4cURmSVPUlsr0kw",
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
			"lastSeenTime": "2020-11-13T18:45:24Z"
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
	fmt.Println("Request Sucessfully unmarshalled")
	fmt.Println("Unmarshalled object:", requestFromGoogle)
	responseToGoogle.Session.ID = requestFromGoogle.Session.ID
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
