package tasks

// Course Type
type Course struct {
	ID        string `json:"id"`
	Name      string `json:"name"`
	Alias     string `json:"alias"`
	Teacher   string `json:"remark"`
	Weekday   string `json:"day_week"`
	StartTime string `json:"start_at"`
	EndTime   string `json:"end_at"`
	OfficeID  string
	RoomID    string
	Month     int
	Year      int
}

// Courses Type
type Courses struct {
	OfficeID  string
	RoomID    string
	Month     int
	Year      int
	Monday    []Course
	Tuesday   []Course
	Wednesday []Course
	Thursday  []Course
	Friday    []Course
	Saturday  []Course
	Sunday    []Course
}
