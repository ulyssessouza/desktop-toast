package toast

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net"
	"net/http"
	"time"
)

const (
	Info = iota
	Warning
	Error
)
type NotificationLevel int

type client struct {
	httpClient *http.Client
}

// NativeNotificationModel is a command
type NativeNotificationModel struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Level NotificationLevel `json:"level"`
}

// Client sends notifications to Docker Desktop's notification socket
type Client interface {
	// Send sends the notification to Docker Desktop.
	Send(NativeNotificationModel) error
}

// NewClient returns a new notification client
func NewClient() Client {
	return &client{
		httpClient: &http.Client{
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return conn()
				},
			},
		},
	}
}

// Send notification to the Docker Desktop server
func (c *client) Send(n NativeNotificationModel) error {
	result := make(chan bool, 1)
	errchan := make(chan error, 1)

	go func() {
		err := c.postNotification(n)
		if err != nil {
			errchan <- err
			return
		}
		result <- true
	}()

	// Consider timeout and/or errors
	select {
	case <-result:
	case err := <-errchan:
		return err
	case <-time.After(50 * time.Millisecond):
		return errors.New("timeout on posting notification")
	}

	return nil
}

func (c *client) postNotification(n NativeNotificationModel) error {
	req, err := json.Marshal(n)
	if err != nil {
		return err
	}
	_, err = c.httpClient.Post("http://localhost/notifications",
		"application/json", bytes.NewBuffer(req))
	if err != nil {
		return err
	}
	return nil
}
