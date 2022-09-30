package notification

import (
	"bytes"

	messages "github.com/go-atomci/atomci/pkg/notification/types"
)

type INotifyTemplate interface {
	GenSubject(buf *bytes.Buffer, m PushNotification) string
	GenContent(buf *bytes.Buffer, m PushNotification) string
	GenFooter(buf *bytes.Buffer, m PushNotification) string
}

type DingRobotMarkdownTemplate struct{}

func (temp *DingRobotMarkdownTemplate) GenSubject(buf *bytes.Buffer, m PushNotification) string {

	buf.WriteString("## ")
	buf.WriteString(messages.StatusCodeToChinese(m.Status))

	return buf.String()
}

func (temp *DingRobotMarkdownTemplate) GenContent(buf *bytes.Buffer, m PushNotification) string {

	buf.WriteString(m.PublishName)
	buf.WriteString("\r\n\r\n")
	buf.WriteString(m.StageName)
	buf.WriteString("\r\n\r\n")
	buf.WriteString(m.StepName)

	return buf.String()
}

func (temp *DingRobotMarkdownTemplate) GenFooter(buf *bytes.Buffer, m PushNotification) string {

	buf.WriteString("\r\n\r\n> by AtomCI")

	return buf.String()
}

type EmailTemplate struct{}

func (temp *EmailTemplate) GenSubject(buf *bytes.Buffer, m PushNotification) string {

	buf.WriteString("流水线")
	buf.WriteString(m.PublishName)
	buf.WriteString("构建")
	buf.WriteString(messages.StatusCodeToChinese(m.Status))

	return buf.String()
}

func (temp *EmailTemplate) GenContent(buf *bytes.Buffer, m PushNotification) string {

	buf.WriteString("<p><h2><span>流水线: </span><b>")
	buf.WriteString(m.PublishName)
	buf.WriteString("</b></h2></p><p><h2><span>阶段: </span><b>")
	buf.WriteString(m.StageName)
	buf.WriteString("</b></h2></p><p><h2><span>步骤: </span><b>")
	buf.WriteString(m.StepName)
	buf.WriteString("</b></h2></p><p><h1>")
	buf.WriteString(messages.StatusCodeToChinese(m.Status))
	buf.WriteString("</h1>")

	return buf.String()
}

func (temp *EmailTemplate) GenFooter(buf *bytes.Buffer, m PushNotification) string {

	buf.WriteString("<p>by AtomCI</p>")

	return buf.String()
}
