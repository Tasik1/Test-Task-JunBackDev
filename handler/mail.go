package handler

import (
	"net/smtp"
	"os"
)

func SendEmail(email string, subject string, HTMLbody string) error {
	to := []string{email}

	fromEmail := os.Getenv("EMAIL_LOGIN")
	SMTPpassword := os.Getenv("SMTP_PSW")
	host := "smtp.gmail.com"
	port := "587"
	address := host + ":" + port

	auth := smtp.PlainAuth("", fromEmail, SMTPpassword, host)

	msg := []byte(
		"From: <" + fromEmail + ">\r\n" +
			"To: " + email + "\r\n" +
			"Subject: " + subject + "\r\n" +
			"MIME-version: 1.0\r\n" +
			"Content-Type: text/html; charset=\"UTF-8\";\r\n" +
			"\r\n" +
			HTMLbody)

	err := smtp.SendMail(address, auth, fromEmail, to, msg)

	if err != nil {
		return err
	}
	return nil
}
