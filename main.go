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
	"github.com/17media/heath-ledger/settings"
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

	router := httprouter.New()
	router.GET("/", Index)
	router.GET("/hello/:name", Hello)

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
