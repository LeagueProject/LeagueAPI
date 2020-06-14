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
	if querry == "user" { ///By UID
		qID, _ := strconv.ParseInt(r.FormValue("id"), 10, 64)
		user, err := getUserByID(qID)
		var printData []byte
		if err == nil {
			printData, _ = json.Marshal(user)
		} else {
			printData, _ = json.Marshal(HTTPResponse{Response: []string{"Unexistend user"}, Code: 404})
		}
		fmt.Fprintln(w, string(printData))
	} else if querry == "post" {
		fmt.Fprintf(w, "post")
	}
}

func addHandler(w http.ResponseWriter, r *http.Request, p httprouter.Params) {
	querry := p.ByName("key")
	if querry == "user" {
		decoder := json.NewDecoder(r.Body)
		var newU User
		decoder.Decode(&newU)
		fmt.Println(newU)
		response := addUser(newU)
		printData, _ := json.Marshal(response)
		fmt.Fprintln(w, string(printData))
	} else if querry == "post" {

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
			/*
				Change status in db
			*/
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
	if err == nil {
		if canLogin(newU.Username, string(newU.PasswordHash)) {
			printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(us.UID, 10)}, Code: 200})
		} else {
			printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(0, 10)}, Code: 400})
		}
	} else {
		printData, _ = json.Marshal(HTTPResponse{Response: []string{strconv.FormatInt(0, 10)}, Code: 404})
	}
	fmt.Fprintln(w, string(printData))
}
