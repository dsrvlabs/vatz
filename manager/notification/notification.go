package notification

import "fmt"

type notification struct {
	//duration time.Duration
	//notificationChannel string
	//message string
}

type Notification interface {
	SendNotification() error
}

func (d notification) SendNotification() error {
	fmt.Println("this is sending Notification")
	return nil
}

func NewDispatcher() Notification {
	return &notification{}
}
