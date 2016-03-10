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

// Activity is a generic object for all activities happens in the app
type Activity struct {
	ID        bson.ObjectId `bson:"_id, omitempty" json:"id, omitempty"`
	Type      string        `bson:"type" json:"type" binding:"required"`
	Creator   User          `bson:"creator" json:"creator" binding:"required"`
	Recipient User          `bson:"recipient" json:"recipient" binding:"required"`
	Created   time.Time     `bson:"created" json:"created"`
}
