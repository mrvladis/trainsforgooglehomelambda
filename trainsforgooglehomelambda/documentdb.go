package main

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
	"github.com/aws/aws-sdk-go/service/dynamodb/dynamodbattribute"
	"github.com/aws/aws-xray-sdk-go/xray"
)

var db = dynamodb.New(session.New(), aws.NewConfig().WithRegion(awsRegion))

func scanStationCodes(ctx context.Context) (*[]appStation, error) {
	fmt.Println("Preparign request for the StationCodes.")
	xray.AWS(db.Client)
	input := &dynamodb.ScanInput{
		TableName: aws.String(applicationParameters.StationCodesStore),
	}
	fmt.Println("Sending the request for the StationCodes")
	result, err := db.ScanWithContext(ctx, input)
	if err != nil {
		if aerr, ok := err.(awserr.Error); ok {
			switch aerr.Code() {
			case dynamodb.ErrCodeProvisionedThroughputExceededException:
				fmt.Println(dynamodb.ErrCodeProvisionedThroughputExceededException, aerr.Error())
			case dynamodb.ErrCodeResourceNotFoundException:
				fmt.Println(dynamodb.ErrCodeResourceNotFoundException, aerr.Error())
			case dynamodb.ErrCodeRequestLimitExceeded:
				fmt.Println(dynamodb.ErrCodeRequestLimitExceeded, aerr.Error())
			case dynamodb.ErrCodeInternalServerError:
				fmt.Println(dynamodb.ErrCodeInternalServerError, aerr.Error())
			default:
				fmt.Println(aerr.Error())
			}
		} else {
			// Print the error, cast err to awserr.Error to get the Code and
			// Message from an error.
			fmt.Println(err.Error())
		}
		return nil, nil
	}
	if result.Items == nil {
		return nil, nil
	}

	hab := new([]appStation)
	fmt.Println("Unmarshaling Result Map Sending the request for the StationCodes")
	err = dynamodbattribute.UnmarshalListOfMaps(result.Items, hab)
	if err != nil {
		fmt.Printf("Erorr: %s", err)
		return nil, err
	}

	return hab, nil
}

func getStation(ctx context.Context, stationName string) (*appStation, error) {
	fmt.Println("Preparign request for the Station.")
	xray.AWS(db.Client)
	input := &dynamodb.GetItemInput{
		TableName: aws.String(applicationParameters.StationCodesStore),
		Key: map[string]*dynamodb.AttributeValue{
			"StationName": {
				S: aws.String(stationName),
			},
		},
	}
	fmt.Println("Sending the request for the Station")
	result, err := db.GetItemWithContext(ctx, input)
	if err != nil {
		return nil, err
	}
	if result.Item == nil {
		return nil, nil
	}

	station := new(appStation)
	fmt.Println("Unmarshaling Result Map Sending the request for the Station")
	err = dynamodbattribute.UnmarshalMap(result.Item, station)
	if err != nil {
		fmt.Printf("Erorr: %s", err)
		return nil, err
	}

	return station, nil
}

/*
// Add a habbit record to DynamoDB.
func putItem(hab *habbit) error {
	input := &dynamodb.PutItemInput{
		TableName: aws.String("habbits"),
		Item: map[string]*dynamodb.AttributeValue{
			"id": {
				S: aws.String(hab.ID),
			},
			"habitName": {
				S: aws.String(hab.HabitName),
			},
			"habitCategory": {
				S: aws.String(hab.HabitCategory),
			},
		},
	}

	_, err := db.PutItem(input)
	return err
}
*/
