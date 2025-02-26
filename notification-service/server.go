package main

import (
	"chat-app/notification-service/config"
	"chat-app/notification-service/database"
	"chat-app/notification-service/handlers"
	"chat-app/notification-service/proto"
	"log"
	"net"
	"google.golang.org/grpc/reflection"
	"google.golang.org/grpc"
)

func main() {
	config.LoadEnv()
	// Connect to DB
	database.Connect()
	database.DB.AutoMigrate(&database.Notification{})

	// Start gRPC server
	server := grpc.NewServer()
	proto.RegisterNotificationServiceServer(server, &handlers.NotificationServiceServer{})

	// Enable gRPC reflection
    reflection.Register(server)

	listener, err := net.Listen("tcp", ":50052")
	if err != nil {
		log.Fatalf("❌ Failed to listen on port 50052: %v", err)
	}

	log.Println("✅ Notification gRPC server started on :50052")
	if err := server.Serve(listener); err != nil {
		log.Fatalf("❌ Failed to start gRPC server: %v", err)
	}
}