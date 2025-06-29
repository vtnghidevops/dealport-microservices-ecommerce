package client

import (
	"context"
	"fmt"
	"log"
	"time"

	pb "mail-service/proto"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// MailServiceClient is a client for the mail service gRPC API
type MailServiceClient struct {
	conn   *grpc.ClientConn
	client pb.MailServiceClient
	logger *log.Logger
}

// NewMailServiceClient creates a new mail service client
func NewMailServiceClient(host string, logger *log.Logger) (*MailServiceClient, error) {
	// Set default logger if none provided
	if logger == nil {
		logger = log.New(log.Writer(), "[MAIL-CLIENT] ", log.LstdFlags)
	}

	// Connect to the gRPC server
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		host,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		return nil, fmt.Errorf("failed to connect to mail service: %w", err)
	}

	// Create the client
	client := pb.NewMailServiceClient(conn)

	return &MailServiceClient{
		conn:   conn,
		client: client,
		logger: logger,
	}, nil
}

// Close closes the client connection
func (c *MailServiceClient) Close() error {
	return c.conn.Close()
}

// SendEmail sends a simple email with basic information
func (c *MailServiceClient) SendEmail(ctx context.Context, from, fromName, to, subject, message string) error {
	req := &pb.SendEmailRequest{
		From:     from,
		FromName: fromName,
		To:       to,
		Subject:  subject,
		Message:  message,
	}

	resp, err := c.client.SendEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling SendEmail: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("mail service reported error: %s", resp.Error)
	}

	c.logger.Printf("Email sent successfully to %s", to)
	return nil
}

// SendTemplateEmail sends an email using a template
func (c *MailServiceClient) SendTemplateEmail(
	ctx context.Context,
	from, fromName, to, subject, template string,
	variables map[string]string,
) error {
	req := &pb.SendTemplateEmailRequest{
		From:      from,
		FromName:  fromName,
		To:        to,
		Subject:   subject,
		Template:  template,
		Variables: variables,
	}

	resp, err := c.client.SendTemplateEmail(ctx, req)
	if err != nil {
		return fmt.Errorf("error calling SendTemplateEmail: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("mail service reported error: %s", resp.Error)
	}

	c.logger.Printf("Template email sent successfully to %s", to)
	return nil
}

// SendOTPEmail sends an OTP verification email
func (c *MailServiceClient) SendOTPEmail(
	ctx context.Context,
	email, otp, purpose string,
	expiresIn int32,
	message, actionText string,
) error {
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
		return fmt.Errorf("error calling SendOTPEmail: %w", err)
	}

	if !resp.Success {
		return fmt.Errorf("mail service reported error: %s", resp.Error)
	}

	c.logger.Printf("OTP email sent successfully to %s", email)
	return nil
}

// SendWelcomeEmail is a helper function to send a welcome email to new users
func (c *MailServiceClient) SendWelcomeEmail(ctx context.Context, to, firstName, lastName string) error {
	variables := map[string]string{
		"first_name": firstName,
		"last_name":  lastName,
		"email":      to,
		"year":       fmt.Sprintf("%d", time.Now().Year()),
	}

	return c.SendTemplateEmail(
		ctx,
		"",
		"",
		to,
		"Welcome to our platform!",
		"welcome.html.gohtml",
		variables,
	)
}

// SendPasswordResetEmail is a helper function to send a password reset email
func (c *MailServiceClient) SendPasswordResetEmail(ctx context.Context, to, tokenHash, expiresAt string) error {
	variables := map[string]string{
		"email":      to,
		"token_hash": tokenHash,
		"expires_at": expiresAt,
		"year":       fmt.Sprintf("%d", time.Now().Year()),
	}

	return c.SendTemplateEmail(
		ctx,
		"",
		"",
		to,
		"Password Reset Request",
		"password_reset.html.gohtml",
		variables,
	)
}

// SendOrderConfirmationEmail is a helper function to send an order confirmation email
func (c *MailServiceClient) SendOrderConfirmationEmail(
	ctx context.Context,
	to, orderID, orderNumber, customerName, orderDate, status,
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
		ctx,
		"",
		"",
		to,
		"Order Confirmation - #"+orderNumber,
		"order_confirmation.html.gohtml",
		variables,
	)
}
