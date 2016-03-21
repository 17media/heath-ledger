package stores

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
	"github.com/maxwellhealth/bongo"
	"github.com/spf13/viper"
	"strings"
)

// MongoDB object for mongodb
var MongoDB *bongo.Connection

// InitMongoDB initializes a mongodb connection
func InitMongoDB() {
	url := strings.Split(viper.GetString("MONGODB_URL"), "/")
	config := &bongo.Config{
		ConnectionString: url[0],
		Database:         url[1],
	}
	conn, err := bongo.Connect(config)

	MongoDB = conn

	if err != nil {
		log.Fatalf("mongoDB connection failed")
		panic(err)
	}
}
