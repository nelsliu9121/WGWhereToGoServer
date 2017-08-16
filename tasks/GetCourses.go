package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

type scheduleCourse struct {
	Subject struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	} `json:"subject"`
	Teacher   string `json:"remark"`
	Weekday   string `json:"day_week"`
	StartTime string `json:"start_at"`
	EndTime   string `json:"end_at"`
}

type scheduleAPIResponse struct {
	Office  Location `json:"office"`
	Room    string   `json:"name"`
	Periods struct {
		Monday    []scheduleCourse `json:"1"`
		Tuesday   []scheduleCourse `json:"2"`
		Wednesday []scheduleCourse `json:"3"`
		Thursday  []scheduleCourse `json:"4"`
		Friday    []scheduleCourse `json:"5"`
		Saturday  []scheduleCourse `json:"6"`
		Sunday    []scheduleCourse `json:"7"`
	} `json:"periods"`
}

var thisTime = time.Now()
var thisYear = thisTime.Year()
var thisMonth = int(thisTime.Month())

var Schedule map[string]Courses

// GetCourses Get courses from the 3rd-party API
func GetCourses() map[string]Courses {

	for k := range Locations {
		location := Locations[k]
		for rk := range location.Rooms {
			room := location.Rooms[rk]
			apiURL := fmt.Sprintf("http://www.worldgymtaiwan.com/api/schedule_period/schedule?classroom_id=%s&office_id=%s&month=%d", location.ID, room.ID, thisMonth)
			resp, err := client.Get(apiURL)
			if err != nil {
				log.WithError(err).Panic("GetCourses")
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.WithError(err).Panic("GetCourses")
			}
			var jsonBody scheduleAPIResponse
			json.Unmarshal(body, &jsonBody)

			Schedule = combCourses(jsonBody)
			pushToFB(Schedule)
		}
	}

	return Schedule
}

func combCourses(body scheduleAPIResponse) map[string]Courses {
	
}

func pushToFB(courses map[string]Course) {
	fb.Child(fmt.Sprintf("Courses/%d/%d", thisYear, thisMonth)).Remove()
	if err := fb.Set() {
		
	}
}
