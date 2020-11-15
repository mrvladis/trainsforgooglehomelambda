package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
)

func processGoogleRequest(requestFromGoogle requestGoogleHome) (requestSoapEnv, error) {
	requestSoap := requestTemplate
	if requestStationTo := requestFromGoogle.Intent.Params.StationTo.Resolved; requestStationTo != "" {
		fmt.Println("Target Station code lookup for station:", requestStationTo)
		destinationStation, err := getStation(requestStationTo)
		if err != nil {
			log.Fatal("Failed to get the Destination Station details", err.Error())
			return requestSoap, err
		}
		requestSoap.Body.Ldb.FilterCrs = destinationStation.CRS
		fmt.Println("Target Station lookup for station:", destinationStation.CRS)
	}
	if requestStationFrom := requestFromGoogle.Intent.Params.StationFrom.Resolved; requestStationFrom != "" {
		fmt.Println("Source Station code identified as:", requestStationFrom)
		sourceStation, err := getStation(requestStationFrom)
		if err != nil {
			log.Fatal("Failed to get the Source Station details", err.Error())
			return requestSoap, err
		}
		requestSoap.Body.Ldb.Crs = sourceStation.CRS
		fmt.Println("Source Station code identified as:", sourceStation.CRS)
	}
	if requestTime := requestFromGoogle.Intent.Params.Time.Resolved; requestTime != nil {

	}
	requestSoap.Header.AccessToken.TokenValue = applicationParameters.LdbwsToken
	requestSoap.Body.Ldb.TimeWindow, err = strconv.Atoi(applicationParameters.DefaultTimeFrame)
	if err != nil {
		log.Fatal("Failed to convert applicationParameters.DefaultTimeFrame value to integer ", err.Error())
		return requestSoap, err

	}
	return requestSoap, nil
}

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
