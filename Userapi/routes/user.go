package routes

import (
	"encoding/json"
	"github.com/jameskeane/bcrypt"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"../db"
	"../models"
	"time"
)

func Signup(w http.ResponseWriter, r *http.Request)  {
	user := &models.User{}
	if !user.Validate(w, r) {
		return
	}
	db := database.Db
	c := db.C("checking")
	//inserting into the DB
	user.ID = bson.NewObjectId()
	salt, _ := bcrypt.Salt(10)
	user.Password, _ = bcrypt.Hash(user.Password, salt)

	user.CreatedAt = time.Now().Local()
	user.UpdatedAt = time.Now().Local()

	insertionErrors := c.Insert(&user)

	if insertionErrors != nil{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		resp := map[string]interface{}{"errors": insertionErrors.Error(), "status": 0}
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}else {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{"data": user, "status": 1}
		json.NewEncoder(w).Encode(resp) //nolint:errcheck
	}

}

func Userprofile(w http.ResponseWriter, r *http.Request)  {
	accessToken := &models.AccessToken{}
	if !accessToken.AuthoriseByToken(w,r) {
		return
	}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	user := accessToken.GetUser()
	resp := map[string]interface{}{"data": user, "status": 1}
	json.NewEncoder(w).Encode(resp)
}

func LogOut(w http.ResponseWriter, r *http.Request)  {
	accessToken := &models.AccessToken{}
	if  !accessToken.AuthoriseByToken(w,r) {
		return

	}
	user := accessToken.GetUser()
	accessToken.Remove()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	resp := map[string]interface{}{"data": map[string]interface{}{"user_id":user.ID,"message": "LoggedOut Successfully"}, "status": 1}
	json.NewEncoder(w).Encode(resp)
	
}

func UserUpdate(w http.ResponseWriter, r *http.Request)  {
	accessToken := &models.AccessToken{}
	if !accessToken.AuthoriseByToken(w,r){
		return
	}
	user := &models.User{}

	if user.UpdateValidate( w,r, "create" ) && user.Save(w,r) {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusOK)
		resp := map[string]interface{}{"data": user, "status":1}
		json.NewEncoder(w).Encode(resp)
	}
	
}

// APIInfo : handler function for / call
func ApiInfo(w http.ResponseWriter, r *http.Request) {

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	response := map[string]interface{}{"hello": "Welcome to Grocery app", "status": 1}
	json.NewEncoder(w).Encode(response)
}