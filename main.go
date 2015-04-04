package main

import (
	"github.com/drone/routes"
	"log"
	"net/http"
    "github.com/maciekmm/curveapi/hooks"
	"strconv"
)
	

func userNameHandler(rw http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get(":name")
	prof, _ := hooks.GetUserProfileByName(id, false)
	routes.ServeJson(rw, prof)
}

func userIdHandler(rw http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(req.URL.Query().Get(":id"))
	prof, _ := hooks.GetUserProfile(id, false)
	routes.ServeJson(rw, prof)
}

func main() {
    log.Println("Starting up curve-api")
	//profile, err := hooks.GetUserProfileByName("maciekmm", true)
	//fmt.Println(profile)
	//fmt.Println(err)
	router := routes.New()
	router.Get("/", userIdHandler)
	router.Get("/user/:id([0-9]+)", userIdHandler)
	router.Get("/user/:name([0-9A-z `._]+)", userNameHandler)
	log.Fatal(http.ListenAndServe(":2000", router))
}
