package main

import (
	"database/sql"
	"fmt"
	"math/rand"
	"net/http"
	"sync"
	"time"

	_ "github.com/lib/pq"

	"github.com/julienschmidt/httprouter"
)

var wg sync.WaitGroup //Wait Group ca sa nu se termine executia programului
var db *sql.DB        //Data Base

/*
	@Note Comunicatia nu este inca secure pentru ca folosim http
	@Note pt api un request arata ceva de genul{
		get ip:port/read/user?id=ID
		-> raspuns un json cu user
		post ip:port/add/user (user ul este descris in body cu ajutorul unui json cu toate campurile)
	}
	@Note For Future reference : {
		Database-ul are 3 table-uri :"league","messages","sessions" deocamdata
		in "league" sunt tinuti userii
	}
	@Note : un tool bun pentru testing pe localhost sau pe server este POSTMAN :
			te lasa sa creezi requesturi http custom
	----TODO----
	cumparat un domeniu + generat certificat SSL ca sa putem folosi https
	ca sa fie cryptate requesturile
	----END-----
*/

/**
* @desc constante pentru database
* @param $host/$port string     -> ip/port pentru conectare
		 $user/$password/$dbane -> "creditentiale" pentru logare / acces la table-uri
* @return none
* @ author Mihai Indreias
*/

const (
	host     = "35.184.233.76"
	port     = 5432
	user     = "postgres"
	password = "test@LEAGUEINC"
	dbname   = "postgres"
)

/**
* @desc realizeaza conexiunea la database si configureaza web handlers
		ruleaza pe port 8080
* @param none
* @return none
* @ author Mihai Indreias
*/

func main() {
	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s "+"password=%s dbname=%s sslmode=disable", host, port, user, password, dbname)
	db, _ = sql.Open("postgres", psqlInfo)

	rand.Seed(time.Now().UnixNano())
	serverCRUD := httprouter.New()
	serverCRUD.GET("/read/:key", readHandler)
	serverCRUD.POST("/add/:key", addHandler)
	serverCRUD.GET("/activate", activationHandler)
	serverCRUD.POST("/login", loginHandler)
	serverCRUD.POST("/check/session", sessionValidHandler)
	serverCRUD.POST("/followStatus/:key", followStatusHandler)
	defer db.Close()
	err := db.Ping()
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
