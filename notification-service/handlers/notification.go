package handlers

import (
	"chat-app/notification-service/database"
	"context"
	"log"
	"time"

	pb "chat-app/notification-service/proto"
)

type NotificationServiceServer struct {
	pb.UnimplementedNotificationServiceServer
}

// SendNotification saves a notification in the database
func (s *NotificationServiceServer) SendNotification(ctx context.Context, req *pb.SendNotificationRequest) (*pb.SendNotificationResponse, error) {
	log.Printf("📩 Sending notification to UserID=%s: %s", req.UserId, req.Message)

	notification := database.Notification{
		UserID:    req.UserId,
		Message:   req.Message,
		Timestamp: time.Now().Format(time.RFC3339),
	}

	if err := database.DB.Create(&notification).Error; err != nil {
		log.Printf("❌ Error saving notification: %v", err)
		return &pb.SendNotificationResponse{
			Success: false,
			Message: "Failed to send notification",
		}, err
	}

	log.Println("✅ Notification saved successfully")
	return &pb.SendNotificationResponse{
		Success: true,
		Message: "Notification sent successfully",
	}, nil
}

// GetNotifications retrieves notifications for a user
func (s *NotificationServiceServer) GetNotifications(ctx context.Context, req *pb.GetNotificationsRequest) (*pb.GetNotificationsResponse, error) {
	log.Printf("🔎 Fetching notifications for UserID=%s", req.UserId)

	var notifications []database.Notification
	if err := database.DB.Where("user_id = ?", req.UserId).Find(&notifications).Error; err != nil {
		log.Printf("❌ Error fetching notifications: %v", err)
		return nil, err
	}

	// Convert DB notifications to gRPC response format
	var grpcNotifications []*pb.Notification
	for _, n := range notifications {
		grpcNotifications = append(grpcNotifications, &pb.Notification{
			UserId:    n.UserID,
			Message:   n.Message,
			Timestamp: n.Timestamp,
		})
	}

	log.Printf("✅ Found %d notifications", len(grpcNotifications))
	return &pb.GetNotificationsResponse{Notifications: grpcNotifications}, nil
}