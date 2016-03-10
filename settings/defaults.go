package settings

/*
.__                   __  .__              .__             .___
|  |__   ____ _____ _/  |_|  |__           |  |   ____   __| _/ ____   ___________
|  |  \_/ __ \\__  \\   __\  |  \   ______ |  | _/ __ \ / __ | / ___\_/ __ \_  __ \
|   Y  \  ___/ / __ \|  | |   Y  \ /_____/ |  |_\  ___// /_/ |/ /_/  >  ___/|  | \/
|___|  /\___  >____  /__| |___|  /         |____/\___  >____ |\___  / \___  >__|
     \/     \/     \/          \/                    \/     \/_____/      \/
*/

import (
	"github.com/spf13/viper"
)

// InitSettings initialize the settings for the app
func InitSettings() {
	// DEBUG mode
	viper.SetDefault("DEBUG", false)

	// Port
	viper.SetDefault("PORT", "8000")

	// AWS Section
	viper.SetDefault("AWS_ACCESS_KEY_ID", "the truth lies beneath")
	viper.SetDefault("AWS_SECRET_KEY", "shh.......")

	// DB Section
	viper.SetDefault("MONGODB_URL", "localhost/test")

	// Cache Section
	viper.SetDefault("REDIS_URL", "localhost:6379/test")

	viper.AutomaticEnv()

	if viper.GetBool("DEVELOPMENT") {
		DevSettings()
	}
}
