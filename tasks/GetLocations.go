package tasks

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"

	utils "github.com/nelsliu9121/WGWhereToGoServer/utils"
	log "github.com/sirupsen/logrus"
)

// Locations received from API
var Locations map[string]Location

var s = utils.Storage{
	Bucket: "wgwheretogo.appspot.com",
}

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
					Long float64 `json:"lng"`
					Lat  float64 `json:"lat"`
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
	channel := make(chan map[string]string)
	go parseLocationImages(channel, jsonBody)
	go putImagesToFirebase(channel)
	go pushLocationsToFirebase(Locations)
}

func pushLocationsToFirebase(locations map[string]Location) {
	if err := fb.Child("Locations").Set(locations); err != nil {
		log.WithError(err).Panic("GetLocations pushToFB")
	} else {
		log.WithFields(log.Fields{"Count": len(locations)}).Info("GetLocations pushToFB")
	}
}

func parseLocation(body officeAPIResponse) map[string]Location {
	locations := make(map[string]Location)
	for _, l := range body.Locations {
		location := Location{
			ID:      l.ID,
			Name:    l.Name,
			Alias:   l.Alias,
			Phone:   l.Phone,
			Zip:     l.Zip,
			Address: l.Address,
			Map:     l.Map,
			Photo:   fmt.Sprintf("Locations/%s.jpg", l.ID),
			City:    l.City,
			Rooms:   l.Rooms,
			Geo:     l.Geocoding[0].Geometry.Location,
		}
		locations[l.ID] = location
	}
	return locations
}

func parseLocationImages(channel chan<- map[string]string, body officeAPIResponse) {
	for _, l := range body.Locations {
		channel <- map[string]string{"ID": l.ID, "Photo": l.Photo}
	}
	close(channel)
}

func putImagesToFirebase(channel chan map[string]string) {
	for p := range channel {
		putImageToFirebase(p["Photo"], p["ID"])
	}
}

func putImageToFirebase(url string, locationID string) string {
	res, err := client.Get(fmt.Sprintf("http://www.worldgymtaiwan.com/imagefly/w330-h240-c/uploads/%s", url))
	if err != nil {
		log.WithError(err).Panic("putImageToFirebase Failed to download image")
	}
	defer res.Body.Close()

	origAttrs, err := s.Attrs(context.Background(), fmt.Sprintf("Locations/%s.jpg", locationID))
	if err == nil && origAttrs.Name != "" {
		log.WithFields(log.Fields{"Name": origAttrs.Name}).Info("putImageToFirebase File exists")
		return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/wgwheretogo.appspot.com/o/%s?alt=media", origAttrs.Name)
	}

	_, attrs, err := s.Put(context.Background(), res.Body, res.Header, fmt.Sprintf("Locations/%s.jpg", locationID))
	if err != nil {
		log.WithError(err).Panic("putImageToFirebase Failed to upload to Firebase")
	} else {
		log.WithFields(log.Fields{"Name": attrs.Name}).Info("putImageToFirebase")
	}
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/wgwheretogo.appspot.com/o/%s?alt=media", attrs.Name)
}
