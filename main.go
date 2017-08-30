package main

import (
	"net/http"

	"google.golang.org/appengine"

	"github.com/nelsliu9121/WGWhereToGoServer/tasks"
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
	appengine.Main()
}
