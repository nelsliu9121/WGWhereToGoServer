package tasks

// Course Type
type Course struct {
	Subject struct {
		ID    string `json:"id"`
		Name  string `json:"name"`
		Alias string `json:"alias"`
	} `json:"subject"`
	Teacher   string `json:"remark"`
	Weekday   string `json:"day_week"`
	StartTime string `json:"start_at"`
	EndTime   string `json:"end_at"`
}
