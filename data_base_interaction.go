package main

import (
	"crypto/md5"
	"errors"
	"fmt"
	"strconv"

	"github.com/lib/pq"
)

/**
* @desc Verifica daca exista sau nu un user in database
* @param $userName (string) : username-ul dupa care se face query-ul
* @return bool -> false daca au aparut erori / nu exista si true daca exista
* @author Mihai Indreias
 */
func findUserByUsername(userName string) bool {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE Username='%s'", userName))
	if err != nil {
		return false
	}
	return userQuery.Next()
}

/**
* @desc Verifica daca exista sau nu un user in database
* @param $userID (int64) : UID-ul dupa care se face query-ul
* @return bool -> false daca au aparut erori / nu exista si true daca exista
* @author Mihai Indreias
 */
func findUserByID(userID int64) bool {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE UID=%v", userID))
	if err != nil {
		return false
	}
	return userQuery.Next()
}

/**
* @desc Gaseste un user in database
* @param $uID (int64) : UID-ul dupa care se face query-ul
* @return User,error -> new User , eroare daca au fost gasite erori sau nu este gasit userul
					 -> User , nil daca exista userul (Field-ul de parola este gol)
* @author Mihai Indreias
*/
func getUserByID(uID int64) (User, error) {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE UID=%v", uID))
	if err != nil {
		return User{}, err
	}
	if userQuery.Next() {
		var us User
		flwi := pq.Int64Array{}
		flwr := pq.Int64Array{}
		userQuery.Scan(&us.UID, &us.InstitutionEmail, &us.PersonalEmail, &us.Username, &us.PasswordHash, &us.YearOfStudy, &us.College, &us.University, &us.Serie, &us.FirstName, &us.LastName, &us.verified, &flwi, &flwr, &us.Major)
		us.FollowingList = flwi
		us.FollowersList = flwr
		return us, nil
	}
	return User{}, errors.New("User does not exist")
}

/**
* @desc Gaseste un user in database
* @param $uername (string) : username-ul dupa care se face query-ul
* @return User,error -> new User , eroare daca au fost gasite erori sau nu este gasit userul
					 -> User , nil daca exista userul (Field-ul de parola este gol)
* @author Mihai Indreias
*/
func getUserByUsername(username string) (User, error) {
	userQuery, err := db.Query(fmt.Sprintf("SELECT * FROM league WHERE Username='%s'", username))
	if err != nil {
		return User{}, err
	}
	if userQuery.Next() {
		var us User
		flwi := pq.Int64Array{}
		flwr := pq.Int64Array{}
		userQuery.Scan(&us.UID, &us.InstitutionEmail, &us.PersonalEmail, &us.Username, &us.PasswordHash, &us.YearOfStudy, &us.College, &us.University, &us.Serie, &us.FirstName, &us.LastName, &us.verified, &flwi, &flwr, &us.Major)
		us.FollowingList = flwi
		us.FollowersList = flwr
		return us, nil
	}
	return User{}, errors.New("User does not exist")
}

/**
* @desc Verifica daca un utilizator se poate autentifica cu "creditentialele" trimise
* @param $uID ,$pass (string) -> user si parola pt query (parola este ne-hash-uita)
* @return bool -> true/fase daca exista in DB un user care are MD5-ul parolei in db egal cu MD5-ul parolei date
* @author Mihai Indreias
 */
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

/**
* @desc Verifica daca un utilizator pote fi adaugat cu "creditentialele" trimise
* @param $newUser (Type User) : noul user cu campurile prinicpale completate , in afara de uID , pe care il dam noi , generat random
								si trimite un mail de verificare pe adresa principala
* @return HTTPResonse -> Code : 409 daca exista deja
					  -> Code : 404 email invalid
					  -> Code : 401 eroare la DB Query
					  -> Code : 200 daca se poate authentifica
* @author Mihai Indreias
*/

func addUser(newUser User) HTTPResponse {
	if findUserByUsername(newUser.Username) {
		return HTTPResponse{Response: []string{"User already exists"}, Code: 409}
	}
	newUser.PasswordHash = passwordType(fmt.Sprintf("%x", md5.Sum([]byte(string(newUser.PasswordHash)))))
	userID := generate18DigitID()
	for findUserByID(userID) {
		userID = generate18DigitID()
	}
	newUser.UID = userID
	if sendVerifcationMail(newUser.InstitutionEmail, "http://35.184.233.76:8080/activate?id="+strconv.FormatInt(userID, 10)) == false {
		return HTTPResponse{Response: []string{"Invalid Email"}, Code: 404}
	}
	sqlStatement := `INSERT INTO league (UID,IEmail,PMail,Username,Password,YearOfStudy,College,University,Major,Serie,FirstName,LastName,verified,following,followers) VALUES($1,$2,$3,$4,$5,$6,$7,$8,$9,$10,$11,$12,$13,$14,$15)`
	_, err := db.Exec(sqlStatement,
		newUser.UID, newUser.InstitutionEmail, newUser.PersonalEmail, newUser.Username,
		newUser.PasswordHash, newUser.YearOfStudy, newUser.College, newUser.University,
		newUser.Major, newUser.Serie, newUser.FirstName, newUser.LastName, 0, pq.Array([]int64{newUser.UID}), pq.Array([]int64{newUser.UID}))
	if err != nil {
		return HTTPResponse{Response: []string{"Error>"}, Code: 401}
	}

	fmt.Println(newUser)
	return HTTPResponse{Response: []string{"User added"}, Code: 200}
}

