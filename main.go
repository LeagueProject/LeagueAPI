package main

import (
	"database/sql"
	"fmt"
	"log"
	"math/rand"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/julienschmidt/httprouter"
)

var wg sync.WaitGroup
var db *sql.DB

/*
	@Note Comunicatia nu este inca secure pentru ca folosim http
	@Note pt api un request arata ceva de genul{
		get ip:port/read/user?id=ID
		-> raspuns un json cu user
		post ip:port/add/user (user ul este descris in body cu ajutorul unui json cu toate campurile)
	}
	----TODO----
	cumparat un domeniu + generat certificat SSL ca sa putem folosi https
	ca sa fie cryptate requesturile
	----END-----
*/

const (
	host     = "35.184.233.76"
	port     = 5432
	user     = "postgres"
	password = "test"
	dbname   = "postgres"
)

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, _ = sql.Open("postgres", psqlInfo)

	rand.Seed(time.Now().UnixNano())
	serverCRUD := httprouter.New()
	serverCRUD.GET("/read/:key", readHandler)
	serverCRUD.POST("/add/:key", addHandler)
	serverCRUD.GET("/activate", activationHandler)
	serverCRUD.POST("/login", loginHandler)
	querryUser, err := db.Query("SELECT * FROM league")
	for querryUser.Next() {
		var us User
		if err := querryUser.Scan(&us.UID, &us.InstitutionEmail, &us.PersonalEmail, &us.Username, &us.PasswordHash, &us.YearOfStudy, &us.College, &us.University, &us.Major, &us.Serie, &us.verified); err != nil {
			log.Fatal(err)
		}
		fmt.Println(us)
	}
	defer db.Close()
	err = db.Ping()
	if err != nil {
		panic(err)
	}
	wg.Add(1)
	go func() {
		http.ListenAndServe(":8080", serverCRUD)
		wg.Done()
	}()
	wg.Wait()
}
