package handlers

import (
	"context"
	"log"
	"net/http"

	"chat-app/gateway/proto/proto" // Import gRPC-generated files
	"github.com/gofiber/fiber/v2"
	"google.golang.org/grpc"
)

// Messaging Service gRPC Connection
var messagingClient proto.MessagingServiceClient
var notificationClient proto.NotificationServiceClient

func InitGRPCClients() {
	// Connect to messaging-service
	msgConn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		log.Fatal("❌ Failed to connect to messaging-service:", err)
	}
	messagingClient = proto.NewMessagingServiceClient(msgConn)

	// Connect to notification-service
	notifConn, err := grpc.Dial("localhost:50052", grpc.WithInsecure())
	if err != nil {
		log.Fatal("❌ Failed to connect to notification-service:", err)
	}
	notificationClient = proto.NewNotificationServiceClient(notifConn)
}

// Send Message (Calls messaging-service gRPC)
func SendMessage(c *fiber.Ctx) error {
	var req proto.SendMessageRequest
	if err := c.BodyParser(&req); err != nil {
		return c.Status(http.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request body"})
	}

	resp, err := messagingClient.SendMessage(context.Background(), &req)
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to send message"})
	}

	return c.JSON(resp)
}

// Get Messages (Calls messaging-service gRPC)
func GetMessages(c *fiber.Ctx) error {
	senderID := c.Query("sender_id")
	receiverID := c.Query("receiver_id")

	resp, err := messagingClient.GetMessages(context.Background(), &proto.GetMessagesRequest{
		SenderId: senderID, ReceiverId: receiverID,
	})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve messages"})
	}

	return c.JSON(resp)
}

// Get Notifications (Calls notification-service gRPC)
func GetNotifications(c *fiber.Ctx) error {
	userID := c.Query("user_id")

	resp, err := notificationClient.GetNotifications(context.Background(), &proto.GetNotificationsRequest{UserId: userID})
	if err != nil {
		return c.Status(http.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to retrieve notifications"})
	}

	return c.JSON(resp)
}