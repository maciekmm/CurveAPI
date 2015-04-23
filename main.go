package main

import (
	"github.com/drone/routes"
	"github.com/maciekmm/curveapi/hooks"
	"github.com/maciekmm/curveapi/models"
	"log"
	"net/http"
	"strconv"
)

func userNameHandler(rw http.ResponseWriter, req *http.Request) {
	id := req.URL.Query().Get(":name")
	prof, _ := hooks.GetUserProfileByName(id, false)
	displayResult(rw, prof)
}

func userIdHandler(rw http.ResponseWriter, req *http.Request) {
	id, _ := strconv.Atoi(req.URL.Query().Get(":id"))
	prof, _ := hooks.GetUserProfile(id, false)
	displayResult(rw, prof)
}

func displayResult(rw http.ResponseWriter, profile *models.Profile) {
	if profile == nil {
		resp := struct {
			Status string `json:"status"`
		}{
			"Profile not found",
		}
		rw.WriteHeader(404)
		routes.ServeJson(rw, resp)
		return
	}
	routes.ServeJson(rw, profile)
}

func main() {
	log.Println("Starting up curve-api")
	router := routes.New()
	router.Get("/", userIdHandler)
	router.Get("/user/:id([0-9]+)", userIdHandler)
	router.Get("/user/:name([0-9A-Za-z `._]+)", userNameHandler)
	router.Get("/username/:name([0-9A-Za-z `._]+)", userNameHandler)
	log.Fatal(http.ListenAndServe(":2000", router))
}
