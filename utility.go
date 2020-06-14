package main

import (
	"log"
	"math/rand"
	"net/smtp"
)

func generate16DigitID() int64 {
	return 1e16 + rand.Int63n(1e15)
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
