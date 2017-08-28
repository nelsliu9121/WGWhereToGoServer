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
		ID     string `json:"id"`
		Name   string `json:"name"`
		Alias  string `json:"alias"`
		Course struct {
			ID    string `json:"id"`
			Name  string `json:"name"`
			Color string `json:"color"`
		} `json:"course"`
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

// CourseTypes storing unique course types
var CourseTypes map[string]string

// GetCourses Get courses from the 3rd-party API
func GetCourses() {
	Schedules = make(map[string]interface{})
	CourseTypes = make(map[string]string)
	channel := make(chan scheduleAPIResponse)
	go requestCourses(channel)
	for jsonBody := range channel {
		if len(jsonBody.OfficeID) > 0 && len(jsonBody.RoomID) > 0 {
			courses := parseCourses(jsonBody)
			Schedules[jsonBody.OfficeID] = map[string]Courses{jsonBody.RoomID: courses}
			go pushCoursesToFirebase(jsonBody.OfficeID, jsonBody.RoomID, courses)
		}
	}
	pushCourseTypesToFirebase(CourseTypes)
}

func requestCourses(channel chan<- scheduleAPIResponse) {
	for _, location := range Locations {
		for _, room := range location.Rooms {
			apiURL := fmt.Sprintf("http://www.worldgymtaiwan.com/api/schedule_period/schedule?classroom_id=%s&office_id=%s&month=%d", room.ID, location.ID, thisMonth)
			resp, err := client.Get(apiURL)
			if err != nil {
				log.WithError(err).Panic("GetCourses FromAPI")
			} else {
				log.WithFields(log.Fields{"Location": location.ID, "Room": room.ID}).Info("GetCourses FromAPI")
			}
			defer resp.Body.Close()
			body, err := ioutil.ReadAll(resp.Body)
			if err != nil {
				log.WithError(err).Panic("GetCourses Parse JSON")
			}
			var jsonBody scheduleAPIResponse
			json.Unmarshal(body, &jsonBody)
			channel <- jsonBody
		}
	}
	close(channel)
}

func pushCoursesToFirebase(locationID string, roomID string, courses Courses) {
	fbURL := fmt.Sprintf("Courses/%s/%s/%d/%d", locationID, roomID, thisYear, thisMonth)
	if err := fb.Child(fbURL).Set(courses); err != nil {
		log.WithError(err).Panic("GetCourses pushCoursesToFirebase")
	} else {
		log.WithFields(log.Fields{"Location": locationID, "Room": roomID}).Info("GetCourses pushCoursesToFirebase")
	}
}

func pushCourseTypesToFirebase(ct map[string]string) {
	fbURL := "Courses/Types"
	if err := fb.Child(fbURL).Set(ct); err != nil {
		log.WithError(err).Panic("GetCourses pushCourseTypesToFirebase")
	} else {
		log.WithFields(log.Fields{"CourseTypes": len(ct)}).Info("GetCourses pushCourseTypesToFirebase")
	}
}

func parseCourses(body scheduleAPIResponse) Courses {
	schedule := make(map[string][]Course)
	for _, day := range body.Periods {
		for ci, c := range day {
			course := Course{
				ID:            c.ID,
				SubjectID:     c.Subject.ID,
				Name:          c.Subject.Name,
				Alias:         c.Subject.Alias,
				Teacher:       c.Teacher,
				Weekday:       string(ci),
				StartTime:     c.StartTime,
				EndTime:       c.EndTime,
				LocationID:    body.OfficeID,
				RoomID:        body.RoomID,
				CategoryID:    c.Subject.Course.ID,
				CategoryName:  c.Subject.Course.Name,
				CategoryColor: c.Subject.Course.Color,
				Month:         thisMonth,
				Year:          thisYear,
			}
			weekday := time.Weekday(ci % 7).String()
			if course.SubjectID != "" {
				CourseTypes[course.SubjectID] = course.Name
			}
			schedule[weekday] = append(schedule[weekday], course)
		}
	}
	courses := Courses{
		LocationName: body.Office.Name,
		LocationID:   body.OfficeID,
		RoomName:     body.Room,
		RoomID:       body.RoomID,
		Month:        thisMonth,
		Year:         thisYear,
		Monday:       schedule["Monday"],
		Tuesday:      schedule["Tuesday"],
		Wednesday:    schedule["Wednesday"],
		Thursday:     schedule["Thursday"],
		Friday:       schedule["Friday"],
		Saturday:     schedule["Saturday"],
		Sunday:       schedule["Sunday"],
	}
	return courses
}
