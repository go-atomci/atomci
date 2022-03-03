package notification

import (
	notification "github.com/go-atomci/atomci/internal/core/notification/impl"
	messages "github.com/go-atomci/atomci/internal/core/notification/types"
	"testing"
)

func TestNotifyEmail(t *testing.T) {
	notification.EmailHandler().Send(messages.StepCallbackResult{
		StageName:   "",
		PublishName: "",
		StepName:    "",
		Status:      1,
	})
}

func TestNotifyDingRobot(t *testing.T) {
	notification.DingRobotHandler().Send(messages.StepCallbackResult{
		StageName:   "aa",
		PublishName: "bb",
		StepName:    "cc",
		Status:      0,
	})
}
