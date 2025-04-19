// Generates and sends an email using SMTP with both HTML (with inline CSS)
package main

import (
	"bytes"
	"html/template"
	"time"

	"github.com/vanng822/go-premailer/premailer" // Used for inlining CSS in HTML email.
	mail "github.com/xhit/go-simple-mail/v2"     // SMTP email send library.
)

// Mail holds the SMTP configuration and default sender information.
type Mail struct {
	Domain     string // Domain of the mail server.
	Host       string // SMTP server hostname.
	Port       int    // SMTP server port.
	Username   string // Username for SMTP authentication.
	Password   string // Password for SMTP authentication.
	Encryption string // Encryption type: "tls", "ssl", or "none".
	FromAddr   string // Default sender's email address.
	FromName   string // Default sender's name.
}

// Message represents the email content and recipient details.
type Message struct {
	From        string                 // Sender's email address.
	FromName    string                 // Sender's name.
	To          string                 // Recipient's email address.
	Subject     string                 // Email subject line.
	Attachments []string               // List of file paths to attach.
	Data        interface{}            // Data to be used render context in email templates .
	DataMap     map[string]interface{} // Map used for template data.
}

// SendSMTPMessage handles the process of building the email (both HTML and plain text),
// establishing an SMTP connection, and sending the email.
func (m *Mail) SendSMTPMessage(msg Message) error {
	// Use default sender email if not explicitly provided.
	if msg.From == "" {
		msg.From = m.FromAddr
	}

	// Use default sender name if not explicitly provided.
	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	// Create a data map for the templates (the "message" key contains the dynamic data).
	data := map[string]interface{}{
		"message": msg.Data,
	}
	msg.DataMap = data

	/*
		Render a context email from template
	*/
	// Build the HTML email body using the HTML template and inline CSS.
	formattedMessage, err := m.buildHTMLMessage(msg)
	if err != nil {
		return err
	}

	// Build the plain text email body using the plain text template.
	plainMessage, err := m.buildPlainTextMessage(msg)
	if err != nil {
		return err
	}

	// Set up the SMTP client with the configuration stored in the Mail struct.
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption(m.Encryption)
	server.KeepAlive = false // Disable persistent (keep-alive) connections.
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Establish the connection to the SMTP server.
	smtpClient, err := server.Connect()
	if err != nil {
		return err
	}

	// Create a new email message.
	email := mail.NewMSG()
	email.SetFrom(msg.From). // Set the sender's email.
					AddTo(msg.To).          // Set the recipient's email.
					SetSubject(msg.Subject) // Set the email subject.

	// Set the primary body as plain text and add the HTML version as an alternative.
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, formattedMessage)

	// Attach files if any attachments are provided.
	if len(msg.Attachments) > 0 {
		for _, x := range msg.Attachments {
			email.AddAttachment(x)
		}
	}

	// Send the email through the established SMTP connection.
	err = email.Send(smtpClient)
	if err != nil {
		return err
	}

	return nil
}

// buildHTMLMessage generates the HTML version of the email by parsing an HTML
// template and inlining CSS styles to ensure compatibility across different email clients.
func (m *Mail) buildHTMLMessage(msg Message) (string, error) {
	// Specify the path to the HTML email template.
	templateToRender := "/app/templates/mail.html.gohtml"

	// Parse the HTML template file.
	t, err := template.New("email-html").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	// Execute the template's "body" block with the provided data map.
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// Convert the executed template to a string.
	formattedMessage := tpl.String()
	// Inline CSS styles (using the premailer library) to ensure better rendering.
	formattedMessage, err = m.inlineCSS(formattedMessage)
	if err != nil {
		return "", err
	}

	return formattedMessage, nil
}

// buildPlainTextMessage generates the plain text version of the email by parsing
// a corresponding template file.
func (m *Mail) buildPlainTextMessage(msg Message) (string, error) {
	// Specify the path to the plain text email template.
	templateToRender := "/app/templates/mail.plain.gohtml"

	// Parse the plain text template file.
	t, err := template.New("email-plain").ParseFiles(templateToRender)
	if err != nil {
		return "", err
	}

	var tpl bytes.Buffer
	// Execute the template's "body" block with the provided data map.
	if err = t.ExecuteTemplate(&tpl, "body", msg.DataMap); err != nil {
		return "", err
	}

	// Return the generated plain text message.
	plainMessage := tpl.String()

	return plainMessage, nil
}

// inlineCSS processes an HTML string to convert CSS styles to inline styles.
// This is crucial for ensuring that styling is correctly applied on various email clients.
func (m *Mail) inlineCSS(s string) (string, error) {
	// Define options for the premailer processing.
	options := premailer.Options{
		RemoveClasses:     false, // Do not remove CSS class attributes from elements.
		CssToAttributes:   false, // Do not convert CSS properties into HTML attributes.
		KeepBangImportant: true,  // Retain "!important" rules in the CSS.
	}

	// Initialize a new premailer instance with the HTML string and processing options.
	prem, err := premailer.NewPremailerFromString(s, &options)
	if err != nil {
		return "", err
	}

	// Transform the HTML to apply inline CSS styles.
	html, err := prem.Transform()
	if err != nil {
		return "", err
	}

	return html, nil
}

// getEncryption determines the correct encryption mode for the SMTP connection
// based on the input string: "tls", "ssl", or "none".
func (m *Mail) getEncryption(s string) mail.Encryption {
	switch s {
	case "tls":
		return mail.EncryptionSTARTTLS // Use STARTTLS encryption.
	case "ssl":
		return mail.EncryptionSSLTLS // Use SSL/TLS encryption.
	case "none", "":
		return mail.EncryptionNone // No encryption.
	default:
		// Default to STARTTLS if an unrecognized value is provided.
		return mail.EncryptionSTARTTLS
	}
}
