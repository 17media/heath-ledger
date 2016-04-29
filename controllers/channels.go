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
        "strings"

        "github.com/17media/heath-ledger/models"
        "github.com/17media/heath-ledger/stores"
        // "github.com/gorilla/context"
        "net/http"
        "strconv"

        "github.com/julienschmidt/httprouter"
        "github.com/maxwellhealth/bongo"
        "gopkg.in/mgo.v2/bson"
)

// ChannelRequest takes in name and particpants to create a channel
type ChannelRequest struct {
        Creator      bson.ObjectId
        Name         string
        Participants []bson.ObjectId
}

// CreateChannel - Create a new Channel
func CreateChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        // user := context.Get(request, "user").(models.User)
        rw.Header().Set("Content-Type", "application/json")
        channelRequest := &ChannelRequest{}
        channel := &models.Channel{}
        decoder := json.NewDecoder(request.Body)
        decodeErr := decoder.Decode(channelRequest)

        users := &[]models.User{}
        usersQuery := bson.M{
                "_id": bson.M{
                        "$in": channelRequest.Participants,
                },
        }

        results := stores.MongoDB.Collection("users").Find(usersQuery)
        results.Query.All(users)

        if decodeErr != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, decodeErr.Error())
        } else {
                query := bson.M{
                        "participants": users,
                }

                mongoErr := stores.MongoDB.Collection("channels").FindOne(query, channel)
                if _, ok := mongoErr.(*bongo.DocumentNotFoundError); ok {
                        //channel.Creator = user
                        channel.Participants = *users
                        saveErr := stores.MongoDB.Collection("channels").Save(channel)
                        if saveErr != nil {
                                rw.WriteHeader(http.StatusBadRequest)
                                fmt.Fprint(rw, saveErr.Error())
                        } else {
                                rw.WriteHeader(http.StatusCreated)
                        }
                } else {
                        rw.WriteHeader(http.StatusOK)
                }
                response, _ := json.Marshal(channel)
                fmt.Fprint(rw, string(response))
        }

}

// GetChannel - Get Channel by mongo ObjectId
func GetChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
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
        rw.Header().Set("Content-Type", "application/json")
        request.ParseForm()

        channels := &[]models.Channel{}
        paginationInfo := models.DefaultPagination()
        query := bson.M{}

        if limit := request.Form.Get("limit"); limit != "" {
                paginationInfo.Limit, _ = strconv.Atoi(limit)
        }
        if offset := request.Form.Get("offset"); offset != "" {
                paginationInfo.Offset, _ = strconv.Atoi(offset)
        }
        sort := request.Form.Get("sort")

        if name := request.Form.Get("name"); name != "" {
                query["name"] = name
        }

        if participants := request.Form.Get("participants"); participants != "" {
                elemMatchParticipants := []bson.M{}
                for _, userID := range strings.Split(participants, ",") {
                        elemMatchParticipants = append(elemMatchParticipants,
                                bson.M{"$elemMatch": bson.M{"_id": bson.ObjectIdHex(userID)}})
                }
                query["participants"] = bson.M{
                        "$all": elemMatchParticipants,
                }
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

        fmt.Fprint(rw, string(response))
}

// UpdateChannel - Update Channel by ID
func UpdateChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")

        channelID := params.ByName("channelID")
        decoder := json.NewDecoder(request.Body)
        channelRequest := ChannelRequest{}
        channel := &models.Channel{}

        mongoErr := stores.MongoDB.Collection("channels").FindById(bson.ObjectIdHex(channelID), channel)
        if dnfError, ok := mongoErr.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                decodeErr := decoder.Decode(&channelRequest)
                if decodeErr != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, decodeErr.Error())
                }

                if channelRequest.Name != "" {
                        channel.Name = channelRequest.Name
                }

                if channelRequest.Creator != "" {
                        creator := models.User{}
                        ceatorErr := stores.MongoDB.Collection("users").FindById(channelRequest.Creator, creator)
                        if creatorNotFound, ok := ceatorErr.(*bongo.DocumentNotFoundError); ok {
                                rw.WriteHeader(http.StatusNotFound)
                                fmt.Fprint(rw, creatorNotFound.Error())
                        }
                        channel.Creator = creator
                }

                if len(channelRequest.Participants) > 0 {
                        participants := []models.User{}
                        usersQuery := bson.M{
                                "_id": bson.M{
                                        "$in": channelRequest.Participants,
                                },
                        }

                        results := stores.MongoDB.Collection("users").Find(usersQuery)
                        results.Query.All(participants)
                        channel.Participants = participants
                }

                updateQuery := bson.M{"$set": channel}

                updateErr := stores.MongoDB.Collection("channels").Collection().UpdateId(bson.ObjectIdHex(channelID), updateQuery)

                if updateErr != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, updateErr.Error())
                } else {
                        response, _ := json.Marshal(&channel)
                        rw.WriteHeader(http.StatusOK)
                        fmt.Fprint(rw, string(response))
                }
        }
}

