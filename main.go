package main

import (
	"fmt"
	"io/ioutil"
	"log"

	"github.com/nelsliu9121/wgwheretogoserver/tasks"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
	"gopkg.in/zabawaba99/firego.v1"
)

func main() {
	d, err := ioutil.ReadFile("./config/firebase.json")
	if err != nil {
		log.Fatal(err)
	}

	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/firebase.database")
	if err != nil {
		log.Fatal(err)
	}
	fb := firego.New("https://wgwheretogo.firebaseio.com/", conf.Client(oauth2.NoContext))

	var tv map[string]interface{}

	if err := fb.Value(&tv); err != nil {
		log.Fatal(err)
	}
	fmt.Println(tv)

	tasks.GetCourses()
}
