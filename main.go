package main

import (
	"encoding/json"
	"github.com/julienschmidt/httprouter"
	"github.com/maciekmm/curveapi/hooks"
	"github.com/maciekmm/curveapi/models"
	"log"
	"net/http"
	"strconv"
)

type status struct {
	Message string `json:"status"`
}

func userNameHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id := ps.ByName("name")
	prof, err := hooks.GetUserProfileByName(id, false)
	if err != nil {
		displayError(rw, "User not found", 404)
		return
	}
	displayResult(rw, prof)
}

func userIdHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	id, err := strconv.Atoi(ps.ByName("id"))
	if err != nil {
		displayError(rw, "Invalid request", 400)
		return
	}
	prof, err := hooks.GetUserProfile(id, false)
	if err != nil {
		displayError(rw, "User not found", 404)
		return
	}
	displayResult(rw, prof)
}

func redirectHandler(rw http.ResponseWriter, req *http.Request, ps httprouter.Params) {
	http.Redirect(rw, req, "https://github.com/maciekmm/CurveAPI", 301)
}

func displayResult(rw http.ResponseWriter, profile *models.Profile) {
	jsonProfile, err := json.MarshalIndent(profile, "", "    ")
	if err != nil {
		displayError(rw, "Could not marshal profile.", 500)
		return
	}
	rw.Header().Set("Content-Type", "application/json")
	rw.Write(jsonProfile)
}

func displayError(rw http.ResponseWriter, message string, errorCode int) {
	rw.Header().Set("Content-Type", "application/json")
	rw.WriteHeader(errorCode)
	mars, err := json.MarshalIndent(status{message}, "", "    ")
	if err != nil {
		panic("Couldn't parse status json")
	}
	rw.Write(mars)
}

func main() {
	log.Println("Starting up curve-api")
	router := httprouter.New()
	router.GET("/user/:id", userIdHandler)
	router.GET("/username/:name", userNameHandler)
	router.GET("/", redirectHandler)
	log.Fatal(http.ListenAndServe(":2000", router))
}
