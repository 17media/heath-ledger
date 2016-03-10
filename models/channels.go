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
	"golang.org/x/crypto/bcrypt"
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

// NewChannel creates a new channel
func newChannel(db *mgo.Database, channel *Channel) error {
	now := time.Now()
	channel.ID = bson.NewObjectId()
	channel.Created = now
	channel.Updated = now
	return db.C("channels").Insert(channel)
}

// GetChannelByID is what it is
func getChannelByID(db *mgo.Database, id string) (channel *Channel) {
	err := db.C("channels").Find(bson.M{"_id": id}).One(&channel)

	if err != nil {
		log.Warningf("Channel does not exists")
	}
	return

}

// GetChannels retrieve channels from given query
func getChannels(db *mgo.Database, field string, value string) (channels []*Channel) {
	err := db.C("channels").Find(bson.M{field: value}).All(&channels)

	if err != nil {
		log.Warningf("No Channels found from given query")
	}
	return
}

// UpdateChannel updates a channel from given channel instance
func updateChannel(db *mgo.Database, channel *Channel) error {
	now := time.Now()
	change := bson.M{"$set": bson.M{"phone": "+86 99 8888 7777", "timestamp": now}}
	return db.C("channels").Update(channel, change)
}

// deleteChannel deletes a channel with ID
func deleteChannel(db *mgo.Database, id string) error {
	// Remove channel
	return db.C("channels").RemoveId(id)
}

// NewMessage creates a new message
func newMessage(db *mgo.Database, message *Message) error {
	message.ID = bson.NewObjectId()
	return db.C("messages").Insert(message)
}
