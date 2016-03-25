package main

/*
.__                   __  .__              .__             .___
|  |__   ____ _____ _/  |_|  |__           |  |   ____   __| _/ ____   ___________
|  |  \_/ __ \\__  \\   __\  |  \   ______ |  | _/ __ \ / __ | / ___\_/ __ \_  __ \
|   Y  \  ___/ / __ \|  | |   Y  \ /_____/ |  |_\  ___// /_/ |/ /_/  >  ___/|  | \/
|___|  /\___  >____  /__| |___|  /         |____/\___  >____ |\___  / \___  >__|
     \/     \/     \/          \/                    \/     \/_____/      \/
*/

import (
	"fmt"
	"github.com/17media/heath-ledger/controllers"
	"github.com/17media/heath-ledger/settings"
	"github.com/17media/heath-ledger/stores"
	"github.com/codegangsta/negroni"
	"github.com/facebookgo/grace/gracehttp"
	"github.com/julienschmidt/httprouter"
	"github.com/spf13/viper"
	"net/http"
)

// Index function
func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, viper.GetBool("DEBUG"))
}

// Hello function
func Hello(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprintf(w, "hello, %s!\n", ps.ByName("name"))
}

func main() {
	settings.InitSettings()
	stores.InitMongoDB()
	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

	// User api
	router.POST("api/1//users/", controllers.CreateUser)
	router.GET("api/1/users/", controllers.ListUsers)
	router.GET("api/1/users/:userID/", controllers.GetUser)
	router.POST("api/1/users/:userID/", controllers.UpdateUser)
	router.PATCH("api/1/users/:userID/", controllers.UpdateUser)
	router.DELETE("api/1//users/:userID/", controllers.DeleteUser)

	// Channel api
	router.POST("api/1/channels/", controllers.CreateChannel)
	router.GET("api/1/channels/", controllers.ListChannels)
	router.GET("api/1/channels/:channelID/", controllers.GetChannel)
	router.POST("api/1/channels/:channelID/", controllers.UpdateChannel)
	router.PATCH("api/1/channels/:channelID/", controllers.UpdateChannel)
	router.DELETE("api/1/channels/:channelID/", controllers.DeleteChannel)

	// Message api
	router.POST("api/1/messages/", controllers.CreateMessage)
	router.GET("api/1/messages/", controllers.ListMessages)
	router.GET("api/1/messages/:messageID/", controllers.GetMessage)
	router.POST("api/1/messages/:messageID/", controllers.UpdateMessage)
	router.PATCH("api/1/messages/:messageID/", controllers.UpdateMessage)
	router.DELETE("api/1/messages/:messageID/", controllers.DeleteMessage)

	// Activity api
	router.POST("api/1/activities/", controllers.CreateActivity)
	router.GET("api/1/activities/", controllers.ListActitivities)
	router.GET("api/1/activities/:activityID/", controllers.GetActivity)
	router.POST("api/1/activities/:activityID/", controllers.UpdateActivity)
	router.PATCH("api/1/activities/:activityID/", controllers.UpdateActivity)
	router.DELETE("api/1/activities/:activityID/", controllers.DeleteActivity)

	// Middleware
	n := negroni.Classic()
	n.UseHandler(router)

	// Start server
	gracehttp.Serve(
		&http.Server{
			Addr:    ":" + viper.GetString("PORT"),
			Handler: n,
		},
	)

}
