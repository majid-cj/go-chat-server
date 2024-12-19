package util

import (
	"bytes"
	"errors"
	"fmt"
	"html/template"
	"net/smtp"
	"os"
	"path/filepath"
)

var emailAuth smtp.Auth

// ReceiveMail ...
type ReceiverMail struct {
	ReceiverName string
	ReceiverMail string
	ReceiverCode   string
}

// ActiveLoginMail ...
type ActiveLoginMail struct {
	ReceiverName string
	ReceiverMail string
	LoginTime   string
	DeviceInfo string
	IPAddress string
}

// SendMail ...
func SendMail(to []string, data interface{}, template string, subject string) error {
	emailHost := os.Getenv("EMAIL_HOST")
	emailFrom := os.Getenv("EMAIL_FROM")
	emailPassword := os.Getenv("EMAIL_PASSWORD")
	emailPort := os.Getenv("EMAIL_PORT")

	emailAuth = smtp.PlainAuth("", emailFrom, emailPassword, emailHost)

	emailBody, err := ParseTemplate(template, data)
	if err != nil {
		return errors.New("unable to parse email template")
	}

	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	mail_subject := "Subject: " + subject + "!\n"
	msg := []byte(mail_subject + mime + "\n" + emailBody)
	addr := fmt.Sprintf("%s:%s", emailHost, emailPort)

	if err := smtp.SendMail(addr, emailAuth, emailFrom, to, msg); err != nil {
		return  err
	}
	return nil
}

// ParseTemplate ...
func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	templatePath, err := filepath.Abs(fmt.Sprintf("util/template/%s", templateFileName))
	if err != nil {
		return "", errors.New("invalid template name")
	}
	t, err := template.ParseFiles(templatePath)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return "", err
	}
	body := buf.String()
	return body, nil
}
