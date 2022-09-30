package notification

type PushNotification struct {
	// dingtalk
	DingURL    string
	DingEnable bool

	// email
	EmailEnable   bool
	EmailHost     string
	EmailPort     int
	EmailUser     string
	EmailPassword string

	// message
	StageName   string
	PublishName string
	StepName    string
	Status      int64
}
