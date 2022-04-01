package controllers

import (
	"fmt"

	"github.com/claudiu/gocron"
)

func Gocron(epgi string, kuser string) {
	gocron.Start()
	// gocron.Every(10).Seconds().Do(gomail.SendMorningMail, epgi, kuser)
	gocron.Every(10).Seconds().Do(Sapa, epgi, kuser)
}

func Sapa(epgi string, kuser string) {
	fmt.Print(epgi + kuser)
}
