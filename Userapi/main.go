package main

import (
	"./db"
	"./routes"
	"github.com/gorilla/mux"
	"log"
	"net/http"
)

func main() {
	database.Connect()
	urouter := mux.NewRouter()
	//SignUp
	urouter.HandleFunc("/Signup", routes.Signup).Methods("POST")

	//OAuth2
	urouter.HandleFunc("/authorise", routes.Authorise).Methods("POST")
	urouter.HandleFunc("/accesstoken", routes.AccessToken).Methods("POST")

	//User Profile
	urouter.HandleFunc("/profile", routes.Userprofile).Methods("GET")
	urouter.HandleFunc("/logut", routes.LogOut).Methods("GET")
	//urouter.HandleFunc("/profile",routes.UserUpdate).Methods("PUT")

	//API info
	urouter.HandleFunc("/", routes.ApiInfo).Methods("GET")

	log.Fatal(http.ListenAndServe(":9050", urouter))

}
