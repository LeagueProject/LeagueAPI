package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

/*
	----TODO----
	Sa primesc in request si UID si SessionID al userului care face request-ul
	Verifica daca sunt compatibile UID si SID si dupa verific ce date poate sa primeasca UID-ul respectiv
	despre UID-ul din request ( grupuri / houses /etc...)
	----END------

	NEVERMIND ??
	Problema era cu grupurile , dar cred ca aia se face in cadrul structurii de grup...
*/
func readHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	querry := p.ByName("key")
	var printData []byte
	if querry == "user" { ///By UID
		qID, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		user, err := getUserByID(qID)
		if err == nil {
			printData, _ = json.Marshal(user)
		} else {
			printData, _ = json.Marshal(HTTPResponse{Response: []string{"Unexistend user"}, Code: 404})
		}
	} else if querry == "message" {
		sID, _ := strconv.ParseInt(r.FormValue("sid"), 10, 64)
		uID, _ := strconv.ParseInt(r.FormValue("uid"), 10, 64)
		mID, _ := strconv.ParseInt(r.FormValue("mid"), 10, 64)
		if uID == getSession(sID) {
			message := getMessageByID(mID)
			if message.TypeOfReceiver == "person" {
				if message.AuthorID == uID || message.Receiver == uID {
					printData, _ = json.Marshal(message)
				} else {
					printData, _ = json.Marshal(*new(User))
				}
			}
		} else {
			printData, _ = json.Marshal(new(Message))
		}
	}
	fmt.Fprintln(w, string(printData))

}

func addHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	querry := p.ByName("key")
	decoder := json.NewDecoder(r.Body)
	if querry == "user" {
		var newU User
		decoder.Decode(&newU)
		response := addUser(newU)
		printData, _ := json.Marshal(response)
		fmt.Fprintln(w, string(printData))
	} else if querry == "message" {
		var newMessage Message
		decoder.Decode(&newMessage)
		uID := newMessage.AuthorID
		sID, _ := strconv.ParseInt(r.FormValue("sid"), 10, 64)

		err := checkUserByID(uID, sID)
		var printData []byte
		if err == nil {
			newMessage.ID = newMessageID()
			sendMessage(newMessage)
			printData, _ = json.Marshal(HTTPResponse{Response: []string{"Sent"}, Code: 200})
		} else {
			printData, _ = json.Marshal(HTTPResponse{Response: []string{"Not Sent"}, Code: 401})
		}
		fmt.Fprintln(w, string(printData))
	}
}

func activationHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	id, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
	user, err := getUserByID(id)
	var printData []byte
	if err != nil {
		printData, _ = json.Marshal(HTTPResponse{Response: []string{"User does not exist"}, Code: 404})
	} else {
		if user.verified == true {
			printData, _ = json.Marshal(HTTPResponse{Response: []string{"User already activated"}, Code: 304})
		} else {
			verifyUser(id)
			printData, _ = json.Marshal(HTTPResponse{Response: []string{"User verified"}, Code: 202})
		}
	}
	fmt.Println("handler user ", id)
	fmt.Fprintln(w, string(printData))
}

func loginHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	decoder := json.NewDecoder(r.Body)
	var newU User
	decoder.Decode(&newU)
	us, err := getUserByUsername(newU.Username)
	var printData []byte
	fmt.Println(us)
	if us.verified == false {
		printData, _ = json.Marshal(HTTPResponse{Response: []string{"0", "Not verified"}, Code: 404})
	} else if err == nil {
		if canLogin(newU.Username, string(newU.PasswordHash)) {
			sID := newSessionID()
			addSession(sID, us.UID)
			printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(us.UID, 10), strconv.FormatInt(sID, 10)}, Code: 200})

		} else {
			printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(0, 10)}, Code: 400})
		}
	} else {
		printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(0, 10)}, Code: 404})
	}
	fmt.Fprintln(w, string(printData))
}

func sessionValidHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sid, _ := strconv.ParseInt(r.FormValue("sid"), 10, 64)
	uid, _ := strconv.ParseInt(r.FormValue("uid"), 10, 64)
	var printData []byte
	if uid == getSession(sid) {
		printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(1, 10)}, Code: 200})
	} else {
		printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(9, 10)}, Code: 404})
	}
	fmt.Fprintln(w, string(printData))
}

func followStatusHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	sid, _ := strconv.ParseInt(r.FormValue("sid"), 10, 64)
	fr, _ := strconv.ParseInt(r.FormValue("from"), 10, 64)
	to, _ := strconv.ParseInt(r.FormValue("to"), 10, 64)
	var printData []byte
	if fr == getSession(sid) {
		if p.ByName("key") == "follow" {
			ok, _ := isFollowing(to, fr)
			if ok {
				printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(1, 10)}, Code: 201})
			} else {
				if findUserByID(to) {
					db.Exec(`UPDATE league SET following = array_append(following,$1) WHERE uid=$2`, to, fr)
					db.Exec(`UPDATE league SET followers = array_append(followers,$2) WHERE uid=$1`, to, fr)
					printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(1, 10)}, Code: 202})
				} else {
					printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(0, 10)}, Code: 404})
				}
			}
		} else if p.ByName("key") == "unfollow" {
			ok, _ := isFollowing(to, fr)
			printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(1, 10)}, Code: 201})
			if ok {
				db.Exec(`UPDATE league SET following = array_remove(following,$2) WHERE uid=$1`, to, fr)
				db.Exec(`UPDATE league SET followers = array_remove(followers,$1) WHERE uid=$2`, to, fr)
			}
		}

	} else {
		printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(0, 10)}, Code: 404})
	}
	fmt.Fprintln(w, string(printData))
}
