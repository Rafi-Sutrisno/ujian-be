package utils

import (
	"fmt"
	"mods/config"

	"gopkg.in/gomail.v2"
)

func SendMail(toEmail string, subject string, body string) error {
	emailConfig, err := config.NewEmailConfig()
	if err != nil {
		fmt.Println("Failed to load email config:", err)
		return err
	}

	fmt.Println("Email Config Loaded Successfully:")
	fmt.Println("Host:", emailConfig.Host)
	fmt.Println("Port:", emailConfig.Port)
	fmt.Println("AuthEmail:", emailConfig.AuthEmail)
	// fmt.Println("AuthPassword:", emailConfig.AuthPassword) // ❗ don't log password in real app

	mailer := gomail.NewMessage()
	mailer.SetHeader("From", emailConfig.SenderName+" <"+emailConfig.AuthEmail+">")
	mailer.SetHeader("To", toEmail)
	mailer.SetHeader("Subject", subject)
	mailer.SetBody("text/html", body)

	dialer := gomail.NewDialer(
		emailConfig.Host,
		emailConfig.Port,
		emailConfig.AuthEmail,
		emailConfig.AuthPassword,
	)

	fmt.Println("Sending email to:", toEmail)
	fmt.Println("Subject:", subject)

	err = dialer.DialAndSend(mailer)
	if err != nil {
		fmt.Println("Failed to send email:", err)
		return err
	}

	fmt.Println("Email sent successfully!")
	return nil
}
