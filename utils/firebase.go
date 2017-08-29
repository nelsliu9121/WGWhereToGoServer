package utils

import (
	"context"

	log "github.com/sirupsen/logrus"
	"github.com/zabawaba99/firego"
	"golang.org/x/oauth2/google"
)

// Firebase Create a Firebase access point
func Firebase() *firego.Firebase {
	client, err := google.DefaultClient(context.Background(), "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/firebase.database")

	if err != nil {
		log.WithError(err).Fatal("firebase")
	}

	return firego.New("https://wgwheretogo.firebaseio.com/", client)
}
