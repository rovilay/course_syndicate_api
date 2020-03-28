package db

import (
	root "course_syndicate_api/pkg"
)

// CourseA ...
var CourseA = &root.Course{
	Title:           "Understanding social media influence on politics and governance",
	NumberOfModules: 3,
}

// CourseB ...
var CourseB = &root.Course{
	Title:           "Go for beginners",
	NumberOfModules: 2,
}

// CourseAModule1 ...
var CourseAModule1 = &root.CourseModule{
	Title: "Module 1: Introduction to Social Media",
}

// CourseAModule2 ...
var CourseAModule2 = &root.CourseModule{
	Title: "Module 2: A Breif History of Social Media",
}

// CourseAModule3 ...
var CourseAModule3 = &root.CourseModule{
	Title: "Module 3: Social Media, Politics and Goverment",
}

// CourseAModule4 ...
var CourseAModule4 = &root.CourseModule{
	Title: "Module 1: Introduction to this Go Course",
}

// CourseAModule5 ...
var CourseAModule5 = &root.CourseModule{
	Title: "Module 2: Datatypes in Go",
}

// Courses ...
var Courses = []*root.Course{CourseA, CourseB}

// CourseAModules ...
var CourseAModules = []*root.CourseModule{CourseAModule1, CourseAModule2, CourseAModule3}

// CourseBModules ...
var CourseBModules = []*root.CourseModule{CourseAModule4, CourseAModule5}
