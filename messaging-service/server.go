package main

import (
	"fmt"
	"log"
	"net"

	"chat-app/messaging-service/config"
	"chat-app/messaging-service/database"
	"chat-app/messaging-service/handlers"
	"chat-app/messaging-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	// Load environment variables
	config.LoadEnv()

	// Connect to the database
	database.Connect()

	// Start gRPC server
	listener, err := net.Listen("tcp", ":50051")
	if err != nil {
		log.Fatalf("Failed to listen on port 50051: %v", err)
	}

	grpcServer := grpc.NewServer()
	proto.RegisterMessagingServiceServer(grpcServer, &handlers.MessagingServiceServer{})

	// Enable gRPC reflection for testing with grpcurl
	reflection.Register(grpcServer)

	fmt.Println("Messaging gRPC server started on :50051")

	if err := grpcServer.Serve(listener); err != nil {
		log.Fatalf("Failed to start gRPC server: %v", err)
	}
}
