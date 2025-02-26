package handlers

import (
	"chat-app/messaging-service/database"

	"chat-app/messaging-service/proto"
	notificationpb "chat-app/notification-service/proto"


	"context"
	"fmt"
	"log"
	"time"

	"google.golang.org/grpc"
)

type MessagingServiceServer struct {
	proto.UnimplementedMessagingServiceServer
}

// SendMessage stores a message in the database and notifies the receiver.
func (s *MessagingServiceServer) SendMessage(ctx context.Context, req *proto.SendMessageRequest) (*proto.SendMessageResponse, error) {
	message := database.Message{
		SenderID:   req.SenderId,
		ReceiverID: req.ReceiverId,
		Content:    req.Content,
		Timestamp:  time.Now().Format(time.RFC3339),
	}

	// Save message in DB
	if err := database.DB.Create(&message).Error; err != nil {
		log.Println("❌ Failed to save message:", err)
		return &proto.SendMessageResponse{
			Success: false,
			Message: "Failed to save message",
		}, err
	}

	log.Println("✅ Message saved successfully:", message)

	// Send Notification via gRPC
	conn, err := grpc.Dial("localhost:50052", grpc.WithInsecure()) // Connect to notification-service
	if err != nil {
		log.Println("❌ Failed to connect to NotificationService:", err)
		return nil, err
	}
	defer conn.Close()

	notificationClient := notificationpb.NewNotificationServiceClient(conn)
	_, err = notificationClient.SendNotification(ctx, &notificationpb.SendNotificationRequest{
		UserId:  req.ReceiverId,
		Message: fmt.Sprintf("New message from User %s", req.SenderId),
	})
	if err != nil {
		log.Println("❌ Failed to send notification:", err)
	} else {
		log.Println("📩 Notification sent successfully to User", req.ReceiverId)
	}

	return &proto.SendMessageResponse{
		Success: true,
		Message: "Message sent successfully",
	}, nil
}

// GetMessages retrieves messages between two users.
func (s *MessagingServiceServer) GetMessages(ctx context.Context, req *proto.GetMessagesRequest) (*proto.GetMessagesResponse, error) {
	var messages []database.Message

	// Query messages from the database
	err := database.DB.Where(
		"(sender_id = ? AND receiver_id = ?) OR (sender_id = ? AND receiver_id = ?)",
		req.SenderId, req.ReceiverId, req.ReceiverId, req.SenderId,
	).Find(&messages).Error
	if err != nil {
		log.Println("❌ Failed to retrieve messages:", err)
		return nil, err
	}

	// Convert to gRPC response format
	var responseMessages []*proto.Message
	for _, msg := range messages {
		responseMessages = append(responseMessages, &proto.Message{
			SenderId:   msg.SenderID,
			ReceiverId: msg.ReceiverID,
			Content:    msg.Content,
			Timestamp:  msg.Timestamp,
		})
	}

	log.Printf("🔎 Retrieved %d messages between User %s and User %s\n", len(responseMessages), req.SenderId, req.ReceiverId)

	return &proto.GetMessagesResponse{Messages: responseMessages}, nil
}