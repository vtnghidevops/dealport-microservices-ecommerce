package event

import (
	"context"
	"fmt"
	"log"
	"os"
	"time"

	pb "listener-service/proto/mail"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MailClient handles interactions with the mail service
type MailClient struct {
	conn   *grpc.ClientConn
	client pb.MailServiceClient
	logger *log.Logger
}

// NewMailClient creates a new mail service client
func NewMailClient(host string) (*MailClient, error) {
	logger := log.New(os.Stdout, "[MAIL-CLIENT] ", log.LstdFlags)

	// If no host is provided, use the default
	if host == "" {
		host = "mail-service:50057"
	}

	// Connect to the mail service
	logger.Printf("Connecting to mail service at %s", host)
	conn, err := grpc.Dial(
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mail service: %w", err)
	}

	// Create the client
	client := pb.NewMailServiceClient(conn)

	return &MailClient{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the gRPC connection
func (c *MailClient) Close() error {
	return c.conn.Close()
}

// SendEmail sends a simple email with basic information
func (c *MailClient) SendEmail(email, subject, message string) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SendEmailRequest{
		To:      email,
		Subject: subject,
		Message: message,
	}

	resp, err := c.client.SendEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling mail service: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("mail service error: %s", resp.Error)
	}

	return nil
}

// SendTemplateEmail sends an email using a template
func (c *MailClient) SendTemplateEmail(
	email, subject, template string,
	variables map[string]string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SendTemplateEmailRequest{
		To:        email,
		Subject:   subject,
		Template:  template,
		Variables: variables,
	}

	resp, err := c.client.SendTemplateEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling mail service: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("mail service error: %s", resp.Error)
	}

	return nil
}

// SendOTPEmail sends an OTP verification email
func (c *MailClient) SendOTPEmail(
	email, otp, purpose string,
	expiresIn int32,
	message, actionText string,
) error {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	req := &pb.SendOTPEmailRequest{
		Email:      email,
		Otp:        otp,
		Purpose:    purpose,
		ExpiresIn:  expiresIn,
		Message:    message,
		ActionText: actionText,
	}

	resp, err := c.client.SendOTPEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling mail service: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("mail service error: %s", resp.Error)
	}

	return nil
}

// SendWelcomeEmail sends a welcome email to new users
func (c *MailClient) SendWelcomeEmail(email, firstName, lastName string) error {
	variables := map[string]string{
		"first_name": firstName,
		"last_name":  lastName,
		"email":      email,
		"year":       fmt.Sprintf("%d", time.Now().Year()),
	}

	return c.SendTemplateEmail(
		email,
		"Welcome to our platform!",
		"welcome.html.gohtml",
		variables,
	)
}

// SendPasswordResetEmail sends a password reset email
func (c *MailClient) SendPasswordResetEmail(email, tokenHash, expiresAt string) error {
	variables := map[string]string{
		"email":      email,
		"token_hash": tokenHash,
		"expires_at": expiresAt,
		"year":       fmt.Sprintf("%d", time.Now().Year()),
	}

	return c.SendTemplateEmail(
		email,
		"Password Reset Request",
		"password_reset.html.gohtml",
		variables,
	)
}

// SendPasswordChangedEmail sends a password changed notification email
func (c *MailClient) SendPasswordChangedEmail(email, changedAt string) error {
	variables := map[string]string{
		"email":      email,
		"changed_at": changedAt,
		"year":       fmt.Sprintf("%d", time.Now().Year()),
	}

	return c.SendTemplateEmail(
		email,
		"Your Password Has Been Changed",
		"password_changed.html.gohtml",
		variables,
	)
}

// SendOrderConfirmationEmail sends an order confirmation email
func (c *MailClient) SendOrderConfirmationEmail(
	email, orderID, orderNumber, customerName, orderDate, status,
	paymentMethod, total, items, shippingAddress, shippingName, shippingPhone string,
) error {
	variables := map[string]string{
		"customer_name":    customerName,
		"order_id":         orderID,
		"order_number":     orderNumber,
		"order_date":       orderDate,
		"status":           status,
		"payment_method":   paymentMethod,
		"total":            total,
		"items":            items,
		"shipping_name":    shippingName,
		"shipping_address": shippingAddress,
		"shipping_phone":   shippingPhone,
		"year":             fmt.Sprintf("%d", time.Now().Year()),
	}

	return c.SendTemplateEmail(
		email,
		"Order Confirmation - #"+orderNumber,
		"order_confirmation.html.gohtml",
		variables,
	)
}
