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
        "gopkg.in/mgo.v2"
)

// Activity is a generic object for all activities happens in the app
type Activity struct {
        bongo.DocumentBase `bson:",inline"`
        Type               string `bson:"type" json:"type" binding:"required"`
        Creator            User   `bson:"creator" json:"creator" binding:"required"`
        Entity             mgo.DBRef
}
