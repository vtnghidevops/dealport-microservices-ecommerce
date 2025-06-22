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

// EmailData represents data for email operations
type EmailData struct {
	Type      string            `json:"type"`
	From      string            `json:"from,omitempty"`
	FromName  string            `json:"from_name,omitempty"`
	To        string            `json:"to"`
	Subject   string            `json:"subject"`
	Template  string            `json:"template,omitempty"`
	Variables map[string]string `json:"variables,omitempty"`
}

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

	// Create data map for handleGenericEmail
	data := map[string]string{
		"from":      req.From,
		"from_name": req.FromName,
		"to":        req.To,
		"subject":   req.Subject,
		"message":   req.Message,
	}

	// Send the email using the helper function
	err := s.handleGenericEmail(data)
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

	// Check if we need to use a specific helper based on template or variables
	emailType := ""
	if req.Variables != nil {
		if typeVal, ok := req.Variables["email_type"]; ok {
			emailType = typeVal
		}
	}

	// Determine which helper function to use based on email type or template
	var err error
	switch {
	case emailType == "registration" || req.Template == "welcome.html.gohtml" || req.Template == "register.html.gohtml":
		err = s.handleRegistrationEmail(req.To, req.Variables)
	case emailType == "password_reset" || req.Template == "password_reset.html.gohtml" || req.Template == "reset_password.html.gohtml":
		err = s.handlePasswordResetEmail(req.To, req.Variables)
	case emailType == "password_change" || emailType == "password_changed" || req.Template == "password_changed.html.gohtml" || req.Template == "password_change.html.gohtml":
		err = s.handlePasswordChangeEmail(req.To, req.Variables)
	case emailType == "order_confirmation" || req.Template == "order_confirmation.html.gohtml":
		err = s.handleOrderConfirmationEmail(req.To, req.Variables)
	case emailType == "order_processing" || req.Template == "order_processing.html.gohtml":
		err = s.handleOrderProcessingEmail(req.To, req.Variables)
	case emailType == "order_shipped" || req.Template == "order_shipped.html.gohtml":
		err = s.handleOrderShippedEmail(req.To, req.Variables)
	case emailType == "order_delivered" || req.Template == "order_delivered.html.gohtml":
		err = s.handleOrderDeliveredEmail(req.To, req.Variables)
	case emailType == "order_cancelled" || req.Template == "order_cancelled.html.gohtml":
		err = s.handleOrderCancelledEmail(req.To, req.Variables)
	case emailType == "order_status_change" || req.Template == "order_status_change.html.gohtml":
		err = s.handleOrderStatusChangeEmail(req.To, req.Variables)
	case emailType == "payment_success" || req.Template == "payment_success.html.gohtml":
		err = s.handlePaymentSuccessEmail(req.To, req.Variables)
	default:
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
		err = s.Config.Mailer.SendSMTPMessage(msg)
	}

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

	// Use the helper function to send OTP email
	err := s.sendOTPEmail(req.Email, req.Otp, req.Purpose, int(req.ExpiresIn), req.Message, req.ActionText)
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

// sendOTPEmail sends an OTP verification email
func (s *MailServer) sendOTPEmail(email, otp, purpose string, expiresIn int, message, actionText string) error {
	// Set default expiration time if not provided
	if expiresIn == 0 {
		expiresIn = 15 // Default: 15 minutes
	}

	// Prepare variables for template
	variables := map[string]interface{}{
		"email":       email,
		"otp":         otp,
		"purpose":     purpose,
		"expires_in":  fmt.Sprintf("%d", expiresIn),
		"message":     message,
		"action_text": actionText,
		"year":        time.Now().Year(),
	}

	// Set subject based on purpose
	subject := "Your Verification Code"
	if purpose == "password_reset" {
		subject = "Password Reset Verification"
	} else if purpose == "registration" {
		subject = "Complete Your Registration"
	}

	// Create the message
	msg := mailer.Message{
		To:       email,
		Subject:  subject,
		Template: "otp_verification.html.gohtml",
		Data:     variables,
	}

	// Send email
	err := s.Config.Mailer.SendSMTPMessage(msg)
	if err != nil {
		s.logger.Printf("Error sending OTP email: %v", err)
		return err
	}

	s.logger.Printf("OTP email sent successfully to %s (purpose: %s)", email, purpose)
	return nil
}

