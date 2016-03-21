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
	"fmt"
	log "github.com/Sirupsen/logrus"
	"github.com/dchest/uniuri"
	"github.com/maxwellhealth/bongo"
	"golang.org/x/crypto/bcrypt"
	"time"
)

// User is an information for the stored User object
type User struct {
	bongo.DocumentBase `bson:",inline"`
	Email              string    `bson:"email" json:"email" binding:"required"`
	Username           string    `bson:"username" json:"username" binding:"required"`
	Verified           bool      `bson:"verified" json:"verified"`
	Password           []byte    `bson:"-" json:"-" binding:"required"`
	LastSeen           time.Time `bson:"last_seen,omitempty" json:"last_seen,omitempty"`
	Location           []float64 `bson:"location,omitempty" json:"location,omitempty"`
}

var passwordBytes = []byte("ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789@!#$%&*()-=[]{}")

// generatePassword generates a safe password string
func generatePassword() string {
	return uniuri.NewLenChars(16, passwordBytes)
}

// SetPassword sets the password to hash
func (user *User) SetPassword(password string) error {
	hash, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	fmt.Println(user)
	user.Password = hash
	if err != nil {
		log.Fatalf("Couldn't hash password: %v", err)
		return err
	}
	return err
}

// checkPassword takes a string and check against the password hash
func (user User) checkPassword(password string) bool {
	err := bcrypt.CompareHashAndPassword(user.Password, []byte(password))
	if err != nil {
		return false
	}
	return true
}
