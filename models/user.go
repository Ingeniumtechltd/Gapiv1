package models

import (
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	time"time"
	"../db"
)

type User struct {
	ID        bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name      string        `bson:"name" json:"name"`
	//Username  string        `bson:"username" json:"username"`
	Email     string        `bson:"email" json:"email"`
	Password  string        `bson:"password" json:"password,omitempty"`
	Phone     int			`bson:"phone"   json:"phone"`
	CreatedAt time.Time     `bson:"created_at" json:"created_at"`
	UpdatedAt time.Time     `bson:"updated_at" json:"updated_at"`
}

var err error

func (user *User) Validate(w http.ResponseWriter, r *http.Request) bool {
	errs := url.Values{}
	db := database.Db

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil {
		errs.Add("data", "Invalid data")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]interface{}{"errors": errs, "status":0}
		json.NewEncoder(w).Encode(resp)
		return false
	}

	if govalidator.IsNull(user.Name){
		errs.Add("name", "Name is Required")
	}
	if govalidator.IsNull(user.Email){
		errs.Add("email", "Email is Required")
	}
	if govalidator.IsNull(user.Password){
		errs.Add("password", "Password is Required")
	}
	if govalidator.IsNull(string(rune(user.Phone))){
		errs.Add("phone", "Phone is Required")
	}

	count,_ := db.C("checking").Find(bson.M{"email": user.Email}).Count()
	if count > 0 {
		errs.Add("email", "E-mail is already in use")
	}
	count, _ = db.C("checking").Find(bson.M{"phone": user.Phone}).Count()
	if count > 0 {
		errs.Add("phone", "Phone number is already in use")
	}
	count, _ = db.C("checking").Find(bson.M{"name": user.Name}).Count()
	if count > 0{
		errs.Add("name", "Name is already in use")
	}

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		res := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(res)
		return false
	}
	return true
}

func (user *User) Save(w http.ResponseWriter, r *http.Request) bool{
	db := database.Db
	c := db.C("User")

	user.UpdatedAt = time.Now().Local()

	if user.ID == ""{
		user.ID = bson.NewObjectId()
		user.CreatedAt = time.Now().Local()
		err = c.Insert(&user)
	}else {
		err = c.Update(bson.M{"_id": user.ID}, bson.M{"$set": user})
	}

	if err != nil {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusInternalServerError)
		response := map[string]interface{}{"errors": err.Error(), "status": 0}
		json.NewEncoder(w).Encode(response)
		return false

	} else {

		return true
	}

}

func (user *User) UpdateValidate(w http.ResponseWriter, r *http.Request, action string) bool {
	errs := url.Values{}
	db := database.Db

	if err := json.NewDecoder(r.Body).Decode(&user); err != nil{
		errs.Add("data", "Invalid data")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(resp)
		return false
	}
	if action == "update" &&  user.ID != ""{
		old_data := User{}
		err = db.C("User").Find(bson.M{"_id": user.ID}).One(&old_data)
		user.CreatedAt = old_data.CreatedAt
		if err != nil{
			errs.Add("id", "Invalid Document")
		}
	}
	// what update is been done
	if govalidator.IsNull(user.Password){
		errs.Add("password", "New password is required")
	}

	count := 0
	if action == "create"{
		//New record
		count, _ = db.C("User").Find(bson.M{"password": user.Password}).Count()
	}else {
		//existing Record in the Database
		count, _ = db.C("User").Find(bson.M{"email": user.Password, "_id": bson.M{"$ne": user.ID}}).Count()
	}

	if count >0 {
		errs.Add("password", "Password is in use already")
	}

	if len(errs) > 0{
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]interface{}{"errors":errs, "status":0}
		json.NewEncoder(w).Encode(resp)
		return false
	}
	return true
}