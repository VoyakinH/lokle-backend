package mailer

import (
	"github.com/VoyakinH/lokle_backend/config"
	"gopkg.in/gomail.v2"
)

func SendVerifiedEmail(to_email string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Mailer.Email)
	msg.SetHeader("To", to_email)
	msg.SetHeader("Subject", "prikol")
	msg.SetBody("text/html", "<b>This is the body of the mail</b>")

	n := gomail.NewDialer("smtp.gmail.com", 587, config.Mailer.Email, config.Mailer.Password)

	// Send the email
	if err := n.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
