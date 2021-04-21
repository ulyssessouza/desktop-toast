package tests

import (
	"io/ioutil"
	"log"
	"net"
	"net/http"
	"strings"

	"github.com/labstack/echo"
)

// MockToastServer a mock registering all notifications POST invocations
type MockToastServer struct {
	socket string
	usage  []string
	e      *echo.Echo
}

// NewNotificationsServer instantiate a new MockToastServer
func NewNotificationsServer(socket string) *MockToastServer {
	return &MockToastServer{
		socket: socket,
		e:      echo.New(),
	}
}

func (s *MockToastServer) handlePostUsage(c echo.Context) error {
	body, err := ioutil.ReadAll(c.Request().Body)
	if err != nil {
		return err
	}
	cliUsage := string(body)
	s.usage = append(s.usage, cliUsage)
	return c.String(http.StatusOK, "")
}

// GetUsage get usage
func (s *MockToastServer) GetUsage() []string {
	return s.usage
}

// ResetUsage reset usage
func (s *MockToastServer) ResetUsage() {
	s.usage = []string{}
}

// Stop stop the mock server
func (s *MockToastServer) Stop() {
	_ = s.e.Close()
}

// Start start the mock server
func (s *MockToastServer) Start() {
	go func() {
		listener, err := net.Listen("unix", strings.TrimPrefix(s.socket, "unix://"))
		if err != nil {
			log.Fatal(err)
		}
		s.e.Listener = listener
		s.e.POST("/notifications", s.handlePostUsage)
		_ = s.e.Start(":1323")
	}()
}
