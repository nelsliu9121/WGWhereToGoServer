package main

import (
	"net/http"

	log "github.com/sirupsen/logrus"

	"github.com/nelsliu9121/wgwheretogoserver/tasks"
)

func getEverythingEndpoints(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(200)
	tasks.GetLocations()
	go tasks.GetCourses()
	go tasks.GetPosts()
}
func getLocationsEndpoint(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(200)
	tasks.GetLocations()
}
func getCoursesEndpoint(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(200)
	tasks.GetLocations()
	tasks.GetCourses()
}
func getPostsEndpoint(writer http.ResponseWriter, req *http.Request) {
	writer.WriteHeader(200)
	tasks.GetLocations()
	tasks.GetPosts()
}

func main() {
	http.HandleFunc("/everything", getEverythingEndpoints)
	http.HandleFunc("/locations", getLocationsEndpoint)
	http.HandleFunc("/courses", getCoursesEndpoint)
	http.HandleFunc("/posts", getPostsEndpoint)
	if err := http.ListenAndServe(":3030", nil); err != nil {
		log.WithError(err).Panic("Start HTTP Services")
	} else {
		log.Info("HTTP Service Started at port 3030")
	}
}
