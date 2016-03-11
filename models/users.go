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
	"github.com/dchest/uniuri"
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
	Email    string        `bson:"email" json:"email" binding:"required"`
	Username string        `bson:"username" json:"username" binding:"required"`
	Password []byte        `bson:"password" json:"password" binding:"required"`
	Created  time.Time     `bson:"created" json:"created"`
	Updated  time.Time     `bson:"updated" json:"updated"`
	LastSeen time.Time     `bson:"last_seen" json:"last_seen"`
	Location GeoJSON       `bson:"location" json:"location"`
}

var passwordBytes = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789@!#$%&*()-=[]{}")

// generatePassword generates a safe password string
func generatePassword() string {
	return uniuri.NewLenChars(16, passwordBytes)
}

// hashPassword is used to bcrypt/hash the actual string
func hashPassword(password string) []byte {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		log.Fatalf("Couldn't hash password: %v", err)
		panic(err)
	}
	return hash
}

// checkPassword takes a string and check against the password hash
func (user User) checkPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))

	if err != nil {
		return false
	}

	return true
}

// NewUser creates user based on given parameters
func NewUser(db *mgo.Database, email string, username string, password string) *User {
	id := bson.NewObjectId()
	if username == "" {
		username = id.Hex()
	}
	if password == "" {
		password = generatePassword()
	}

	user := &User{
		ID:       id,
		Email:    email,
		Username: username,
		Password: hashPassword(password),
		Created:  time.Now(),
		Updated:  time.Now(),
	}

	err := db.C("users").Insert(user)

	if err != nil {
		log.Fatalf("Unable to create user")
		panic(err)
	}

	return user
}

// GetUser based on the given field
func GetUser(db *mgo.Database, field string, value string) (user *User) {
	err := db.C("users").Find(bson.M{field: value}).One(&user)

	if err != nil {
		log.Warningf("User does not exists")
	}
	return
}

// GetUsers based on given bson.M query
func GetUsers(db *mgo.Database, query bson.M) (users []*User) {
	err := db.C("users").Find(query).All(&users)

	if err != nil {
		log.Warningf("No Channels found from given query")
	}
	return
}
