package utils

import (
	"io/ioutil"

	log "github.com/sirupsen/logrus"
	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Firebase Create a Firebase access point
func Firebase() *firego.Firebase {
	d, err := ioutil.ReadFile("./config/firebase.json")
	if err != nil {
		log.WithError(err).Panic()
	}

	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		log.Fatal(err)
	}

	return firego.New("https://wgwheretogo.firebaseio.com/", conf.Client(oauth2.NoContext))
}
