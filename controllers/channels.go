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
)

// CreateChannel - Create a new Channel
func CreateChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        decoder := json.NewDecoder(request.Body)

        channel := &models.Channel{}

        err := decoder.Decode(&channel)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, err.Error())
        }

        saveErr := stores.MongoDB.Collection("channels").Save(channel)
        if saveErr != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, saveErr.Error())
        }

        rw.WriteHeader(http.StatusCreated)
        fmt.Fprint(rw, "")
}

// GetChannel - Get Channel by mongo ObjectId
func GetChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        channel := &models.Channel{}
        channelID := params.ByName("channelID")

        err := stores.MongoDB.Collection("channels").FindById(bson.ObjectIdHex(channelID), channel)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        response, _ := json.Marshal(channel)
                        rw.Header().Set("Content-Type", "application/json")
                        fmt.Fprint(rw, string(response))
                }
        }
}

// ListChannels - Channel listing
func ListChannels(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        request.ParseForm()

        channels := &[]models.Channel{}

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

        results := models.ResultSet{stores.MongoDB.Collection("channels").Find(query)}

        if sort != "" {
                results.Query.Sort(sort)
        }

        meta, _ := results.Paginate(paginationInfo, request.Host, request.URL)
        results.Query.All(channels)

        tempJSON := &ListResponse{
                Meta:    meta,
                Objects: channels,
        }

        response, _ := json.Marshal(tempJSON)

        rw.Header().Set("Content-Type", "application/json")
        fmt.Fprint(rw, string(response))
}

// UpdateChannel - Update Channel by ID
func UpdateChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        decoder := json.NewDecoder(request.Body)
        channel := models.Channel{}

        err := decoder.Decode(&channel)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, err.Error())
        }

        updateQuery := bson.M{"$set": channel}

        channelID := params.ByName("channelID")
        dnfErr := stores.MongoDB.Collection("channels").Collection().UpdateId(bson.ObjectIdHex(channelID), updateQuery)

        if dnfErr != nil {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfErr.Error())
        } else {
                rw.WriteHeader(http.StatusNoContent)
                fmt.Fprint(rw, "")
        }
}

// DeleteChannel - Delete Channel
func DeleteChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        channelID := params.ByName("channelID")
        channel := &models.Channel{}
        err := stores.MongoDB.Collection("channels").FindById(bson.ObjectIdHex(channelID), channel)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        delError := stores.MongoDB.Collection("channels").DeleteDocument(channel)
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

// CreateMessage - Create a new Message
func CreateMessage(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        decoder := json.NewDecoder(request.Body)

        message := &models.Message{}

        err := decoder.Decode(&message)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, err.Error())
        }

        saveErr := stores.MongoDB.Collection("messages").Save(message)
        if saveErr != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, saveErr.Error())
        }

        rw.WriteHeader(http.StatusCreated)
        fmt.Fprint(rw, "")
}

// GetMessage - Get Message by mongo ObjectId
func GetMessage(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        message := &models.Message{}
        messageID := params.ByName("messageID")

        err := stores.MongoDB.Collection("messages").FindById(bson.ObjectIdHex(messageID), message)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        response, _ := json.Marshal(message)
                        rw.Header().Set("Content-Type", "application/json")
                        fmt.Fprint(rw, string(response))
                }
        }
}

// ListMessages - Message listing
func ListMessages(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        request.ParseForm()

        messages := &[]models.Message{}

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

        results := models.ResultSet{stores.MongoDB.Collection("messages").Find(query)}

        if sort != "" {
                results.Query.Sort(sort)
        }

        meta, _ := results.Paginate(paginationInfo, request.Host, request.URL)
        results.Query.All(messages)

        tempJSON := &ListResponse{
                Meta:    meta,
                Objects: messages,
        }

        response, _ := json.Marshal(tempJSON)

        rw.Header().Set("Content-Type", "application/json")
        fmt.Fprint(rw, string(response))
}

// UpdateMessage - Update Message by ID
func UpdateMessage(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        decoder := json.NewDecoder(request.Body)
        message := models.Message{}

        err := decoder.Decode(&message)
        if err != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, err.Error())
        }

        updateQuery := bson.M{"$set": message}

        messageID := params.ByName("messageID")
        dnfErr := stores.MongoDB.Collection("messages").Collection().UpdateId(bson.ObjectIdHex(messageID), updateQuery)

        if dnfErr != nil {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfErr.Error())
        } else {
                rw.WriteHeader(http.StatusNoContent)
                fmt.Fprint(rw, "")
        }
}

// DeleteMessage - Delete Message
func DeleteMessage(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        messageID := params.ByName("messageID")
        message := &models.Message{}
        err := stores.MongoDB.Collection("messages").FindById(bson.ObjectIdHex(messageID), message)
        if dnfError, ok := err.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if err != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, err.Error())
                } else {
                        delError := stores.MongoDB.Collection("messages").DeleteDocument(message)
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
