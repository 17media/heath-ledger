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

// GeoJSON is used to store geolocation
type GeoJSON struct {
	Type        string    `json:"-"`
	Coordinates []float64 `json:"coordinates"`
}

// User is an information for the stored User object
type User struct {
	ID       bson.ObjectId `bson:"_id,omitempty" json:"id"`
	Email    string        `bson:"email" json:"email"`
	Username string        `bson:"username" json:"username"`
	Password []byte        `bson:"password" json:"password"`
	Created  time.Time     `bson:"created" json:"created"`
	Updated  time.Time     `bson:"updated" json:"updated"`
	LastSeen time.Time     `bson:"last_seen" json:"last_seen"`
	Location GeoJSON       `bson:"location" json:"location"`
}

// HashPassword is used to bcrypt/hash the actual string
func (user *User) HashPassword(password string) {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Couldn't hash password: %v", err)
		panic(err)
	}
	user.Password = hash
}

// newUser is inserting a new user to the db
func newUser(db *mgo.Database, user *User) error {
	user.ID = bson.NewObjectId()
	return db.C("users").Insert(user)
}

// getUser based on the given field
func getUser(db *mgo.Database, field string, value string) (user *User) {
	err := db.C("users").Find(bson.M{field: value}).One(&user)

	if err != nil {
		log.Warningf("User does not exists")
	}
	return
}
