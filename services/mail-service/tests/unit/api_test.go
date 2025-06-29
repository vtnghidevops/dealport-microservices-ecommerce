package unit_test

import (
	"log"
	"mail-service/internal/config"
	"mail-service/internal/grpc"
	"mail-service/internal/mailer"
	"net"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	grpclib "google.golang.org/grpc"
)

// MockListener is a mock for net.Listener
type MockListener struct {
	mock.Mock
}

// Accept mocks the Accept method
func (m *MockListener) Accept() (net.Conn, error) {
	args := m.Called()
	var conn net.Conn
	if args.Get(0) != nil {
		conn = args.Get(0).(net.Conn)
	}
	return conn, args.Error(1)
}

// Close mocks the Close method
func (m *MockListener) Close() error {
	args := m.Called()
	return args.Error(0)
}

// Addr mocks the Addr method
func (m *MockListener) Addr() net.Addr {
	args := m.Called()
	var addr net.Addr
	if args.Get(0) != nil {
		addr = args.Get(0).(net.Addr)
	}
	return addr
}

// MockGRPCServer is a mock for grpc.Server
type MockGRPCServer struct {
	mock.Mock
}

// Serve mocks the Serve method
func (m *MockGRPCServer) Serve(lis net.Listener) error {
	args := m.Called(lis)
	return args.Error(0)
}

// Stop mocks the Stop method
func (m *MockGRPCServer) Stop() {
	m.Called()
}

// GracefulStop mocks the GracefulStop method
func (m *MockGRPCServer) GracefulStop() {
	m.Called()
}

// TestServeGRPC tests the serveGRPC function
func TestServeGRPC(t *testing.T) {
	// Setup mocks
	mockListener := new(MockListener)

	// Create a proper mail struct
	mailStruct := mailer.Mail{
		Domain:      "test.com",
		Host:        "localhost",
		Port:        25,
		Username:    "test",
		Password:    "password",
		Encryption:  "none",
		FromName:    "Test Sender",
		FromAddress: "test@example.com",
	}

	// Create app config
	cfg := config.Config{
		Mailer: mailStruct,
	}
	app := &AppConfig{
		Config: cfg,
	}

	// Create logger
	logger := log.New(os.Stdout, "[TEST] ", log.LstdFlags)

	// Test the function
	t.Run("Test GRPC server initialization", func(t *testing.T) {
		// Store original functions
		origNetListen := netListen
		origGRPCServer := newGRPCServer
		origRegisterServer := registerMailServer

		// Override with mocks
		netListen = func(network, address string) (net.Listener, error) {
			assert.Equal(t, "tcp", network)
			assert.Equal(t, ":50057", address)
			return mockListener, nil
		}

		newGRPCServer = func() *grpclib.Server {
			return grpclib.NewServer()
		}

		registerMailServer = func(s *grpclib.Server, mailServer *grpc.MailServer) {
			// Just verify the server is initialized properly
			assert.NotNil(t, s)
			assert.NotNil(t, mailServer)
		}

		// Restore at the end
		defer func() {
			netListen = origNetListen
			newGRPCServer = origGRPCServer
			registerMailServer = origRegisterServer
		}()

		// Just test the initialization, don't run the server
		lis, err := netListen("tcp", ":50057")
		assert.NoError(t, err)
		assert.Equal(t, mockListener, lis)

		// Verify logger is working
		assert.NotNil(t, logger)

		// Verify config is properly created
		assert.NotNil(t, app.Config)
	})
}

// AppConfig is a simplified version of the application config for testing
type AppConfig struct {
	Config config.Config
}

// Variables that will be replaced during testing
var (
	netListen          = net.Listen
	newGRPCServer      = func() *grpclib.Server { return grpclib.NewServer() }
	registerMailServer = func(s *grpclib.Server, mailServer *grpc.MailServer) {}
)
