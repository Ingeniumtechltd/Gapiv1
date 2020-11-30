package models

import (
	"../db"
	"encoding/json"
	"github.com/asaskevich/govalidator"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"net/url"
	time "time"
)

type Driver struct {
	ID   bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
	Name string        `bson:"name" json:"name"`
	//Username  string        `bson:"username" json:"username"`
	Email         string    `bson:"email" json:"email"`
	Password      string    `bson:"password" json:"password,omitempty"`
	Phone         string    `bson:"phone"   json:"phone"`
	Address       string    `bson:"address" json:"address"`
	Driverlicense string    `bson:"driverlicense json:"driverlicense""`
	CreatedAt     time.Time `bson:"created_at" json:"created_at"`
	UpdatedAt     time.Time `bson:"updated_at" json:"updated_at"`
}

var err error

func (driver *Driver) Validate(w http.ResponseWriter, r *http.Request) bool {
	errs := url.Values{}
	db := database.Db

	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		errs.Add("data", "Invalid data")

		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(resp)
		return false
	}

	if govalidator.IsNull(driver.Name) {
		errs.Add("name", "Name is Required")
	}
	if govalidator.IsNull(driver.Email) {
		errs.Add("email", "Email is Required")
	}
	if govalidator.IsNull(driver.Password) {
		errs.Add("password", "Password is Required")
	}
	if govalidator.IsNull(driver.Phone) {
		errs.Add("phone", "Phone is Required")
	}
	if govalidator.IsNull(driver.Address) {
		errs.Add("address", "Address is required")
	}
	if govalidator.IsNull(driver.Driverlicense) {
		errs.Add("driverlicense", "Driver license is required")
	}

	count, _ := db.C("drvtest").Find(bson.M{"email": driver.Email}).Count()
	if count > 0 {
		errs.Add("email", "E-mail is already in use")
	}
	count, _ = db.C("drvtest").Find(bson.M{"phone": driver.Phone}).Count()
	if count > 0 {
		errs.Add("phone", "Phone number is already in use")
	}
	count, _ = db.C("drvtest").Find(bson.M{"name": driver.Name}).Count()
	if count > 0 {
		errs.Add("name", "Name is already in use")
	}
	count, _ = db.C("drvtest").Find(bson.M{"address": driver.Address}).Count()
	if count > 0 {
		errs.Add("address", "Address is use")
	}
	count, _ = db.C("drvtest").Find(bson.M{"driverlicense": driver.Driverlicense}).Count()
	if count > 0 {
		errs.Add("driverlicense", "Driver license is use")
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

func (driver *Driver) Save(w http.ResponseWriter, r *http.Request) bool {
	db := database.Db
	c := db.C("User")

	driver.UpdatedAt = time.Now().Local()

	if driver.ID == "" {
		driver.ID = bson.NewObjectId()
		driver.CreatedAt = time.Now().Local()
		err = c.Insert(&driver)
	} else {
		err = c.Update(bson.M{"_id": driver.ID}, bson.M{"$set": driver})
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

func (driver *Driver) UpdateValidate(w http.ResponseWriter, r *http.Request, action string) bool {
	errs := url.Values{}
	db := database.Db

	if err := json.NewDecoder(r.Body).Decode(&driver); err != nil {
		errs.Add("data", "Invalid data")
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(resp)
		return false
	}
	if action == "update" && driver.ID != "" {
		old_data := Driver{}
		err = db.C("User").Find(bson.M{"_id": driver.ID}).One(&old_data)
		driver.CreatedAt = old_data.CreatedAt
		if err != nil {
			errs.Add("id", "Invalid Document")
		}
	}
	// what update is been done
	if govalidator.IsNull(driver.Password) {
		errs.Add("password", "New password is required")
	}

	count := 0
	if action == "create" {
		//New record
		count, _ = db.C("User").Find(bson.M{"password": driver.Password}).Count()
	} else {
		//existing Record in the Database
		count, _ = db.C("User").Find(bson.M{"email": driver.Password, "_id": bson.M{"$ne": driver.ID}}).Count()
	}

	if count > 0 {
		errs.Add("password", "Password is in use already")
	}

	if len(errs) > 0 {
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(http.StatusBadRequest)
		resp := map[string]interface{}{"errors": errs, "status": 0}
		json.NewEncoder(w).Encode(resp)
		return false
	}
	return true
}
