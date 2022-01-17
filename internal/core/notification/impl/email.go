package notification

import (
	"bytes"
	"github.com/astaxie/beego/logs"
	"strings"

	"github.com/astaxie/beego"
	messages "github.com/go-atomci/atomci/internal/core/notification/types"
	"github.com/go-gomail/gomail"
)

type Email struct{}

func EmailHandler() INotify {
	notifyHandler := &Email{}
	return notifyHandler
}

func emailEventMessage(template INotifyTemplate, result messages.StepCallbackResult) messages.EventMessage {

	smtpHost := beego.AppConfig.String("notification::smtpHost")
	smtpAccount := beego.AppConfig.String("notification::smtpAccount")
	smtpPassword := beego.AppConfig.String("notification::smtpPassword")
	smtpPort, _ := beego.AppConfig.Int("notification::smtpPort")

	var buf bytes.Buffer
	subject := template.GenSubject(&buf, result)
	buf.Reset()
	template.GenContent(&buf, result)
	template.GenFooter(&buf, result)

	mailMessage := &messages.MailMessage{
		SmtpPort:     smtpPort,
		SmtpHost:     smtpHost,
		SmtpAccount:  smtpAccount,
		SmtpPassword: smtpPassword,
		Body:         buf.String(),
		Subject:      subject,
	}
	msg := messages.EventMessage{
		Mail: mailMessage,
	}

	return msg
}

func (email *Email) Send(result messages.StepCallbackResult) error {

	template := &EmailTemplate{}

	message := emailEventMessage(template, result)

	body := message.Mail.Body
	body = strings.Replace(body, "\n", "<br>", -1)

	subject := message.Mail.Subject

	m := gomail.NewMessage()
	m.SetHeader("From", message.Mail.SmtpAccount)
	m.SetHeader("To", message.Mail.SmtpAccount)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	d := gomail.NewDialer(message.Mail.SmtpHost, message.Mail.SmtpPort, message.Mail.SmtpAccount, message.Mail.SmtpPassword)

	defer func() {
		if r := recover(); r != nil {
			logs.Error("%v", r)
		}
	}()

	if err := d.DialAndSend(m); err != nil {
		panic(err)
	}

	return nil
}
