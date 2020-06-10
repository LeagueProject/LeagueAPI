package main

import (
	"math/rand"
	"net/http"
	"sync"
	"time"

	"github.com/julienschmidt/httprouter"
)

var wg sync.WaitGroup

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

func main() {
	rand.Seed(time.Now().UnixNano())
	serverCRUD := httprouter.New()
	serverCRUD.GET("/read/:key", readHandler)
	serverCRUD.POST("/add/:key", addHandler)
	serverCRUD.GET("/activate", activationHandler)
	serverCRUD.GET("/login", loginHandler)
	wg.Add(1)
	go func() {
		http.ListenAndServe(":8555", serverCRUD)
		wg.Done()
	}()
	wg.Wait()
}
