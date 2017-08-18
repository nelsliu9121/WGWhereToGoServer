package tasks

// Post Type
type Post struct {
	ID           string `json:"id"`
	Name         string `json:"name"`
	Content      string `json:"content"`
	CreatedAt    string `json:"createdAt"`
	LocationID   string `json:"locationID"`
	CategoryID   string `json:"categoryID"`
	CategoryName string `json:"categoryName"`
}
