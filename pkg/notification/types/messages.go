package types

import (
	"fmt"

	"github.com/go-atomci/atomci/internal/models"
)

const (
	Text     MessageType = "text"
	MarkDown MessageType = "markdown"
)

type MessageAt struct {
	AtMobiles []string `json:"atmobiles"`
	IsAtAll   bool     `json:"isAtAll"`
}

type MessageType string

type MessageContent struct {
	Content string `json:"content"`
}

type MessageMarkdown struct {
	Title string `json:"title"`
	Text  string `json:"text"`
}

// PlainMessage 消息格式，钉钉专用
type PlainMessage struct {
	MsgType MessageType    `json:"msgtype"`
	At      MessageAt      `json:"at"`
	Text    MessageContent `json:"text"`
}

// MarkdownMessage 消息格式，钉钉专用
type MarkdownMessage struct {
	MsgType      MessageType     `json:"msgtype"`
	At           MessageAt       `json:"at"`
	MarkdownText MessageMarkdown `json:"markdown"`
}

// MailMessage 邮件消息
type MailMessage struct {
	SmtpPort     int
	SmtpHost     string
	SmtpAccount  string
	SmtpPassword string
	Body         string
	Subject      string
}

type DingMessage struct {
	RobotHost    []string
	EventMessage interface{}
}

type EventMessage struct {
	Mail *MailMessage
	Ding *DingMessage
}

func StatusCodeToChinese(status int64) string {
	switch status {
	case models.Success:
		return "成功"
	case models.Failed:
		return "失败"
	default:
		return fmt.Sprintf("未知状态：%d", status)
	}
}
