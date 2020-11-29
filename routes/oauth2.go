package routes

import (
	"../models"
	"encoding/json"
	"net/http"
)

//Authorise handler function for authorise
func Authorise(w http.ResponseWriter, r *http.Request)  {

	authRequest := &models.AuthoriseRequestBody{}
	if !authRequest.Validate(w,r){
		return
	}
	authcode := authRequest.GenerateAuthCode(w)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{"data": map[string]interface{}{
		"authorisation_code" : authcode.Code,
		"expires_at" : authcode.ExpiresAt,
	}, "status": 1}
	json.NewEncoder(w).Encode(response)
}

//AccessToken
func AccessToken(w http.ResponseWriter, r *http.Request)  {
	accessTokenRequest := &models.AccessTokenRequestBody{}
	if !accessTokenRequest.Validate(w, r){
		return
	}
	accesstoken := accessTokenRequest.GenerateAccessToken(w)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)

	response := map[string]interface{}{"data": map[string]interface{}{
		"access_token": accesstoken.Token,
		"expirees_at": accesstoken.ExpiresAt,
	}, "status":1}
	json.NewEncoder(w).Encode(response)

}