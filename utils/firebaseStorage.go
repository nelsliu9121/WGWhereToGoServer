package utils

import (
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http"

	log "github.com/sirupsen/logrus"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

// Storage Reference
type Storage struct {
	Bucket string
}

func (s *Storage) resource(path string) string {
	return fmt.Sprintf("https://firebasestorage.googleapis.com/v0/b/%s/o?name=%s", s.Bucket, path)
}

func (s *Storage) client() *http.Client {
	d, err := ioutil.ReadFile("./config/firebase.json")
	if err != nil {
		log.WithError(err).Panic()
	}

	conf, err := google.JWTConfigFromJSON(d, "https://www.googleapis.com/auth/userinfo.email", "https://www.googleapis.com/auth/devstorage.read_write")
	if err != nil {
		log.Fatal(err)
	}

	return conf.Client(oauth2.NoContext)
}

func (s *Storage) request(verb string, loc string, data io.Reader) (map[string]interface{}, error) {
	req, err := http.NewRequest(verb, loc, data)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Content-Type", "application/json")

	client := s.client()
	res, err := client.Do(req)
	if err != nil {
		return nil, err
	}

	defer res.Body.Close()

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil, err
	}

	r := new(map[string]interface{})
	err = json.Unmarshal(body, &r)
	if err != nil {
		return nil, err
	}

	return *r, nil
}

// Object will fetch the storage object metadata from Firebase
func (s *Storage) Object(path string) (map[string]interface{}, error) {
	return s.request("GET", s.resource(path), nil)
}

// Put will store a file in Firebase Storage
func (s *Storage) Put(data io.Reader, path string) (map[string]interface{}, error) {
	res, err := s.request("POST", s.resource(path), data)
	if err != nil {
		return nil, err
	}
	return res, err
}
