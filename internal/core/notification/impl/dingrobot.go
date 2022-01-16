package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
	messages "github.com/go-atomci/atomci/internal/core/notification/types"
)

type DingRobot struct{}

func DingRobotHandler() INotify {
	notifyHandler := &DingRobot{}
	return notifyHandler
}

func dingEventMessage(template INotifyTemplate, result messages.StepCallbackResult) messages.EventMessage {

	robotHost := beego.AppConfig.String("notification::ding")

	var buf bytes.Buffer
	template.GenSubject(&buf, result)
	template.GenContent(&buf, result)
	template.GenFooter(&buf, result)

	markdownText := &messages.MarkdownMessage{
		MsgType: messages.MarkDown,
		At: messages.MessageAt{
			AtMobiles: []string{},
			IsAtAll:   false,
		},
		MarkdownText: messages.MessageMarkdown{
			Title: "流水线通知",
			Text:  buf.String(),
		},
	}

	dingMsg := &messages.DingMessage{
		RobotHost:    []string{robotHost},
		EventMessage: markdownText,
	}

	msg := messages.EventMessage{
		Ding: dingMsg,
	}

	return msg
}

func (dingtalk *DingRobot) Send(result messages.StepCallbackResult) error {

	template := &DingRobotMarkdownTemplate{}

	message := dingEventMessage(template, result)

	body, err := json.Marshal(message.Ding.EventMessage)
	if err != nil {
		return fmt.Errorf("序列化消息失败 err:%s message:%+v", err, message)
	}

	for _, host := range message.Ding.RobotHost {
		res, err := http.Post(host, "application/json", bytes.NewBuffer(body))
		if err != nil {
			fmt.Println("钉钉消息发送失败 ", err)
			return fmt.Errorf("钉钉消息发送失败 err:%s", err)
		}
		if res != nil && res.Body != nil {
			defer res.Body.Close()
			content, err := ioutil.ReadAll(res.Body)
			if err != nil {
				return fmt.Errorf("读取钉钉响应失败 err:%s", err)
			}
			if res.StatusCode == http.StatusFound {
				return fmt.Errorf("机器人可能已被限流 请注意发送频率不要过高 会自行恢复")
			}
			if res.StatusCode != http.StatusOK {
				return fmt.Errorf("钉钉消息发送失败 err:%s", content)
			}
			return nil
		}
	}

	return fmt.Errorf("钉钉响应为空")
}
