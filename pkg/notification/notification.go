package notification

type INotify interface {
	Send(m PushNotification) error
}

func NewHandlers(notify *PushNotification) []INotify {

	var handlers []INotify

	if notify.DingEnable && len(notify.DingURL) > 0 {
		handlers = append(handlers, DingRobotHandler(notify.DingURL))
	}

	if notify.EmailEnable && len(notify.EmailHost) > 0 && len(notify.EmailUser) > 0 && len(notify.EmailPassword) > 0 {
		handlers = append(handlers, EmailHandler(notify.EmailHost, notify.EmailUser, notify.EmailPassword, notify.EmailPort))
	}

	return handlers
}

func Send(options PushNotification) {

	notify := &PushNotification{
		DingURL:       options.DingURL,
		DingEnable:    options.DingEnable,
		EmailEnable:   options.EmailEnable,
		EmailHost:     options.EmailHost,
		EmailPort:     options.EmailPort,
		EmailUser:     options.EmailUser,
		EmailPassword: options.EmailPassword,
		PublishName:   options.PublishName,
		StageName:     options.StageName,
		StepName:      options.StepName,
		Status:        options.Status,
	}

	handlers := NewHandlers(notify)

	for _, handler := range handlers {
		go handler.Send(*notify)
	}

}
