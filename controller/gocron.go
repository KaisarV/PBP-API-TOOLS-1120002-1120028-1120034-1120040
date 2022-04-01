package controllers

import (
	gomail "GolangTools/gomail"

	"github.com/claudiu/gocron"
)

func Gocron(epgi string, kuser string) {
	gocron.Start()
	gocron.Every(10).Seconds().Do(gomail.SendMorningMail, epgi, kuser)
}
