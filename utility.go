package main

import (
	"log"
	"math/rand"
	"net/smtp"
)

/**
* @desc Returneaza un ID unic
* @param none
* @return int64 , un numar random de 18 cifre , sub 9e18 , deoarece ar da overflow pe negativ
* @author Mihai Indreias
 */

func generate18DigitID() int64 {
	return (rand.Int63n(7)+1)*1e18 + rand.Int63n(1e17)
}

/**
* @desc HelperMethod care face field-ul de parola nul la transformarea unei variable User in JSON
* @param passwordType ( string ) -> nu conteaza param , deoarece oricum returneaza ""
* @return []byte,error : [] si nil (nu pot exista erori)
* @author Mihai Indreias
 */

func (passwordType) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

/**
* @desc Trimite un email(principal pentru verificare)
* @param $to (string)    :adresa email a utlizatorului
		 $body (string ) :continutul emailului
* @return bool -> returneaza daca a fost trimis cu succes email ul
* @author Mihai Indreias
*/

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
