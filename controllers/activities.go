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

        "github.com/17media/heath-ledger/models"
        "github.com/17media/heath-ledger/stores"
        "github.com/julienschmidt/httprouter"
        "github.com/maxwellhealth/bongo"
        "gopkg.in/mgo.v2/bson"
)

// CreateActivity - Create a new Activity
func CreateActivity(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        decoder := json.NewDecoder(request.Body)

        activity := &models.Activity{}

        decodeErr := decoder.Decode(&activity)
        if decodeErr != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, decodeErr.Error())
        }

        saveErr := stores.MongoDB.Collection("activitys").Save(activity)
        if saveErr != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, saveErr.Error())
        }

        response, _ := json.Marshal(activity)
        rw.WriteHeader(http.StatusCreated)
        fmt.Fprint(rw, string(response))
}

// GetActivity - Get Activity by mongo ObjectId
func GetActivity(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        activity := &models.Activity{}
        activityID := params.ByName("activityID")

        err := stores.MongoDB.Collection("activitys").FindById(bson.ObjectIdHex(activityID), activity)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        response, _ := json.Marshal(activity)
                        rw.Header().Set("Content-Type", "application/json")
                        fmt.Fprint(rw, string(response))
                }
        }
}

// ListActivities - Activity listing
func ListActivities(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        request.ParseForm()

        activitys := &[]models.Activity{}

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

        results := models.ResultSet{stores.MongoDB.Collection("activitys").Find(query)}

        if sort != "" {
                results.Query.Sort(sort)
        }

        meta, _ := results.Paginate(paginationInfo, request.Host, request.URL)
        results.Query.All(activitys)

        tempJSON := &ListResponse{
                Meta:    meta,
                Objects: activitys,
        }

        response, _ := json.Marshal(tempJSON)

        rw.Header().Set("Content-Type", "application/json")
        fmt.Fprint(rw, string(response))
}

// UpdateActivity - Update Activity by ID
func UpdateActivity(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        decoder := json.NewDecoder(request.Body)
        activity := models.Activity{}

        err := decoder.Decode(&activity)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, err.Error())
        }

        updateQuery := bson.M{"$set": activity}

        activityID := params.ByName("activityID")
        dnfErr := stores.MongoDB.Collection("activitys").Collection().UpdateId(bson.ObjectIdHex(activityID), updateQuery)

        if dnfErr != nil {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfErr.Error())
        } else {
                rw.WriteHeader(http.StatusNoContent)
                fmt.Fprint(rw, "")
        }
}

// DeleteActivity - Delete Activity
func DeleteActivity(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
        activityID := params.ByName("activityID")
        activity := &models.Activity{}
        err := stores.MongoDB.Collection("activitys").FindById(bson.ObjectIdHex(activityID), activity)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        delError := stores.MongoDB.Collection("activitys").DeleteDocument(activity)
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
