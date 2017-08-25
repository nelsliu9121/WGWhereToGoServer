package tasks

import (
	"encoding/json"
	"io/ioutil"

	log "github.com/sirupsen/logrus"
)

// Locations received from API
var Locations map[string]Location

type officeAPIResponse struct {
	Locations []struct {
		ID      string `json:"id"`
		Name    string `json:"name"`
		Alias   string `json:"alias"`
		Phone   string `json:"phone"`
		Zip     string `json:"zipcode"`
		Address string `json:"address"`
		Map     string `json:"map"`
		Photo   string `json:"image_path"`
		City    string `json:"zip_city"`
		Rooms   []struct {
			ID   string `json:"id"`
			Name string `json:"name"`
		} `json:"classrooms"`
		Geocoding []struct {
			Geometry struct {
				Location struct {
					Long float32 `json:"lng"`
					Lat  float32 `json:"lat"`
				} `json:"location"`
			} `json:"geometry"`
		} `json:"geocoding"`
	} `json:"data"`
}

// GetLocations Get all locations available from 3rd-party API
func GetLocations() {
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

	Locations = parseLocation(jsonBody)

	if err := fb.Child("Locations").Set(Locations); err != nil {
		log.WithError(err).Panic("GetLocations pushToFB")
	} else {
		log.WithFields(log.Fields{"Count": len(Locations)}).Info("GetLocations pushToFB")
	}
}

func parseLocation(body officeAPIResponse) map[string]Location {
	locations := make(map[string]Location)
	for l := range body.Locations {
		location := Location{
			ID:      l.ID,
			Name:    l.Name,
			Alias:   l.Alias,
			Phone:   l.Phone,
			Zip:     l.Zip,
			Address: l.Address,
			Map:     l.Map,
			Photo:   l.Photo,
			City:    l.City,
			Rooms:   l.Rooms,
			Geo:     l.Geocoding[0].Geometry.Location,
		}
		locations[l.ID] = location
	}
	return locations
}
