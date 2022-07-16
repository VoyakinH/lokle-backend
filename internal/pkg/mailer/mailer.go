package mailer

import (
	"fmt"

	"github.com/VoyakinH/lokle_backend/config"
	"gopkg.in/gomail.v2"
)

func SendVerifiedEmail(to_email string, first_name string, second_name string, token string) error {
	msg := gomail.NewMessage()
	msg.SetHeader("From", config.Mailer.Email)
	msg.SetHeader("To", to_email)
	msg.SetHeader("Subject", "Подтвержение почты Столичный-КИТ")
	msg.SetBody("text/html", fmt.Sprintf("Приветствуем, %s %s! <br/> Для подтверждения электронной почты, пройдите, пожалуйста, по ссылке: <br/>  http://185.225.34.197/login?verification_email_token=%s <br/> Если Вы получили это письмо по ошибке, просто игнорируйте его. <br/> Ссылка активна в течении 7 дней.", first_name, second_name, token))

	n := gomail.NewDialer("smtp.mail.ru", 465, config.Mailer.Email, config.Mailer.Password)

	if err := n.DialAndSend(msg); err != nil {
		return err
	}
	return nil
}
