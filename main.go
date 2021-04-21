package main

import "github.com/ulyssessouza/desktop-toast/toast"

func main() {
	c := toast.NewClient()
	err := c.Send(toast.NativeNotificationModel{
		Title: "MyToastTitle",
		Body:  "My Toast Body",
		Level: toast.Warning,
	})
	if err != nil {
		panic(err)
	}
}
