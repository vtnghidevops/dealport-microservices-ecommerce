package mocks

import (
	"mail-service/internal/mailer"

	"github.com/stretchr/testify/mock"
)

// MockMailer is a mock implementation of the mailer.Mail struct
type MockMailer struct {
	mock.Mock
	Domain      string
	Host        string
	Port        int
	Username    string
	Password    string
	Encryption  string
	FromName    string
	FromAddress string
}

// SendSMTPMessage mocks the SendSMTPMessage method
func (m *MockMailer) SendSMTPMessage(msg mailer.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}

// TestSend mocks the TestSend method for testing connection
func (m *MockMailer) TestSend(msg mailer.Message) error {
	args := m.Called(msg)
	return args.Error(0)
}
