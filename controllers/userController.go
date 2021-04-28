package controllers

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
	"github.com/hienviluong125/trello_clone_api/database"
	"github.com/hienviluong125/trello_clone_api/models"
)

type JsonResponse map[string]interface{}
type UserController struct{}

func (controller *UserController) Register(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &params)

	var user = models.User{Email: params["email"]}
	hashedPassword, err := user.HashPassword(params["password"])
	user.HashedPassword = hashedPassword

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": err.Error()})
		return
	}

	err = user.Validate()

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": err.Error()})
		return
	}

	err = database.DBConn.Create(&user).Error

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": err.Error()})
		return
	}

	var token string
	token, err = user.GenerateToken()

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(JsonResponse{"success": true, "token": token})
}

func (controller *UserController) Login(w http.ResponseWriter, r *http.Request) {
	var params map[string]string
	reqBody, _ := ioutil.ReadAll(r.Body)
	json.Unmarshal(reqBody, &params)

	var user models.User
	err := database.DBConn.Model(models.User{}).Where("email=?", params["email"]).Take(&user).Error

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "error": err.Error()})
		return
	}

	err = user.VerifyPassword(string(user.HashedPassword), params["password"])

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "error": err.Error()})
		return
	}

	token, err := user.GenerateToken()

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "error": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(JsonResponse{"success": true, "token": token})
}

func (controller *UserController) Profile(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "jwt").(jwt.MapClaims)
	user_id := (claims)["user_id"]

	var user models.User
	err := database.DBConn.Select("name, email").Find(&user, fmt.Sprint(user_id)).Error

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(JsonResponse{"success": true, "user": user})
}

func (controller *UserController) UpdateProfile(w http.ResponseWriter, r *http.Request) {
	claims := context.Get(r, "jwt").(jwt.MapClaims)
	user_id := (claims)["user_id"]

	reqBody, _ := ioutil.ReadAll(r.Body)

	var user models.User
	var err error
	err = database.DBConn.Find(&user, fmt.Sprint(user_id)).Error

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": err.Error()})
		return
	}

	json.Unmarshal(reqBody, &user)

	err = database.DBConn.Model(&user).Updates(&user).Error

	if err != nil {
		json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": err.Error()})
		return
	}

	json.NewEncoder(w).Encode(JsonResponse{"success": true, "user": user})
}
