package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

type scheduleCourse struct {
	ID      string `json:"id"`
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
	OfficeID string                      `json:"office_id"`
	Room     string                      `json:"name"`
	RoomID   string                      `json:"id"`
	Periods  map[string][]scheduleCourse `json:"periods"`
	Office   struct {
		Name string `json:"name"`
	} `json:"office"`
}

var thisTime = time.Now()
var thisYear = thisTime.Year()
var thisMonth = int(thisTime.Month())

// Schedules entire schedule for all locations
var Schedules map[string]interface{}

// GetCourses Get courses from the 3rd-party API
func GetCourses() {
	Schedules = make(map[string]interface{})
	for k := range Locations {
		location := Locations[k]
		for rk := range location.Rooms {
			room := location.Rooms[rk]
			apiURL := fmt.Sprintf("http://www.worldgymtaiwan.com/api/schedule_period/schedule?classroom_id=%s&office_id=%s&month=%d", room.ID, location.ID, thisMonth)
			resp, err := client.Get(apiURL)
			if err != nil {
				log.WithError(err).Panic("GetCourses FromAPI")
			} else {
				log.WithFields(log.Fields{"Status": resp.StatusCode, "Location": location.ID, "Room": room.ID}).Info("GetCourses FromAPI")
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.WithError(err).Panic("GetCourses Parse JSON")
			}
			var jsonBody scheduleAPIResponse
			json.Unmarshal(body, &jsonBody)

			courses := combCourses(jsonBody)
			pushToFB(courses)
			Schedules[location.ID] = map[string]Courses{room.ID: courses}
		}
	}
}

func combCourses(body scheduleAPIResponse) Courses {
	schedule := make(map[string][]Course)
	for _, day := range body.Periods {
		for ci, c := range day {
			course := Course{
				ID:        c.ID,
				Name:      c.Subject.Name,
				Alias:     c.Subject.Alias,
				Teacher:   c.Teacher,
				Weekday:   string(ci),
				StartTime: c.StartTime,
				EndTime:   c.EndTime,
				OfficeID:  body.OfficeID,
				RoomID:    body.RoomID,
				Month:     thisMonth,
				Year:      thisYear,
			}
			weekday := time.Weekday(ci % 7).String()
			schedule[weekday] = append(schedule[weekday], course)
		}
	}
	courses := Courses{
		OfficeName: body.Office.Name,
		OfficeID:   body.OfficeID,
		RoomName:   body.Room,
		RoomID:     body.RoomID,
		Month:      thisMonth,
		Year:       thisYear,
		Monday:     schedule["Monday"],
		Tuesday:    schedule["Tuesday"],
		Wednesday:  schedule["Wednesday"],
		Thursday:   schedule["Thursday"],
		Friday:     schedule["Friday"],
		Saturday:   schedule["Saturday"],
		Sunday:     schedule["Sunday"],
	}
	return courses
}

func pushToFB(courses Courses) {
	fbURL := fmt.Sprintf("Courses/%d/%d/%s/%s", thisYear, thisMonth, courses.OfficeID, courses.RoomID)
	fb.Child(fbURL).Remove()
	if err := fb.Child(fbURL).Set(courses); err != nil {
		log.WithError(err).Panic("GetCourses pushToFB")
	} else {
		log.WithFields(log.Fields{"Location": courses.OfficeName, "Room": courses.RoomName}).Info("GetCourses pushToFB")
	}
}
