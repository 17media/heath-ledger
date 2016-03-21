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
)

// Channel is a generic object for all pub/sub purpose
type Channel struct {
        bongo.DocumentBase `bson:",inline"`
        Name               string `bson:"name" json:"name"`
        Users              []User `bson:"users" json:"users"`
}

// Message is an archvie of messages in a channel
type Message struct {
        bongo.DocumentBase `bson:",inline"`
        Sender             User    `bson:"sender" json:"sender" binding:"required"`
        Channel            Channel `bson:"channel" json:"channel"`
        Archived           bool    `bson:"archived" json:"archived"`
}
