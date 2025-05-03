package mailer

import (
	"bytes"
	"fmt"
	"html/template"
	"log"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"time"

	"github.com/vanng822/go-premailer/premailer"
	mail "github.com/xhit/go-simple-mail/v2"
)

// Mail holds mail server configuration
type Mail struct {
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromName    string
	FromAddress string
}

// Message is the data required to send an email
type Message struct {
	From        string
	FromName    string
	To          string
	Subject     string
	Template    string
	Attachments []string
	Data        interface{}
}

// NewMail creates a Mail instance with the provided configuration
func NewMail(domain, host string, port int, username, password, encryption, fromName, fromAddress string) Mail {
	return Mail{
		Domain:      domain,
		Host:        host,
		Port:        port,
		Username:    username,
		Password:    password,
		Encryption:  encryption,
		FromName:    fromName,
		FromAddress: fromAddress,
	}
}

// SendSMTPMessage sends an email message via SMTP
func (m *Mail) SendSMTPMessage(msg Message) error {
	if msg.From == "" {
		msg.From = m.FromAddress
	}

	if msg.FromName == "" {
		msg.FromName = m.FromName
	}

	log.Printf("Preparing to send email: To=%s, Subject=%s, From=%s", msg.To, msg.Subject, msg.From)
	log.Printf("Using SMTP server: %s:%d, User: %s", m.Host, m.Port, m.Username)

	// Convert message data to map if possible
	var dataMap map[string]interface{}

	// Convert message data to map if possible
	switch v := msg.Data.(type) {
	case map[string]interface{}:
		dataMap = v
	case map[string]string:
		// Convert string map to interface{} map
		dataMap = make(map[string]interface{})
		for key, value := range v {
			dataMap[key] = value
		}
	default:
		// Default to simple message wrapper
		dataMap = map[string]interface{}{
			"message": msg.Data,
		}
	}

	// Add current year (used in footers)
	dataMap["year"] = time.Now().Year()

	// Create email template
	htmlMessage, err := m.buildHTMLMessage(msg, dataMap)
	if err != nil {
		log.Printf("Error building HTML message: %v", err)
		return err
	}

	// Create plain text version
	plainMessage, err := m.buildPlainTextMessage(msg, dataMap)
	if err != nil {
		log.Printf("Error building plain text message: %v", err)
		return err
	}

	// Set up SMTP server
	server := mail.NewSMTPClient()
	server.Host = m.Host
	server.Port = m.Port
	server.Username = m.Username
	server.Password = m.Password
	server.Encryption = m.getEncryption()
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	// Enable SMTP authentication debugging
	server.Authentication = mail.AuthLogin // Explicitly set authentication method

	log.Printf("Using encryption type: %v", server.Encryption)

	// Connect to SMTP server with retry
	var smtpClient *mail.SMTPClient
	var retryCount int
	var connectErr error

	for retryCount < 3 {
		log.Printf("Connecting to SMTP server, attempt %d", retryCount+1)
		smtpClient, connectErr = server.Connect()
		if connectErr == nil {
			break
		}
		retryCount++
		log.Printf("SMTP connection attempt %d failed: %v", retryCount, connectErr)
		time.Sleep(time.Duration(retryCount) * time.Second)
	}

	if connectErr != nil {
		log.Printf("Failed to connect to SMTP server after 3 attempts: %v", connectErr)
		return fmt.Errorf("failed to connect to SMTP server after 3 attempts: %w", connectErr)
	}

	log.Printf("Successfully connected to SMTP server")

	// Create email message
	email := mail.NewMSG()
	email.SetFrom(fmt.Sprintf("%s <%s>", msg.FromName, msg.From))
	email.AddTo(msg.To)
	email.SetSubject(msg.Subject)
	email.SetBody(mail.TextPlain, plainMessage)
	email.AddAlternative(mail.TextHTML, htmlMessage)

	// Add attachments if any
	if len(msg.Attachments) > 0 {
		for _, attachment := range msg.Attachments {
			email.AddAttachment(attachment)
		}
	}

	// Check email validity before sending
	if email.Error != nil {
		log.Printf("Error in email format: %v", email.Error)
		return email.Error
	}

	// Send email
	log.Printf("Attempting to send email to %s", msg.To)
	err = email.Send(smtpClient)
	if err != nil {
		log.Printf("Error sending email: %v", err)
		return err
	}

	log.Printf("Email successfully sent to %s", msg.To)
	return nil
}

// getEncryption converts encryption string to mail package constant
func (m *Mail) getEncryption() mail.Encryption {
	switch m.Encryption {
	case "tls":
		return mail.EncryptionSTARTTLS
	case "ssl":
		return mail.EncryptionSSLTLS
	default:
		return mail.EncryptionNone
	}
}

