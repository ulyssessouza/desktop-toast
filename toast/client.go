package toast

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net"
	"net/http"
	"time"
)

type client struct {
	httpClient *http.Client
}

// NewClient returns a new notification client
func NewClient(timeout time.Duration) Client {
	return &client{
		httpClient: &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				DialContext: func(_ context.Context, _, _ string) (net.Conn, error) {
					return conn()
				},
			},
		},
	}
}

type NotificationLevel int

const (
	Info NotificationLevel = iota
	Warning
	Error
)

func (n NotificationLevel) validate() error {
	if n < Info || n > Error {
		return fmt.Errorf("toast 'level' '%d' is invalid", n)
	}
	return nil
}

// NativeNotificationModel is a command
type NativeNotificationModel struct {
	Title string            `json:"title"`
	Body  string            `json:"body"`
	Level NotificationLevel `json:"level"`
}

func (n NativeNotificationModel) validate() error {
	if n.Title == "" {
		return errors.New("toast 'title' cannot be empty")
	}
	if n.Body == "" {
		return errors.New("toast 'body' cannot be empty")
	}
	return n.Level.validate()
}

// Client sends notifications to Docker Desktop's notification socket
type Client interface {
	// SendToast sends the notification to Docker Desktop.
	SendToast(NativeNotificationModel) error
}

// SendToast notification to the Docker Desktop server
func (c *client) SendToast(n NativeNotificationModel) error {
	err := n.validate()
	if err != nil {
		return err
	}

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
	case <-time.After(time.Minute): // Maximal timeout
		return errors.New("reached maximal timeout of one minute on posting notification")
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
