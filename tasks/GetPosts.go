package tasks

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"time"

	log "github.com/sirupsen/logrus"
)

var startOfThisMonth = time.Date(thisYear, time.Month(thisMonth), 1, 0, 0, 0, 0, thisTime.Location())

type postsAPIResponse struct {
	LocationID  string `json:"office_id"`
	RecordCount int    `json:"take"`
	Posts       []struct {
		OfficeID  string `json:"office_id"`
		PostID    string `json:"id"`
		Name      string `json:"title"`
		Content   string `json:"details"`
		CreatedAt string `json:"created_at"`
		Category  struct {
			CategoryID   string `json:"id"`
			CategoryName string `json:"name"`
		} `json:"category"`
	} `json:"data"`
}

// GetPosts reads posts from API and store to Firebase
func GetPosts() {
	channel := make(chan postsAPIResponse)
	go requestPosts(channel)
	for jsonBody := range channel {
		if jsonBody.RecordCount > 0 {
			posts := parsePosts(jsonBody)
			pushPostsToFirebase(jsonBody.LocationID, posts)
		}
	}
}

func requestPosts(channel chan<- postsAPIResponse) {
	for _, location := range Locations {
		apiURL := fmt.Sprintf("http://www.worldgymtaiwan.com/api/post?office_id=%s&sort=-released_at", location.ID)
		resp, err := client.Get(apiURL)
		if err != nil {
			log.WithError(err).Panic("GetPosts FromAPI")
		} else {
			log.WithFields(log.Fields{"Status": resp.StatusCode, "Location": location.ID}).Info("GetPosts FromAPI")
		}
		defer resp.Body.Close()
		body, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			log.WithError(err).Panic("GetCourses Parse JSON")
		}
		var jsonBody postsAPIResponse
		json.Unmarshal(body, &jsonBody)
		channel <- jsonBody
	}
	close(channel)
}

func parsePosts(body postsAPIResponse) []Post {
	var posts []Post
	skippedCount := 0
	for _, p := range body.Posts {
		createdAt, err := time.Parse("2006-01-02 15:04:05 -0700", fmt.Sprintf("%s +0800", p.CreatedAt))
		if err != nil {
			log.WithError(err).Panic("GetPosts parsePost Time Parse")
			continue
		}
		if createdAt.Before(startOfThisMonth) {
			skippedCount++
			continue
		}
		post := Post{
			ID:           p.PostID,
			Name:         p.Name,
			Content:      p.Content,
			CreatedAt:    p.CreatedAt,
			LocationID:   p.OfficeID,
			CategoryID:   p.Category.CategoryID,
			CategoryName: p.Category.CategoryName,
		}
		posts = append(posts, post)
	}
	log.WithFields(log.Fields{"Total": len(body.Posts), "Parsed": len(posts), "Skipped": skippedCount}).Info("GetPosts parsePost")
	return posts
}

func pushPostsToFirebase(locationID string, posts []Post) {
	fbURL := fmt.Sprintf("Posts/%s/%d/%d", locationID, thisYear, thisMonth)
	if err := fb.Child(fbURL).Set(posts); err != nil {
		log.WithError(err).Panic("GetPosts pushPostsToFirebase")
	} else {
		log.WithFields(log.Fields{"Location": locationID}).Info("GetPosts pushPostsToFirebase")
	}
}
