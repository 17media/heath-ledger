package models

/*
.__                   __  .__              .__             .___
|  |__   ____ _____ _/  |_|  |__           |  |   ____   __| _/ ____   ___________
|  |  \_/ __ \\__  \\   __\  |  \   ______ |  | _/ __ \ / __ | / ___\_/ __ \_  __ \
|   Y  \  ___/ / __ \|  | |   Y  \ /_____/ |  |_\  ___// /_/ |/ /_/  >  ___/|  | \/
|___|  /\___  >____  /__| |___|  /         |____/\___  >____ |\___  / \___  >__|
     \/     \/     \/          \/                    \/     \/_____/      \/
*/

import (
        log "github.com/Sirupsen/logrus"
        mgo "gopkg.in/mgo.v2"
        "gopkg.in/mgo.v2/bson"
        "time"
)

// Channel is a generic object for all pub/sub purpose
type Channel struct {
        ID      bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
        Name    string        `bson:"name" json:"name"`
        Users   []User        `bson:"users" json:"users"`
        Created time.Time     `bson:"created" json:"created"`
        Updated time.Time     `bson:"updated" json:"updated"`
}

// Message is an archvie of messages in a channel
type Message struct {
        ID       bson.ObjectId `bson:"_id,omitempty" json:"id,omitempty"`
        Sender   User          `bson:"sender" json:"sender" binding:"required"`
        Channel  Channel       `bson:"channel" json:"channel"`
        Created  time.Time     `bson:"created" json:"created"`
        Archived bool          `bson:"archived" json:"archived"`
}

// newChannel creates a new channel
func newChannel(db *mgo.Database, channel *Channel) error {
        now := time.Now()
        channel.ID = bson.NewObjectId()
        channel.Created = now
        channel.Updated = now
        return db.C("channels").Insert(channel)
}

// getChannel returns only one channel from the given query
func getChannel(db *mgo.Database, query bson.M) (channel *Channel) {
        err := db.C("channels").Find(query).One(&channel)

        if err != nil {
                log.Warningf("Channel does not exists")
        }
        return

}

// getChannels retrieve channels from given query
func getChannels(db *mgo.Database, query bson.M) (channels []*Channel) {
        err := db.C("channels").Find(query).All(&channels)

        if err != nil {
                log.Warningf("No Channels found from given query")
        }
        return
}

// updateChannel updates a channel from given channel instance
func updateChannel(db *mgo.Database, channel *Channel, query bson.M) error {
        now := time.Now()
        query["updated"] = now
        change := bson.M{"$set": query}
        return db.C("channels").Update(channel, change)
}

// deleteChannel deletes a channel with ID
func deleteChannel(db *mgo.Database, id string) error {
        // Remove channel
        return db.C("channels").RemoveId(id)
}

// newMessage creates a new message
func newMessage(db *mgo.Database, message *Message) error {
        message.ID = bson.NewObjectId()
        return db.C("messages").Insert(message)
}

// getMessage returns one message from the given bson.M query
func getMessage(db *mgo.Database, query bson.M) (message *Message) {
        err := db.C("messages").Find(query).One(&message)

        if err != nil {
                log.Warningf("Message does not exists")
        }

        return
}

// getMessages retireves messages from the given bson.M query
func getMessages(db *mgo.Database, query bson.M) (messages []*Message) {
        err := db.C("messages").Find(query).All(&messages)

        if err != nil {
                log.Warningf("No messages found")
        }
        return
}

// updateMessage updates Message with the given bson.M query to the provided Message instance
func updateMessage(db *mgo.Database, message *Message, query bson.M) error {
        now := time.Now()
        query["updated"] = now
        change := bson.M{"$set": query}
        return db.C("messages").Update(&message, change)
}

// deleteMessage from given message id
func deleteMessage(db *mgo.Database, id string) error {
        return db.C("messages").RemoveId(id)
}

// archiveMessage is a shortcut helper for archiving a message
func archiveMessage(db *mgo.Database, message *Message) error {
        query := bson.M{"archived": true}
        return updateMessage(db, message, query)
}
