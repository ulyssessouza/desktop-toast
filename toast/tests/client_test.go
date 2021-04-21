package tests

import (
	"os"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/ulyssessouza/desktop-toast/toast"
)

func TestClient(t *testing.T) {
	s := NewNotificationsServer(toast.Socket)
	s.Start()
	defer s.Stop()

	_, err := os.Stat(toast.Socket)
	for err != nil {
		time.Sleep(time.Second)
		_, err = os.Stat(toast.Socket)
	}

	c := toast.NewClient()
	assert.Nil(t, c.Send(toast.NativeNotificationModel{
		Title: "mytest_title",
		Body:  "mytest_body",
		Level: toast.Warning,
	}))
	usage := s.GetUsage()
	assert.Equal(t, []string{`{"title":"mytest_title","body":"mytest_body","level":1}`}, usage)
}