// buildHTMLMessage builds an HTML email message from template
func (m *Mail) buildHTMLMessage(msg Message, data map[string]interface{}) (string, error) {
	// Ensure template has proper extension
	templateName := msg.Template
	if !strings.HasSuffix(templateName, ".html.gohtml") && !strings.HasSuffix(templateName, ".gohtml") {
		templateName = templateName + ".html.gohtml"
		log.Printf("DEBUG: Added .html.gohtml extension to template name: %s", templateName)
	}

	// Check for template path with fallback
	templatePath := fmt.Sprintf("./templates/%s", templateName)
	log.Printf("DEBUG: Looking for template at: %s", templatePath)

	// Try different paths in this order:
	possiblePaths := []string{
		fmt.Sprintf("./templates/%s", templateName),
		fmt.Sprintf("../templates/%s", templateName),
		fmt.Sprintf("/app/templates/%s", templateName),
		fmt.Sprintf("./templates/%s.html.gohtml", strings.TrimSuffix(msg.Template, ".html.gohtml")),
		fmt.Sprintf("../templates/%s.html.gohtml", strings.TrimSuffix(msg.Template, ".html.gohtml")),
		fmt.Sprintf("/app/templates/%s.html.gohtml", strings.TrimSuffix(msg.Template, ".html.gohtml")),
		"./templates/mail.html.gohtml",
		"../templates/mail.html.gohtml",
		"/app/templates/mail.html.gohtml",
	}

	templateFound := false
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			templatePath = path
			templateFound = true
			log.Printf("DEBUG: Template found at: %s", templatePath)
			break
		}
	}

	if !templateFound {
		log.Printf("ERROR: Template not found after trying multiple paths")
		return "", fmt.Errorf("template not found: %s", msg.Template)
	}

	// Log the data we're providing to the template
	log.Printf("DEBUG: Template data keys: %v", getMapKeys(data))

	t, err := template.New(filepath.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		log.Printf("ERROR: Failed to parse template: %v", err)
		return "", err
	}

	var tpl bytes.Buffer
	// Try to execute with a named block first (for templates with {{define "body"}} blocks)
	err = t.ExecuteTemplate(&tpl, "body", data)
	if err != nil {
		log.Printf("DEBUG: Failed to execute template with 'body' block: %v, trying without block", err)
		// If failing, try executing directly
		err = t.Execute(&tpl, data)
		if err != nil {
			log.Printf("ERROR: Failed to execute template: %v", err)
			return "", err
		}
	}

	// Inline CSS with premailer
	prem, err := premailer.NewPremailerFromString(tpl.String(), premailer.NewOptions())
	if err != nil {
		log.Printf("ERROR: Failed to create premailer: %v", err)
		return "", err
	}

	html, err := prem.Transform()
	if err != nil {
		log.Printf("ERROR: Failed to transform HTML: %v", err)
		return "", err
	}

	return html, nil
}

// Helper function to get sorted keys from a map
func getMapKeys(m map[string]interface{}) []string {
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	return keys
}

// buildPlainTextMessage builds a plain-text email from template
func (m *Mail) buildPlainTextMessage(msg Message, data map[string]interface{}) (string, error) {
	// Try to derive plain text template name from HTML template name
	plainTemplate := ""

	// If template ends with .html.gohtml, replace with .plain.gohtml
	if strings.HasSuffix(msg.Template, ".html.gohtml") {
		plainTemplate = strings.Replace(msg.Template, ".html.gohtml", ".plain.gohtml", 1)
	} else if strings.HasSuffix(msg.Template, ".gohtml") {
		// If template ends with .gohtml, replace with .plain.gohtml
		plainTemplate = strings.Replace(msg.Template, ".gohtml", ".plain.gohtml", 1)
	} else {
		// If no extension, add .plain.gohtml
		plainTemplate = msg.Template + ".plain.gohtml"
	}

	// As a fallback, split by dot and use first part
	if plainTemplate == msg.Template || plainTemplate == "" {
		parts := strings.Split(msg.Template, ".")
		if len(parts) > 0 {
			plainTemplate = parts[0] + ".plain.gohtml"
		} else {
			plainTemplate = "mail.plain.gohtml"
		}
	}

	log.Printf("DEBUG: Looking for plain text template at: %s", plainTemplate)

	// Try different paths in this order:
	possiblePaths := []string{
		fmt.Sprintf("./templates/%s", plainTemplate),
		fmt.Sprintf("../templates/%s", plainTemplate),
		fmt.Sprintf("/app/templates/%s", plainTemplate),
		fmt.Sprintf("./templates/%s.plain.gohtml", strings.TrimSuffix(msg.Template, ".html.gohtml")),
		fmt.Sprintf("./templates/%s.plain.gohtml", strings.TrimSuffix(msg.Template, ".gohtml")),
		fmt.Sprintf("./templates/%s.plain.gohtml", msg.Template),
		"./templates/mail.plain.gohtml",
		"../templates/mail.plain.gohtml",
		"/app/templates/mail.plain.gohtml",
	}

	templatePath := ""
	templateFound := false
	for _, path := range possiblePaths {
		if _, err := os.Stat(path); !os.IsNotExist(err) {
			templatePath = path
			templateFound = true
			log.Printf("DEBUG: Plain text template found at: %s", templatePath)
			break
		}
	}

	if !templateFound {
		log.Printf("ERROR: Plain text template not found after trying multiple paths")
		return "", fmt.Errorf("plain text template not found: %s", plainTemplate)
	}

	t, err := template.New(filepath.Base(templatePath)).ParseFiles(templatePath)
	if err != nil {
		log.Printf("ERROR: Failed to parse plain text template: %v", err)
		return "", err
	}

	var tpl bytes.Buffer
	// Try to execute with a named block first (for templates with {{define "body"}} blocks)
	err = t.ExecuteTemplate(&tpl, "body", data)
	if err != nil {
		log.Printf("DEBUG: Failed to execute plain text template with 'body' block: %v, trying without block", err)
		// If failing, try executing directly
		err = t.Execute(&tpl, data)
		if err != nil {
			log.Printf("ERROR: Failed to execute plain text template: %v", err)
			return "", err
		}
	}

	plainText := tpl.String()
	return plainText, nil
}