// DeleteChannel - Delete Channel
func DeleteChannel(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")

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
        rw.Header().Set("Content-Type", "application/json")
        decoder := json.NewDecoder(request.Body)

        message := &models.Message{}

        decodeErr := decoder.Decode(&message)
        if decodeErr != nil {
                rw.WriteHeader(http.StatusBadRequest)
                fmt.Fprint(rw, decodeErr.Error())
        } else {
                saveErr := stores.MongoDB.Collection("messages").Save(message)
                if saveErr != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, saveErr.Error())
                } else {
                        response, _ := json.Marshal(message)
                        rw.WriteHeader(http.StatusCreated)
                        fmt.Fprint(rw, string(response))
                }
        }
}

// GetMessage - Get Message by mongo ObjectId
func GetMessage(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        message := &models.Message{}
        messageID := params.ByName("messageID")

        mongoErr := stores.MongoDB.Collection("messages").FindById(bson.ObjectIdHex(messageID), message)
        if dnfError, ok := mongoErr.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                if mongoErr != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, mongoErr.Error())
                } else {
                        response, _ := json.Marshal(message)
                        rw.Header().Set("Content-Type", "application/json")
                        fmt.Fprint(rw, string(response))
                }
        }
}

// ListMessages - Message listing
func ListMessages(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")
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

        archived, _ := strconv.ParseBool(request.Form.Get("archived"))

        query := bson.M{
                "archived": archived,
        }

        if channel := request.Form.Get("channel"); channel != "" {
                query["channel"] = bson.ObjectIdHex(channel)
        }

        if creatorID := request.Form.Get("creator"); creatorID != "" {
                query["creator"] = bson.M{
                        "$all": bson.M{"$elemMatch": bson.M{"_id": bson.ObjectIdHex(creatorID)}},
                }
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
        fmt.Fprint(rw, string(response))
}

// UpdateMessage - Update Message by ID
func UpdateMessage(rw http.ResponseWriter, request *http.Request, params httprouter.Params) {
        rw.Header().Set("Content-Type", "application/json")

        messageID := params.ByName("messageID")

        newMessage := models.Message{}
        message := models.Message{}

        mongoErr := stores.MongoDB.Collection("messages").FindById(bson.ObjectIdHex(messageID), messageID)
        if dnfError, ok := mongoErr.(*bongo.DocumentNotFoundError); ok {
                rw.WriteHeader(http.StatusNotFound)
                fmt.Fprint(rw, dnfError.Error())
        } else {
                decoder := json.NewDecoder(request.Body)
                decodeErr := decoder.Decode(&newMessage)

                if decodeErr != nil {
                        rw.WriteHeader(http.StatusBadRequest)
                        fmt.Fprint(rw, decodeErr.Error())
                } else {

                        if newMessage.Content != "" {
                                message.Content = newMessage.Content
                        }

                        if newMessage.Sender != "" {
                                sender := models.User{}
                                senderErr := stores.MongoDB.Collection("users").FindById(newMessage.Sender, sender)
                                if senderNotFound, ok := senderErr.(*bongo.DocumentNotFoundError); ok {
                                        rw.WriteHeader(http.StatusNotFound)
                                        fmt.Fprint(rw, senderNotFound.Error())
                                } else {
                                        message.Sender = sender
                                }
                        }

                        if newMessage.Channel != "" {
                                message.Channel = newMessage.Channel
                        }

                        updateQuery := bson.M{"$set": message}
                        updateErr := stores.MongoDB.Collection("messages").Collection().UpdateId(bson.ObjectIdHex(messageID), updateQuery)

                        if updateErr != nil {
                                rw.WriteHeader(http.StatusBadRequest)
                                fmt.Fprint(rw, updateErr.Error())
                        } else {
                                rw.WriteHeader(http.StatusNoContent)
                                response, _ := json.Marshal(message)
                                fmt.Fprint(rw, string(response))
                        }
                }
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
