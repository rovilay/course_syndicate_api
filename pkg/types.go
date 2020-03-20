package root

import "time"

// User ...
type User struct {
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Email     string `json:"email"`
	Password  string `json:"password"`
}

// Course ...
type Course struct {
	Title           string    `json:"title"`
	NumberOfModules int       `json:"numberOfModules"`
	CreatedAt       time.Time `json:"createdAt"`
}

// CourseModule ...
type CourseModule struct {
	Title string `json:"title"`
}
