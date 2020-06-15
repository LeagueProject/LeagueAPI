package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"
)

/*
	Verifica daca exista user cu userName in db
*/
func findUserByUsername(userName string) bool {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE Username='%s'", userName))
	if err != nil {
		return false
	}
	return userQuery.Next()
}

/*
	Verifica daca exista uid cu userID in db
*/
func findUserByID(userID int64) bool {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE UID=%v", userID))
	if err != nil {
		return false
	}
	return userQuery.Next()
}

/*
	Request la DB dupa uID
	Returneaza un user gol la eroare sau user si nil
*/
func getUserByID(uID int64) (User, error) {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE UID=%v", uID))
	if err != nil {
		return User{}, err
	}
	if userQuery.Next() {
		var us User
		userQuery.Scan(&us.UID, &us.InstitutionEmail, &us.PersonalEmail, &us.Username, &us.PasswordHash, &us.YearOfStudy, &us.College, &us.University, &us.Major, &us.Serie, &us.FirstName, &us.LastName, &us.verified)
		return us, nil
	}
	return User{}, errors.New("User does not exist")
}

/*
	Request la DB dupa user
	Returneaza un user gol la eroare sau user si nil
*/
func getUserByUsername(username string) (User, error) {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE Username='%s'", username))
	if err != nil {
		return User{}, err
	}
	if userQuery.Next() {
		var us User
		userQuery.Scan(&us.UID, &us.InstitutionEmail, &us.PersonalEmail, &us.Username, &us.PasswordHash, &us.YearOfStudy, &us.College, &us.University, &us.Major, &us.Serie, &us.FirstName, &us.LastName, &us.verified)
		return us, nil
	}
	return User{}, errors.New("User does not exist")
}

func canLogin(user, pass string) bool {
	userQuery, err := db.Query(fmt.Sprintf("SELECT Username,Password FROM league WHERE Username='%s'", user))
	if err != nil {
		return false
	}
	if userQuery.Next() == false {
		return false
	}
	var fsUser, fsPassword string
	userQuery.Scan(&fsUser, &fsPassword)
	fmt.Println(fsPassword)
	fmt.Println(fmt.Sprintf("%x", md5.Sum([]byte(pass))))
	return fsUser == user && fsPassword == fmt.Sprintf("%x", md5.Sum([]byte(pass)))
}

func addUser(newUser User) HTTPResponse {
	if findUserByUsername(newUser.Username) {
		return HTTPResponse{Response: []string{"User already exists"}, Code: 409}
	}
	newUser.PasswordHash = passwordType(fmt.Sprintf("%x", md5.Sum([]byte(string(newUser.PasswordHash)))))
	userID := generate16DigitID()
	for findUserByID(userID) {
		userID = generate16DigitID()
	}
	newUser.UID = userID
	if sendVerifcationMail(newUser.InstitutionEmail, "http://35.184.233.76:8080/activate?id="+strconv.FormatInt(userID, 10)) == false {
		return HTTPResponse{Response: []string{"Invalid Email"}, Code: 404}
	}
	sqlStatement := `INSERT INTO league (UID,IEmail,PMail,Username,Password,YearOfStudy,College,University,Major,Serie,FirstName,LastName,verified) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13)`
	_, err := db.Exec(sqlStatement,
		newUser.UID, newUser.InstitutionEmail, newUser.PersonalEmail, newUser.Username,
		newUser.PasswordHash, newUser.YearOfStudy, newUser.College, newUser.University,
		newUser.Major, newUser.Serie, newUser.FirstName, newUser.LastName, 0)
	if err != nil {
		panic(err)
	}

	fmt.Println(newUser)
	return HTTPResponse{Response: []string{"User added"}, Code: 200}
}

func checkUserByID(uID, sID int64) error {
	return nil
}

func seesionExist(sID int64) bool {
	sessionQuery, err := db.Query(fmt.Sprintf("SELECT * FROM sessions WHERE sid=%v", sID))
	if err != nil {
		return false
	}
	return sessionQuery.Next()
}

func getUIDFromSession(sID int64) int64 {
	sessionQuery, err := db.Query(fmt.Sprintf("SELECT UID FROM sessions WHERE sid=%v", sID))
	if err != nil {
		return 0
	}
	if sessionQuery.Next() == false {
		return -1
	}
	var uID int64
	sessionQuery.Scan(&uID)
	return uID
}

func sendMessage(newMessage Message) {
	sqlStatement := `INSERT INTO messages (id,authorid,text,mediafilepath,date,receiver,typeofreceiver) VALUES($1,$2,$3,$4,$5,$6,$7)`
	_, err := db.Exec(sqlStatement,
		newMessage.ID, newMessage.AuthorID, newMessage.Text, newMessage.MediaFilePath, newMessage.Date, newMessage.Receiver, newMessage.TypeOfReceiver)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(newMessage)
}

func verifyUser(uID int64) {
	db.Exec(`UPDATE league SET verified=1 WHERE uid=$1`, uID)
}

func newSessionID() int64 {
	id := generate16DigitID()
	for seesionExist(id) {
		id = generate16DigitID()
	}
	return id
}

func addSession(sID, uID int64) {
	db.Exec(`INSERT INTO sessions (sid,uid) VALUES ($1,$2)`, sID, uID)
}
