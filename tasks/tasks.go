package tasks

import (
	"net/http"

	utils "github.com/nelsliu9121/wgwheretogoserver/utils"
	"github.com/robfig/cron"
)

var client = &http.Client{}
var fb = utils.Firebase()
var scheduler *cron.Cron

// Tasks Start cron jobs
func Tasks() {
	scheduler = cron.New()

	scheduler.AddFunc("0 0 2 1 * *", GetLocations)
	scheduler.AddFunc("0 0 2 1 * *", GetCourses)
	scheduler.AddFunc("@daily", GetPosts)

	scheduler.Start()
}
