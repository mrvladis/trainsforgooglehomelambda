package main

import (
	"reflect"
	"testing"

	"github.com/aws/aws-lambda-go/events"
)

func Test_router(t *testing.T) {
	type args struct {
		req events.APIGatewayProxyRequest
	}
	tests := []struct {
		name    string
		args    args
		want    events.APIGatewayProxyResponse
		wantErr bool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := router(tt.args.req)
			if (err != nil) != tt.wantErr {
				t.Errorf("router() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("router() = %v, want %v", got, tt.want)
			}
		})
	}
}
