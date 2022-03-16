package notification

var (
	notificationInstance Notification
	DManager             dispatcher_manager
)

func init() {
	notificationInstance = NewDispatcher()
}

type dispatcher_manager struct {
}

func (s *dispatcher_manager) SendNotification() error {
	return nil
}
