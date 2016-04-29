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
        "github.com/maxwellhealth/bongo"
        "gopkg.in/mgo.v2/bson"
)

// Channel is a generic object for all pub/sub purpose
type Channel struct {
        bongo.DocumentBase `bson:",inline"`
        Name               string `bson:"name" json:"name"`
        Creator            User   `bson:"creator" json:"creator" binding:"required"`
        Participants       []User `bson:"participants" json:"participants" binding:"required"`
}

// Message is an archive of messages in a channel
type Message struct {
        bongo.DocumentBase `bson:",inline"`
        Sender             User          `bson:"sender" json:"sender" binding:"required"`
        Channel            bson.ObjectId `bson:"channel" json:"channel"`
        Content            string        `bson:"content" json:"content"`
        Archived           bool          `bson:"archived" json:"archived"`
}
