package main

import (
	"log"
	"math/rand"
	"net/smtp"
)

func generate18DigitID() int64 {
	return (rand.Int63n(7)+1)*1e18 + rand.Int63n(1e17)
}

/*
	Hide password field at marshal
*/

func (passwordType) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

func sendVerifcationMail(to, body string) bool {
	from := "league.noreply@gmail.com"
	pass := "indreias@leagueINC"
	msg := "From: " + from + "\n" +
		"To: " + to + "\n" +
		"Subject: League INC verification\n\n" + body
	err := smtp.SendMail("smtp.gmail.com:587",
		smtp.PlainAuth("", from, pass, "smtp.gmail.com"),
		from, []string{to}, []byte(msg))
	if err != nil {
		log.Printf("smtp error: %s", err)
		return false
	}
	return true
}