func checkUserByID(uID, sID int64) error {
	return nil
}

/**
* @desc Verifica daca exista o sesiune
* @param $sID (in64): id-ul sessiunii
* @return bool -> exista sau nu
* @author Mihai Indreias
 */

func seesionExist(sID int64) bool {
	sessionQuery, err := db.Query(fmt.Sprintf("SELECT * FROM sessions WHERE sid=%v", sID))
	if err != nil {
		return false
	}
	return sessionQuery.Next()
}

/**
* @desc Trimite un mesaj in DB
* @param $newMessage (Message)
* @return None
* @author Mihai Indreias
 */

func sendMessage(newMessage Message) {
	sqlStatement := `INSERT INTO messages (id,authorid,text,mediafilepath,date,receiver,typeofreceiver) VALUES($1,$2,$3,$4,$5,$6,$7)`
	_, err := db.Exec(sqlStatement,
		newMessage.ID, newMessage.AuthorID, newMessage.Text, newMessage.MediaFilePath, newMessage.Date, newMessage.Receiver, newMessage.TypeOfReceiver)
	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(newMessage)
}

/**
* @desc Activeaza (verifica) mailul principal al unui user si modifica in DB
* @param $uID(int64) : ID-ul userului verificat
* @return None
* @author Mihai Indreias
 */

func verifyUser(uID int64) {
	db.Exec(`UPDATE league SET verified=1 WHERE uid=$1`, uID)
}

/**
* @desc Gaseste un ID pentru o noua sesiune
* @param None
* @return int64 = id-ul unei noi sesiuni care nu exista deja in DB
* @author Mihai Indreias
 */

func newSessionID() int64 {
	id := generate18DigitID()
	for seesionExist(id) {
		id = generate18DigitID()
	}
	return id
}

/**
* @desc Adauga o sesiune in DB
* @param $sID,$uID (int64) -> sessionID si userID , care descriu sesiunea
* @return None
* @author Mihai Indreias
 */

func addSession(sID, uID int64) {
	db.Exec(`INSERT INTO sessions (sid,uid) VALUES ($1,$2)`, sID, uID)
}

/**
* @desc Verifica daca exista un mesaj in DB
* @param $ID (int64) : ID-ul mesajului dupa care se face Query-ul
* @return bool -> daca exista sau nu
* @author Mihai Indreias
 */

func messageExist(ID int64) bool {
	messageQuery, err := db.Query(fmt.Sprintf("SELECT * FROM messages WHERE id=%v", ID))
	if err != nil {
		return false
	}
	return messageQuery.Next()
}

/**
* @desc Gaseste un ID pentru un nou mesaj
* @param None
* @return int64 = id-ul unui mesaj nou care nu exista deja in DB
* @author Mihai Indreias
 */

func newMessageID() int64 {
	id := generate18DigitID()
	for messageExist(id) {
		id = generate18DigitID()
	}
	return id

}

/**
* @desc Gaseste o sesiune in DB
* @param $sID (in64): id-ul sessiunii
* @return int64 = id-ul userului care are o sesiune sau 0 daca nu exista sesiunea
* @author Mihai Indreias
 */

func getSession(sID int64) int64 {
	sessionQuery, err := db.Query(fmt.Sprintf("SELECT * FROM sessions WHERE sid=%v", sID))
	if err != nil {
		return 0
	}
	sessionQuery.Next()
	var uid int64
	sessionQuery.Scan(&sID, &uid)
	return uid
}

/**
* @desc Gaseste un mesaj in DB
* @param $mID (in64): id-ul mesajului
* @return Message -> mesajul propriu zis (daca exista)
				  -> un mesaj nou (gol) daca nu exista
* @author Mihai Indreias
*/

func getMessageByID(mID int64) Message {
	messageQuery, err := db.Query(fmt.Sprintf("SELECT * FROM messages WHERE id=%v", mID))
	if err != nil {
		return *new(Message)
	}
	messageQuery.Next()
	var m Message
	messageQuery.Scan(&m.ID, &m.AuthorID, &m.Text, &m.MediaFilePath, &m.Date, &m.Receiver, &m.TypeOfReceiver)
	return m
}

/**
* @desc Verifica daca un user da follow altuia
* @param $from , $to (int64) id-urile celor doi utlizatori
* @return bool,int -> true , pozitia in lista de following daca exista $from ii da follow lui $to
				   -> false,0 daca nu exista in lista de following
* @author Mihai Indreias
*/

func isFollowing(from, to int64) (bool, int) {
	fr, _ := getUserByID(from)
	for i, v := range fr.FollowingList {
		if v == to {
			return true, i
		}
	}
	return false, 0
}
