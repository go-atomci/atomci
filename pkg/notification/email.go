package notification

import (
	"bytes"
	"strings"

	messages "github.com/go-atomci/atomci/pkg/notification/types"
	"github.com/go-gomail/gomail"
)

type Email struct {
	smtpHost     string
	stmpAccount  string
	smtpPassword string
	smtpPort     int
}

func EmailHandler(host, user, password string, port int) INotify {
	notifyHandler := &Email{
		smtpHost:     host,
		stmpAccount:  user,
		smtpPassword: password,
		smtpPort:     port,
	}
	return notifyHandler
}

func emailEventMessage(template INotifyTemplate, result PushNotification) messages.EventMessage {
	var buf bytes.Buffer
	subject := template.GenSubject(&buf, result)
	buf.Reset()
	template.GenContent(&buf, result)
	template.GenFooter(&buf, result)

	mailMessage := &messages.MailMessage{
		SmtpPort:     result.EmailPort,
		SmtpHost:     result.EmailHost,
		SmtpAccount:  result.EmailUser,
		SmtpPassword: result.EmailPassword,
		Body:         buf.String(),
		Subject:      subject,
	}
	msg := messages.EventMessage{
		Mail: mailMessage,
	}

	return msg
}

func (email *Email) Send(result PushNotification) error {

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

	return d.DialAndSend(m)
}
