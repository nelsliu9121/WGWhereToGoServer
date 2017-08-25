package tasks

// Location Type
type Location struct {
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
	Geo struct {
		Long float32 `json:"lng"`
		Lat  float32 `json:"lat"`
	} `json:"geo"`
}
