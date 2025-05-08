package grpc

import (
	"context"
	"fmt"
	"log"
	"time"

	"mail-service/internal/config"
	"mail-service/internal/mailer"
	pb "mail-service/proto"
)

// MailServer implements the gRPC MailService interface
type MailServer struct {
	pb.UnimplementedMailServiceServer
	Config config.Config
	logger *log.Logger
}

// NewMailServer creates a new MailServer instance
func NewMailServer(cfg config.Config, logger *log.Logger) *MailServer {
	return &MailServer{
		Config: cfg,
		logger: logger,
	}
}

// SendEmail sends a simple email with basic information
func (s *MailServer) SendEmail(ctx context.Context, req *pb.SendEmailRequest) (*pb.SendEmailResponse, error) {
	s.logger.Printf("gRPC SendEmail request received: To=%s, Subject=%s", req.To, req.Subject)

	// Create the message
	msg := mailer.Message{
		From:     req.From,
		FromName: req.FromName,
		To:       req.To,
		Subject:  req.Subject,
		Data: map[string]interface{}{
			"message": req.Message,
		},
		Template: "mail.html.gohtml", // Use default template
	}

	// Send the email
	err := s.Config.Mailer.SendSMTPMessage(msg)
	if err != nil {
		s.logger.Printf("Error sending email: %v", err)
		return &pb.SendEmailResponse{
			Success: false,
			Message: "Failed to send email",
			Error:   err.Error(),
		}, nil
	}

	s.logger.Printf("Email sent successfully to %s", req.To)
	return &pb.SendEmailResponse{
		Success: true,
		Message: "Email sent successfully",
	}, nil
}

// SendTemplateEmail sends an email using a template
func (s *MailServer) SendTemplateEmail(ctx context.Context, req *pb.SendTemplateEmailRequest) (*pb.SendEmailResponse, error) {
	s.logger.Printf("gRPC SendTemplateEmail request received: To=%s, Subject=%s, Template=%s",
		req.To, req.Subject, req.Template)

	// Convert variables to map[string]interface{}
	variables := make(map[string]interface{})
	for k, v := range req.Variables {
		variables[k] = v
	}

	// Create the message
	msg := mailer.Message{
		From:     req.From,
		FromName: req.FromName,
		To:       req.To,
		Subject:  req.Subject,
		Template: req.Template,
		Data:     variables,
	}

	// Send the email
	err := s.Config.Mailer.SendSMTPMessage(msg)
	if err != nil {
		s.logger.Printf("Error sending template email: %v", err)
		return &pb.SendEmailResponse{
			Success: false,
			Message: "Failed to send email",
			Error:   err.Error(),
		}, nil
	}

	s.logger.Printf("Template email sent successfully to %s", req.To)
	return &pb.SendEmailResponse{
		Success: true,
		Message: "Email sent successfully",
	}, nil
}

// SendOTPEmail sends an OTP verification email
func (s *MailServer) SendOTPEmail(ctx context.Context, req *pb.SendOTPEmailRequest) (*pb.SendEmailResponse, error) {
	s.logger.Printf("gRPC SendOTPEmail request received: To=%s, Purpose=%s", req.Email, req.Purpose)

	// Set default expiration time if not provided
	expiresIn := req.ExpiresIn
	if expiresIn == 0 {
		expiresIn = 15 // Default: 15 minutes
	}

	// Set subject based on purpose
	subject := "Your Verification Code"
	if req.Purpose == "password_reset" {
		subject = "Password Reset Verification"
	} else if req.Purpose == "registration" {
		subject = "Complete Your Registration"
	}

	// Prepare variables for template
	variables := map[string]interface{}{
		"email":       req.Email,
		"otp":         req.Otp,
		"purpose":     req.Purpose,
		"expires_in":  fmt.Sprintf("%d", expiresIn),
		"message":     req.Message,
		"action_text": req.ActionText,
		"year":        time.Now().Year(),
	}

	// Create the message
	msg := mailer.Message{
		To:       req.Email,
		Subject:  subject,
		Template: "otp_verification.html.gohtml",
		Data:     variables,
	}

	// Send the email
	err := s.Config.Mailer.SendSMTPMessage(msg)
	if err != nil {
		s.logger.Printf("Error sending OTP email: %v", err)
		return &pb.SendEmailResponse{
			Success: false,
			Message: "Failed to send OTP email",
			Error:   err.Error(),
		}, nil
	}

	s.logger.Printf("OTP email sent successfully to %s", req.Email)
	return &pb.SendEmailResponse{
		Success: true,
		Message: "OTP email sent successfully",
	}, nil
}
