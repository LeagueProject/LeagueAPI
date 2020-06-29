package main

type passwordType string //Am folosit passwordType ca sa ascund parola cand formez un JSON ? bad idea idk

/**
* @desc clasa pentru tipul basic User
* @param $UID = user id : numar random de 18 cifre
		 $InstitutionEmail se foloseste ca primary email , iar
		 $PersonalEmail se va folosesi ca back-up in cazul in care userul nu are cont de mail de la facultate
		 $PasswordHash nu o primim hash-uita o vom hash-ui noi local inainte sa o punem in database (nu am mai modificat filed-ul)
		 $YearOfStudy,$College,$University,$Major,$Serie,$FirstName,$LastName -> informatii personale despre user
		 $verified : true / false , semnifica daca user ul si-a verificat email-ul principal
		 $FollowingList,$FollowersList : liste cu UID-uri care descriu relatiile sociale dintre useri
* @ author Mihai Indreias
*/
type User struct {
	UID              int64        `json:"UID"`
	InstitutionEmail string       `json:"IEmail"`
	PersonalEmail    string       `json:"PMail"`
	Username         string       `json:"Username"`
	PasswordHash     passwordType `json:"Password"`
	YearOfStudy      int32        `json:"YearOfStudy"`
	College          string       `json:"College"`
	University       string       `json:"University"`
	Major            string       `json:"Major"`
	Serie            string       `json:"Serie"`
	FirstName        string       `json:"FirstName"`
	LastName         string       `json:"LastName"`
	verified         int
	FollowingList    []int64 `json:"Following"`
	FollowersList    []int64 `json:"Followers"`
}

/**
* @desc data type pentru raspunsuri la aproape orice web request
* @param $Response : informatii relevante pentru raspuns
		 $Code     : cod (aproximativ de format standard http ex : 2XX-succes ,4XX-eroare)
* @author Mihai Indreias
*/

type HTTPResponse struct {
	Response []string `json:"Response"`
	Code     int64    `json:"Code"`
}

/**
* @desc data type pentru mesaje basic (grupuri de orice fel sau persoane)
* @param $ID : id unic de 18 cifre(generat random) al mesajului
		 $AuthorID : uid-ul autorului
		 $Text , $MediaFilePath - mesaul in sine + path-ul pe server al fisierului media atasat
		 $Date : string reprezentatnd data la care a fost trimis mesajul
		 $Receiver : ID de 18 cifre al unui grup / persoane care a primit mesajul
		 $TypeOfReceiver : string care descrie tipul de receiver ("person"/"group")
* @Note : $TypeOfReceiver ar trebui facut un int8
* @author Mihai Indreias
*/

type Message struct {
	ID             int64  `json:"ID"`
	AuthorID       int64  `json:"AuthorID"`
	Text           string `json:"Text"`
	MediaFilePath  string `json:"Media"`
	Date           string `json:"Date"`
	Receiver       int64  `json:"Receiver"`
	TypeOfReceiver string `json:"TypeOfReceiver"`
}
