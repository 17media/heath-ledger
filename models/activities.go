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

// newActivity creates a new Activity
func newActivity(db *mgo.Database, activity *Activity) error {
	now := time.Now()
	activity.ID = bson.NewObjectId()
	activity.Created = now
	return db.C("activities").Insert(activity)
}

// getActivity returns only one Activity from the given query
func getActivity(db *mgo.Database, query bson.M) (activity *Activity) {
	err := db.C("activities").Find(query).One(&activity)

	if err != nil {
		log.Warningf("Activity does not exists")
	}
	return

}

// getActivitys retrieve Activitys from given query
func getActivitys(db *mgo.Database, query bson.M) (activities []*Activity) {
	err := db.C("activities").Find(query).All(&activities)

	if err != nil {
		log.Warningf("No Activitys found from given query")
	}
	return
}

// updateActivity updates a Activity from given Activity instance
func updateActivity(db *mgo.Database, activity *Activity, query bson.M) error {
	now := time.Now()
	query["updated"] = now
	change := bson.M{"$set": query}
	return db.C("activities").Update(activity, change)
}

// deleteActivity deletes a Activity with ID
func deleteActivity(db *mgo.Database, id string) error {
	// Remove Activity
	return db.C("activities").RemoveId(id)
}
