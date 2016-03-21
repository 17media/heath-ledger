package controllers

/*
.__                   __  .__              .__             .___
|  |__   ____ _____ _/  |_|  |__           |  |   ____   __| _/ ____   ___________
|  |  \_/ __ \\__  \\   __\  |  \   ______ |  | _/ __ \ / __ | / ___\_/ __ \_  __ \
|   Y  \  ___/ / __ \|  | |   Y  \ /_____/ |  |_\  ___// /_/ |/ /_/  >  ___/|  | \/
|___|  /\___  >____  /__| |___|  /         |____/\___  >____ |\___  / \___  >__|
     \/     \/     \/          \/                    \/     \/_____/      \/
*/

import (
	"encoding/json"
	"fmt"
	"github.com/17media/heath-ledger/models"
	"github.com/17media/heath-ledger/stores"
	"github.com/julienschmidt/httprouter"
	"github.com/maxwellhealth/bongo"
	"gopkg.in/mgo.v2/bson"
	"net/http"
	"strconv"
	"time"
)

// CreateUser - Create a new user
func CreateUser(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
	decoder := json.NewDecoder(request.Body)
	requestJSON := make(map[string]string)

	err := decoder.Decode(&requestJSON)
	if err != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "")
	}

	user := &models.User{
		Username: requestJSON["username"],
		Email:    requestJSON["email"],
		LastSeen: time.Now(),
	}
	user.SetPassword(requestJSON["password"])

	dbErr := stores.MongoDB.Collection("users").Save(user)
	if dbErr != nil {
		rw.WriteHeader(http.StatusBadRequest)
		fmt.Fprint(rw, "")
	}

	rw.WriteHeader(http.StatusCreated)
	fmt.Fprint(rw, "")
}

// GetUser - Get User by mongo ObjectId
func GetUser(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
	user := &models.User{}
	userID := params.ByName("userID")

	err := stores.MongoDB.Collection("users").FindById(bson.ObjectIdHex(userID), user)
	if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
		rw.WriteHeader(http.StatusNotFound)
		fmt.Fprint(rw, dnfError)
	} else {
		if err != nil {
			rw.WriteHeader(http.StatusBadRequest)
			fmt.Fprint(rw, "")
		} else {
			response, _ := json.Marshal(user)
			rw.Header().Set("Content-Type", "application/json")
			fmt.Fprint(rw, string(response))
		}
	}
}

// ListUsers - User listing
func ListUsers(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
	request.ParseForm()

	users := &[]models.User{}

	limit, _ := strconv.Atoi(request.Form.Get("limit"))
	offset, _ := strconv.Atoi(request.Form.Get("offset"))
	sort := request.Form.Get("sort")

	verified, _ := strconv.ParseBool(request.Form.Get("verified"))

	query := bson.M{
		"verified": verified,
	}

	results := models.ResultSet{stores.MongoDB.Collection("users").Find(query)}

	if sort != "" {
		results.Query.Sort(sort)
	}

	meta, _ := results.Paginate(limit, offset, request.Host, request.URL)
	results.Query.All(users)

	tempJSON := &ListResponse{
		Meta:    meta,
		Objects: users,
	}

	response, _ := json.Marshal(tempJSON)

	rw.Header().Set("Content-Type", "application/json")
	fmt.Fprint(rw, string(response))
}

/*
// updateUser - Update User by ID
func updateUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprint(w, "Update User!\n")
}

// deleteUser - Delete User by ID
func deleteUser(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
        fmt.Fprint(w, "Delete User!\n")
}
*/
