package main

import (
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

	fmt.Println("Getting the Application Parameters")
	fmt.Printf("Request Body: %v", gRequest)
	//Get the Application Parameters

	err := json.Unmarshal([]byte(gRequest), &requestFromGoogle)
	if err != nil {
		fmt.Println("Couldn't unmarshal Google Request")

	}

}
