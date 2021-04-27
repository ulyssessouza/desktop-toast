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

	// Wait for socket to be created
	_, err := os.Stat(toast.Socket)
	for err != nil {
		time.Sleep(time.Second)
		_, err = os.Stat(toast.Socket)
	}

	c := toast.NewClient(50 * time.Millisecond)

	t.Run("Successful cases", func(t *testing.T) {
		assert.Nil(t, c.SendToast(toast.NativeNotificationModel{
			Title: "mytest_title1",
			Body:  "mytest_body1",
			Level: toast.Info,
		}))
		assert.Nil(t, c.SendToast(toast.NativeNotificationModel{
			Title: "mytest_title2",
			Body:  "mytest_body2",
			Level: toast.Warning,
		}))
		assert.Nil(t, c.SendToast(toast.NativeNotificationModel{
			Title: "mytest_title3",
			Body:  "mytest_body3",
			Level: toast.Error,
		}))
		assert.Equal(t, []string{
			`{"title":"mytest_title1","body":"mytest_body1","level":0}`,
			`{"title":"mytest_title2","body":"mytest_body2","level":1}`,
			`{"title":"mytest_title3","body":"mytest_body3","level":2}`,
		}, s.GetUsage())

	})

	t.Run("Fail cases", func(t *testing.T) {
		s.ResetUsage()
		assert.NotNil(t, c.SendToast(toast.NativeNotificationModel{
			Title: "", // invalid Title
			Body:  "some_body",
			Level: toast.Info,
		}))
		assert.NotNil(t, c.SendToast(toast.NativeNotificationModel{
			Title: "some_title",
			Body:  "", // invalid Body
			Level: toast.Info,
		}))
		assert.NotNil(t, c.SendToast(toast.NativeNotificationModel{
			Title: "mytest_title4",
			Body:  "mytest_body4",
			Level: 99, // invalid Level
		}))
		assert.Equal(t, []string{}, s.GetUsage())
	})
}
