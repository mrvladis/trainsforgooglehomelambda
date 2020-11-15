package main

import (
	"fmt"
	"strings"
)

func prepareGoogleResponse(responseXMLObject *responseSoapEnv) (responseGoogleHome, error) {
	var googleHomeMessage, message string
	var simpleR gSimple
	var promptR gPrompt
	var responseToGoogle responseGoogleHome

	fmt.Println("Preparing Result")
	trainsCount := len(responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.TrainServices.Service)
	fmt.Printf("There are %v trains", trainsCount)
	currentServices := responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.TrainServices.Service
	fmt.Printf("There are %v trains", trainsCount)
	fmt.Printf("Processing Trains Information")
	if trainsCount > 0 {
		googleHomeMessage = fmt.Sprintln("There are currently", trainsCount, "services scheduled from", responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.LocationName, "within the next ", applicationParameters.DefaultTimeFrame, " minutes:")
		for _, trainService := range currentServices {
			if strings.EqualFold(trainService.Etd, "Cancelled") {
				message = fmt.Sprintln(trainService.Std, trainService.Operator, trainService.Destination.Location.LocationName, "service has been", trainService.Etd)
			} else {
				message = fmt.Sprintln(trainService.Std, trainService.Operator, trainService.Destination.Location.LocationName, "service running", trainService.Etd, "formed of", trainService.Length, "coaches.")
			}
			googleHomeMessage += message
		}
	}

	simpleR.Speech = &googleHomeMessage
	promptR.FirstSimple = &simpleR
	responseToGoogle.Prompt = &promptR
	return responseToGoogle, nil

}
