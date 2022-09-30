package notification

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	messages "github.com/go-atomci/atomci/pkg/notification/types"
)

type DingRobot struct {
	DingUrl string
}

func DingRobotHandler(dingUrl string) INotify {
	notifyHandler := &DingRobot{
		DingUrl: dingUrl,
	}
	return notifyHandler
}

func (dingtalk *DingRobot) dingEventMessage(template INotifyTemplate, result PushNotification) messages.EventMessage {

	robotHost := dingtalk.DingUrl

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

func (dingtalk *DingRobot) Send(result PushNotification) error {

	template := &DingRobotMarkdownTemplate{}

	message := dingtalk.dingEventMessage(template, result)

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
