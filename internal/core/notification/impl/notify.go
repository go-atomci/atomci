package notification

import (
	"github.com/astaxie/beego"
	messages "github.com/go-atomci/atomci/internal/core/notification/types"
	"github.com/go-atomci/atomci/internal/core/publish"
)

type INotify interface {
	Send(m messages.StepCallbackResult) error
}

func NewHandlers() []INotify {

	dingEnable, _ := beego.AppConfig.Int("notification::dingEnable")
	mailEnable, _ := beego.AppConfig.Int("notification::mailEnable")

	var handlers []INotify

	if dingEnable > 0 {
		_ = append(handlers, DingRobotHandler())
	}

	if mailEnable > 0 {
		_ = append(handlers, EmailHandler())
	}

	return handlers
}

func Send(publishId int64, status int64) {

	pm := publish.NewPublishManager()

	pub, _ := pm.GetPublishInfo(publishId)

	handlers := NewHandlers()

	callbackResult := messages.StepCallbackResult{
		PublishName: pub.Name,
		StageName:   pub.StageName,
		StepName:    pub.Step,
		Status:      status,
	}

	if handlers != nil {
		for _, handler := range handlers {
			go handler.Send(callbackResult)
		}
	}
}
