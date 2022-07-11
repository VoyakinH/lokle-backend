package mailer

import (
	"fmt"

	"github.com/VoyakinH/lokle_backend/config"
	"gopkg.in/gomail.v2"
)

func SendVerifiedEmail(to_email string, token string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Mailer.Email)
	msg.SetHeader("To", to_email)
	msg.SetHeader("Subject", "Подтверждение почты на lokle")
	msg.SetBody("text/html", fmt.Sprintf("<b>Перейдите по ссылке http://185.225.34.197/login?verification_email_token=%s и подтвердите почту. Ссылка активна в течении 7 дней.</b>", token))

	n := gomail.NewDialer("smtp.mail.ru", 465, config.Mailer.Email, config.Mailer.Password)

	if err := n.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
