package main

type passwordType string

/*
	@Note : stochez doar hashul parolei pe MD5
	@Note : hasul parolei nu intra la marshal , dar poate fi "unmarshel"ed dintr un
			text json
	@Note : profilurile sunt publice
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
	verified         bool
}

/*
	@Note : helper struct for Create/Delete/Update requests
	@Note : folosesc coduri HTTP pt status
*/
type HTTPResponse struct {
	Response string `json:"Response"`
	Code     int64  `json:"Code"`
}
