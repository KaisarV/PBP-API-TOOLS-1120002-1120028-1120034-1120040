package gomail

import (
	"bytes"
	"fmt"
	"html/template"
	"log"

	"gopkg.in/gomail.v2"
)

const CONFIG_SMTP_HOST = "smtp.gmail.com"
const CONFIG_SMTP_PORT = 587
const CONFIG_SENDER_NAME = "CentStore<nomen.test123@gmail.com>"
const CONFIG_AUTH_EMAIL = "nomen.test123@gmail.com"
const CONFIG_AUTH_PASSWORD = "tes12345"

type BodylinkEmail struct {
	Name string
}

func SendMail(email string, name string) {
	templateData := BodylinkEmail{
		Name: name,
	}

	result, _ := ParseTemplate("gomail/email_template_verifikasi.html", templateData)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", email, email)
	mailer.SetAddressHeader("Cc", email, "Pemberitahuan Pendaftaran Akun")
	mailer.SetHeader("Subject", "Pemberitahuan Pendaftaran Akun")
	mailer.SetBody("text/html", result)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

}

func ParseTemplate(templateFileName string, data interface{}) (string, error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return "", err
	}
	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		fmt.Println(err)
		return "", err
	}

	return buf.String(), nil
}

func SendMorningMail(epgi string, kuser string) {
	templateData := BodylinkEmail{
		Name: kuser,
	}
	result, _ := ParseTemplate("gomail/email_template_hai.html", templateData)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", kuser /*"xenon.avius@gmail.com"*/)
	mailer.SetAddressHeader("Cc" /*emailUser*/, "if-20034@students.ithb.ac.id", "Pemberitahuan Penting dari IF-20")
	mailer.SetHeader("Subject", epgi)

	mailer.SetBody("text/html", result)
}

func SendPromoMail(email string, name string) {
	templateData := BodylinkEmail{
		Name: name,
	}

	result, _ := ParseTemplate("gomail/email_template_promo.html", templateData)
	mailer := gomail.NewMessage()
	mailer.SetHeader("From", CONFIG_SENDER_NAME)
	mailer.SetHeader("To", email, email)
	mailer.SetAddressHeader("Cc", email, "Pemberitahuan Promo")
	mailer.SetHeader("Subject", "PROMO TERBATAS")
	mailer.SetBody("text/html", result)

	dialer := gomail.NewDialer(
		CONFIG_SMTP_HOST,
		CONFIG_SMTP_PORT,
		CONFIG_AUTH_EMAIL,
		CONFIG_AUTH_PASSWORD,
	)

	err := dialer.DialAndSend(mailer)
	if err != nil {
		log.Fatal(err.Error())
	}

}
