package tasks

// Course Type
type Course struct {
	ID            string `json:"id"`
	SubjectID     string `json:"subjectId"`
	Name          string `json:"name"`
	Alias         string `json:"alias"`
	Teacher       string `json:"teacher"`
	Weekday       string `json:"weekday"`
	StartTime     string `json:"startTime"`
	EndTime       string `json:"endTime"`
	LocationID    string `json:"locationID"`
	RoomID        string `json:"roomID"`
	CategoryID    string `json:"categoryID"`
	CategoryName  string `json:"categoryName"`
	CategoryColor string `json:"categoryColor"`
	Month         int    `json:"month"`
	Year          int    `json:"year"`
}

// Courses Type
type Courses struct {
	LocationName string   `json:"locationName"`
	LocationID   string   `json:"locationID"`
	RoomName     string   `json:"roomName"`
	RoomID       string   `json:"roomID"`
	Month        int      `json:"month"`
	Year         int      `json:"year"`
	Monday       []Course `json:"monday"`
	Tuesday      []Course `json:"tuesday"`
	Wednesday    []Course `json:"wednesday"`
	Thursday     []Course `json:"thursday"`
	Friday       []Course `json:"friday"`
	Saturday     []Course `json:"saturday"`
	Sunday       []Course `json:"sunday"`
}

// CourseType Course Type
type CourseType struct {
	ID      string   `json:"id"`
	Name    string   `json:"name"`
	Classes []Course `json:"classes"`
}

// AppendClass helper func to append to Classes
func (ct *CourseType) AppendClass(item Course) []Course {
	ct.Classes = append(ct.Classes, item)
	return ct.Classes
}
