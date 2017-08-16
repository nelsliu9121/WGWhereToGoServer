package tasks

// Course Type
type Course struct {
	ID        string
	Name      string
	Alias     string
	Teacher   string
	Weekday   string
	StartTime string
	EndTime   string
	OfficeID  string
	RoomID    string
	Month     int
	Year      int
}

// Courses Type
type Courses struct {
	OfficeName string
	OfficeID   string
	RoomName   string
	RoomID     string
	Month      int
	Year       int
	Monday     []Course
	Tuesday    []Course
	Wednesday  []Course
	Thursday   []Course
	Friday     []Course
	Saturday   []Course
	Sunday     []Course
}
