package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
)

// ScheduleAPIResponse Type
type ScheduleAPIResponse struct {
	Office  Location `json:"office"`
	Room    string   `json:"name"`
	Periods struct {
		Monday    []Course `json:"1"`
		Tuesday   []Course `json:"2"`
		Wednesday []Course `json:"3"`
		Thursday  []Course `json:"4"`
		Friday    []Course `json:"5"`
		Saturday  []Course `json:"6"`
		Sunday    []Course `json:"7"`
	} `json:"periods"`
}

var client = &http.Client{}

// GetCourses Get courses from the 3rd-party API
func GetCourses() {
	resp, err := client.Get("http://www.worldgymtaiwan.com/api/schedule_period/schedule?classroom_id=4&office_id=1&month=8")
	if err != nil {
		log.Fatal(err)
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	var jsonBody ScheduleAPIResponse
	json.Unmarshal(body, &jsonBody)
	fmt.Print(jsonBody)
}
