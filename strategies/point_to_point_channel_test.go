package strategies

import (
	"context"
	"errors"
	"testing"

	pb "github.com/buoyantio/bb/gen"
	"github.com/buoyantio/bb/service"
)

func TestPointToPointChannelStrategy(t *testing.T) {
	t.Run("forwards all requests to strategy and returns its response", func(t *testing.T) {
		expectedResponse := &pb.TheResponse{Payload: "1"}
		mockClient := &service.MockClient{IDToReturn: "1", ResponseToReturn: expectedResponse}

		expectedRequest := &pb.TheRequest{
			RequestUID: "expected-req",
		}

		strategy, err := NewPointToPointChannel(&service.Config{}, []service.Server{service.MockServer{}}, []service.Client{mockClient})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		actualResponse, err := strategy.Do(context.TODO(), expectedRequest)
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		actualRequest := mockClient.RequestReceived
		if actualRequest != expectedRequest {
			t.Fatalf("Expected client [%s] to receive request [%v], but got [%v]", mockClient.GetID(), expectedRequest, actualRequest)
		}

		if actualResponse != expectedResponse {
			t.Fatalf("Expected to return response [%v], but got [%v]", expectedRequest, actualRequest)
		}
	})

	t.Run("forwards errors returned by clients", func(t *testing.T) {
		mockClient := &service.MockClient{IDToReturn: "1", ErrorToReturn: errors.New("expected")}

		expectedRequest := &pb.TheRequest{
			RequestUID: "expected-req",
		}

		strategy, err := NewPointToPointChannel(&service.Config{}, []service.Server{service.MockServer{}}, []service.Client{mockClient})
		if err != nil {
			t.Fatalf("Unexpected error: %v", err)
		}

		_, err = strategy.Do(context.TODO(), expectedRequest)
		if err == nil {
			t.Fatalf("Expecting error, got nothing")
		}

		actualRequest := mockClient.RequestReceived
		if actualRequest != expectedRequest {
			t.Fatalf("Expected client [%s] to receive request [%v], but got [%v]", mockClient.GetID(), expectedRequest, actualRequest)
		}
	})
}
