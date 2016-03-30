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
        "net/http"
        "strconv"
        "time"

        "github.com/17media/heath-ledger/models"
        "github.com/17media/heath-ledger/stores"
        "github.com/julienschmidt/httprouter"
        "github.com/maxwellhealth/bongo"
        "gopkg.in/mgo.v2/bson"
)

// CreateUser - Create a new user
func CreateUser(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        decoder := json.NewDecoder(request.Body)
        requestJSON := make(map[string]string)

        err := decoder.Decode(&requestJSON)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, err.Error())
        }

        user := &models.User{
                Username: requestJSON["username"],
                Email:    requestJSON["email"],
                LastSeen: time.Now(),
        }
        user.SetPassword(requestJSON["password"])

        saveErr := stores.MongoDB.Collection("users").Save(user)
        if saveErr != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, saveErr.Error())
        }

        //responseJSON = map[string]string

        //responseJSON["access_token"] =

        rw.WriteHeader(http.StatusCreated)
        fmt.Fprint(rw, "")
}

// GetUser - Get User by mongo ObjectId
func GetUser(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        user := &models.User{}
        userID := params.ByName("userID")

        err := stores.MongoDB.Collection("users").FindById(bson.ObjectIdHex(userID), user)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        response, _ := json.Marshal(user)
                        rw.Header().Set("Content-Type", "application/json")
                        fmt.Fprint(rw, string(response))
                }
        }
}

// ListUsers - User listing
func ListUsers(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        request.ParseForm()

        users := &[]models.User{}

        paginationInfo := models.DefaultPagination()

        if limit := request.Form.Get("limit"); limit != "" {
                paginationInfo.Limit, _ = strconv.Atoi(limit)
        }
        if offset := request.Form.Get("offset"); offset != "" {
                paginationInfo.Offset, _ = strconv.Atoi(offset)
        }
        sort := request.Form.Get("sort")

        verified, _ := strconv.ParseBool(request.Form.Get("verified"))

        query := bson.M{
                "verified": verified,
        }

        results := models.ResultSet{stores.MongoDB.Collection("users").Find(query)}

        if sort != "" {
                results.Query.Sort(sort)
        }

        meta, _ := results.Paginate(paginationInfo, request.Host, request.URL)
        results.Query.All(users)

        tempJSON := &ListResponse{
                Meta:    meta,
                Objects: users,
        }

        response, _ := json.Marshal(tempJSON)

        rw.Header().Set("Content-Type", "application/json")
        fmt.Fprint(rw, string(response))
}

// UpdateUser - Update User by ID
func UpdateUser(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        decoder := json.NewDecoder(request.Body)
        requestJSON := make(map[string]string)

        err := decoder.Decode(&requestJSON)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, err.Error())
        }

        updateQuery := bson.M{"$set": requestJSON}

        userID := params.ByName("userID")
        dnfErr := stores.MongoDB.Collection("users").Collection().UpdateId(bson.ObjectIdHex(userID), updateQuery)

        if dnfErr != nil {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfErr.Error())
        } else {
                rw.WriteHeader(http.StatusNoContent)
                fmt.Fprint(rw, "")
        }
}

// DeleteUser - Delete User
func DeleteUser(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        userID := params.ByName("userID")
        user := &models.User{}
        err := stores.MongoDB.Collection("users").FindById(bson.ObjectIdHex(userID), user)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        delError := stores.MongoDB.Collection("users").DeleteDocument(user)
                        if delError != nil {
                                rw.WriteHeader(http.StatusBadRequest)
                                fmt.Fprint(rw, delError.Error())
                        } else {
                                rw.WriteHeader(http.StatusNoContent)
                                fmt.Fprint(rw, "")

                        }
                }
        }
}
