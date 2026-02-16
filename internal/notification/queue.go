//Queues
package notifications

type Job struct {
	UserID  string
	Message string
	Type    string
}

var NotificationQueue = make(chan Job, 100)
