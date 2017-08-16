package tasks

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Locations received from API
var Locations map[string]Location

type officeAPIResponse struct {
	Locations []Location `json:"data"`
}

// GetLocations Get all locations available from 3rd-party API
func GetLocations() map[string]Location {
	resp, err := client.Get("http://www.worldgymtaiwan.com/api/office")
	if err != nil {
		log.WithError(err).Panic("GetLocations")
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.WithError(err).Panic()
	}
	var jsonBody officeAPIResponse
	json.Unmarshal(body, &jsonBody)

	Locations = make(map[string]Location)

	for i := 0; i < len(jsonBody.Locations); i++ {
		d := jsonBody.Locations[i]
		Locations[d.ID] = d
	}

	fb.Child("Locations").Remove()
	if err := fb.Set(map[string]interface{}{"Locations": Locations}); err != nil {
		log.WithError(err).Panic("GetLocations pushToFB")
	} else {
		log.WithFields(log.Fields{"Count": len(Locations)}).Info("GetLocations pushToFB")
	}
	return Locations
}
