package db

import (
	"time"

	root "github.com/rovilay/course_syndicate_api/pkg"
)

// CourseA ...
var CourseA = &root.Course{
	Title:           "Understanding social media influence on politics and governance",
	NumberOfModules: 5,
	CreatedAt:       time.Now(),
}

// CourseB ...
var CourseB = &root.Course{
	Title:           "Go for beginners",
	NumberOfModules: 5,
	CreatedAt:       time.Now(),
}

// CourseA Modules
var CourseAModule1 = &root.CourseModule{
	Title: "Module 1: Introduction to Course"
}
var CourseAModule1 = &root.CourseModule{
	Title: "Module 1: Introduction to Course"
}
var CourseAModule1 = &root.CourseModule{
	Title: "Module 1: Introduction to Course"
}
var CourseAModule1 = &root.CourseModule{
	Title: "Module 1: Introduction to Course"
}
var CourseAModule1 = &root.CourseModule{
	Title: "Module 1: Introduction to Course"
}

// Courses ...
var Courses = []*root.Course{CourseA, CourseA}

// CourseModules ...
var CourseModules = []*root.CourseModule{}
