package main

import (
	"context"
	"log"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	pb "github.com/troydai/grpcprober/gen/api/protos"
)

func main() {
	// Set up a connection to the gRPC server.
	conn, err := grpc.NewClient("localhost:8080", grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("Failed to connect: %v", err)
	}
	defer conn.Close()

	// Create a new gRPC client for the beacon service.
	client := pb.NewBeaconClient(conn)

	// Make a gRPC request to the beacon service.
	response, err := client.Signal(context.Background(), &pb.SignalRequest{})
	if err != nil {
		log.Fatalf("Failed to get beacon: %v", err)
	}

	// Process the response from the beacon service.
	log.Printf("Received beacon: %v", response.GetDetails())
}
