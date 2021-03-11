package main

import (
	"context"
	"encoding/xml"
	"fmt"
	"log"
	"os"
	"strconv"
	"strings"
)

func prepareRequestToNationalRail(ctx context.Context, requestFromGoogle requestGoogleHome) (requestSoapEnv, error) {
	requestSoap := requestTemplate
	fmt.Println("Initial Saop request body")
	output, err := xml.MarshalIndent(requestSoap, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
	fmt.Println()
	fmt.Println("Preparing to process StationTo Parameter, if available.")
	if requestFromGoogle.Intent.Params.StationTo != nil {
		if requestStationTo := requestFromGoogle.Intent.Params.StationTo.Resolved; requestStationTo != "" {
			fmt.Println("Target Station code lookup for station:", requestStationTo)
			destinationStation, err := getStation(ctx, requestStationTo)
			if err != nil {
				log.Fatal("Failed to get the Destination Station details", err.Error())
				return requestSoap, err
			}
			requestSoap.Body.Ldb.FilterCrs = destinationStation.CRS
			fmt.Println("Target Station code identified as:", destinationStation.CRS)
			fmt.Println("Updated Saop request body", requestSoap)
		}
	} else {
		requestSoap.Body.Ldb.FilterCrs = ""
	}
	fmt.Println("Preparing to process StationFrom Parameter, if available.")
	if requestFromGoogle.Intent.Params.StationFrom != nil {
		if requestStationFrom := requestFromGoogle.Intent.Params.StationFrom.Resolved; requestStationFrom != "" {
			fmt.Println("Source Station code lookup for station:", requestStationFrom)
			sourceStation, err := getStation(ctx, requestStationFrom)
			if err != nil {
				log.Fatal("Failed to get the Source Station details", err.Error())
				return requestSoap, err
			}
			requestSoap.Body.Ldb.Crs = sourceStation.CRS
			fmt.Println("Source Station code identified as:", sourceStation.CRS)
		}
	}
	fmt.Println("Preparing to process Time Parameter, if available.")
	if requestFromGoogle.Intent.Params.Time != nil {
		fmt.Println("Seems Time Parameter exist. Let's try to decouple...")
		if requestTime := requestFromGoogle.Intent.Params.Time.Resolved; requestTime != nil {
			fmt.Println("Request time exist")
		}
	}
	requestSoap.Header.AccessToken.TokenValue = applicationParameters.LdbwsToken
	requestSoap.Body.Ldb.TimeWindow, err = strconv.Atoi(applicationParameters.DefaultTimeFrame)
	fmt.Println("Prepared Saop request body")
	output, err = xml.MarshalIndent(requestSoap, "  ", "    ")
	if err != nil {
		fmt.Printf("error: %v\n", err)
	}
	os.Stdout.Write(output)
	fmt.Println()

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
	fmt.Println()
	if trainsCount > 0 {
		currentServices := responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.TrainServices.Service
		fmt.Println("Processing Trains Information")
		googleHomeMessage = fmt.Sprintln("There are currently", trainsCount, "services scheduled from", responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.LocationName, "within the next ", applicationParameters.DefaultTimeFrame, " minutes:")
		for _, trainService := range currentServices {
			if strings.EqualFold(trainService.Etd, "Cancelled") {
				message = fmt.Sprintln(trainService.Std, trainService.Operator, trainService.Destination.Location.LocationName, "service has been", trainService.Etd)
			} else {
				message = fmt.Sprintln(trainService.Std, trainService.Operator, trainService.Destination.Location.LocationName, "service running", trainService.Etd, "formed of", trainService.Length, "coaches.")
			}
			googleHomeMessage += message
		}

	} else {
		googleHomeMessage = fmt.Sprintln("There are currently", trainsCount, "services scheduled from", responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.LocationName, "within the next ", applicationParameters.DefaultTimeFrame, " minutes.")
	}
	if len(responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.NrccMessages.Message) > 0 {
		fmt.Println("There is a message on the board: ", responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.NrccMessages.Message)
		googleHomeMessage += fmt.Sprintln("There is a message on the board: ")
		googleHomeMessage += strings.ReplaceAll(strings.ReplaceAll(responseXMLObject.Body.GetDepBoardWithDetailsResponse.GetStationBoardResult.NrccMessages.Message, "<P>", ""), "</P>", "")
	}
	simpleR.Speech = &googleHomeMessage
	promptR.FirstSimple = &simpleR
	responseToGoogle.Prompt = &promptR
	return responseToGoogle, nil

}
