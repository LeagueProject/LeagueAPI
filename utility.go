package main

import (
	"fmt"
	"log"
	"math/rand"
	"net/smtp"
	"strconv"
)

var (
	lastUserID int64
)

/*
	Request la DB dupa uID
	Formeaza un user pregenerat
*/
func getUserByID(uID int64) (User, error) {
	return User{
		UID:              uID,
		InstitutionEmail: "mihai.indreias@unibuc.ro",
		YearOfStudy:      1,
		University:       "Politehnica",
		College:          "Politehnica Bucuresti",
		Serie:            "123",
		Major:            "Inginerie",
		PersonalEmail:    "mihai.indreias@gmail.com",
		Username:         "g0g05arui",
	}, nil
}

func getUserByUsername(username string) (User, error) {
	return User{
		UID:              12345,
		InstitutionEmail: "vlad.cainamisir@gmail.com",
		YearOfStudy:      1,
		University:       "Politehnica",
		College:          "Politehnica Bucuresti",
		Serie:            "123",
		Major:            "Inginerie",
		PersonalEmail:    "mihai.indreias@gmail.com",
		Username:         username,
	}, nil
}

func generate16DigitID() int64 {
	return 1e16 + rand.Int63n(1e15)
}

func findUserByUsername(userName string) bool {
	return false
}

func findUserByID(userID int64) bool {
	return false
}

/*
	Hide password field at marshal
*/

func (passwordType) MarshalJSON() ([]byte, error) {
	return []byte(`""`), nil
}

func addUser(newUser User) HTTPResponse {
	if findUserByUsername(newUser.Username) {
		return HTTPResponse{Response: "User already exists", Code: 409}
	}
	/*
		Add User to DB
	*/

	userID := generate16DigitID()
	for findUserByID(userID) {
		userID = generate16DigitID()
	}
	newUser.UID = userID

	sqlStatement := `INSERT INTO league (UID,IEmail,PMail,Username,Password,YearOfStudy,College,University,Major,Serie,verified) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11)`
	_, err := db.Exec(sqlStatement,
		newUser.UID, newUser.InstitutionEmail, newUser.PersonalEmail, newUser.Username,
		newUser.PasswordHash, newUser.YearOfStudy, newUser.College, newUser.University, newUser.Major, newUser.Serie, 0)
	if err != nil {
		panic(err)
	}

	sendVerifcationMail(newUser.InstitutionEmail, "http://34.67.7.77:8555/activate?id="+strconv.FormatInt(userID, 10))
	fmt.Println(newUser)
	return HTTPResponse{Response: "User added", Code: 200}
}

func seesionExist(sID int64) bool {
	return false
}

func canLogin(user, pass string) bool {
	return user == "cnmsr" && pass == "e10adc3949ba59abbe56e057f20f883e"
}

func sendVerifcationMail(to, body string) {
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
	}
}