// handleGenericEmail processes a generic email request
func (s *MailServer) handleGenericEmail(data map[string]string) error {
	// Set default template if not specified
	template := data["template"]
	if template == "" {
		template = "welcome.html.gohtml"
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       data["to"],
		Subject:  data["subject"],
		Template: template,
		Data:     data,
	}

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handleRegistrationEmail processes registration email
func (s *MailServer) handleRegistrationEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "welcome.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		msg.Subject = "Registration Confirmation"
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending registration email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handlePasswordResetEmail processes password reset email
func (s *MailServer) handlePasswordResetEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "reset_password.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		msg.Subject = "Password Reset Request"
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending password reset email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handlePasswordChangeEmail processes password change email
func (s *MailServer) handlePasswordChangeEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "password_changed.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		msg.Subject = "Password Changed"
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending password change email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handleOrderConfirmationEmail processes order confirmation email
func (s *MailServer) handleOrderConfirmationEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "order_confirmation.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		// Use order number in subject if available
		if orderNumber, ok := data["order_number"]; ok && orderNumber != "" {
			msg.Subject = fmt.Sprintf("Order Confirmation - #%s", orderNumber)
		} else {
			msg.Subject = "Order Confirmation"
		}
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending order confirmation email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handleOrderProcessingEmail processes order processing notification email
func (s *MailServer) handleOrderProcessingEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "order_processing.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		// Use order number in subject if available
		if orderNumber, ok := data["order_number"]; ok && orderNumber != "" {
			msg.Subject = fmt.Sprintf("Your Order #%s is Being Processed", orderNumber)
		} else {
			msg.Subject = "Your Order is Being Processed"
		}
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending order processing email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handleOrderShippedEmail processes order shipped notification email
func (s *MailServer) handleOrderShippedEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "order_shipped.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		// Use order number in subject if available
		if orderNumber, ok := data["order_number"]; ok && orderNumber != "" {
			msg.Subject = fmt.Sprintf("Your Order #%s Has Been Shipped", orderNumber)
		} else {
			msg.Subject = "Your Order Has Been Shipped"
		}
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending order shipped email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handleOrderDeliveredEmail processes order delivered notification email
func (s *MailServer) handleOrderDeliveredEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "order_delivered.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		// Use order number in subject if available
		if orderNumber, ok := data["order_number"]; ok && orderNumber != "" {
			msg.Subject = fmt.Sprintf("Your Order #%s Has Been Delivered", orderNumber)
		} else {
			msg.Subject = "Your Order Has Been Delivered"
		}
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending order delivered email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handleOrderCancelledEmail processes order cancelled notification email
func (s *MailServer) handleOrderCancelledEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "order_cancelled.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		// Use order number in subject if available
		if orderNumber, ok := data["order_number"]; ok && orderNumber != "" {
			msg.Subject = fmt.Sprintf("Your Order #%s Has Been Cancelled", orderNumber)
		} else {
			msg.Subject = "Your Order Has Been Cancelled"
		}
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending order cancelled email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handleOrderStatusChangeEmail processes generic order status change notification email
func (s *MailServer) handleOrderStatusChangeEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "order_status_change.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		// Use order number in subject if available
		if orderNumber, ok := data["order_number"]; ok && orderNumber != "" {
			if status, ok := data["status"]; ok && status != "" {
				msg.Subject = fmt.Sprintf("Order #%s Status Update: %s", orderNumber, status)
			} else {
				msg.Subject = fmt.Sprintf("Order #%s Status Update", orderNumber)
			}
		} else {
			msg.Subject = "Order Status Update"
		}
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending order status change email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}

// handlePaymentSuccessEmail processes payment success notification email
func (s *MailServer) handlePaymentSuccessEmail(email string, data map[string]string) error {
	// Convert data map to interface map
	dataInterface := make(map[string]interface{})
	for k, v := range data {
		dataInterface[k] = v
	}

	msg := mailer.Message{
		From:     data["from"],
		FromName: data["from_name"],
		To:       email,
		Subject:  data["subject"],
		Template: "payment_success.html.gohtml",
		Data:     dataInterface,
	}

	// Use custom subject if not provided
	if msg.Subject == "" {
		// Use order number in subject if available
		if orderNumber, ok := data["order_number"]; ok && orderNumber != "" {
			msg.Subject = fmt.Sprintf("Payment Confirmation for Order #%s", orderNumber)
		} else {
			msg.Subject = "Payment Confirmation"
		}
	}

	// Override template if specified
	if template := data["template"]; template != "" {
		msg.Template = template
	}

	s.logger.Printf("Sending payment success email to %s", email)

	// Send email
	return s.Config.Mailer.SendSMTPMessage(msg)
}
