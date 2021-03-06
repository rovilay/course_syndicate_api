package server

import (
	"net/http"

	root "course_syndicate_api/pkg"
	"course_syndicate_api/pkg/controllers"
	"course_syndicate_api/pkg/middlewares"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewCourseRouter creates courses routes
func NewCourseRouter(r *mux.Router, c *mongo.Client, config *root.MongoConfig) {
	// create subrouter
	courseSubrouter := r.PathPrefix("/api/v1/courses").Subrouter()

	// attach handlers to subrouter
	courseController := controllers.NewCourseController(c, config)
	v := middlewares.NewValidator(c, config)

	courseSubrouter.HandleFunc(
		"/seed",
		middlewares.Authenticate(http.HandlerFunc(courseController.SeedCoursesData)),
	).Methods("GET")

	courseSubrouter.HandleFunc(
		"",
		middlewares.Authenticate(http.HandlerFunc(courseController.FetchCourses)),
	).Methods("GET")

	courseSubrouter.HandleFunc(
		"/{id}",
		middlewares.Authenticate(http.HandlerFunc(courseController.FetchSingleCourse)),
	).Methods("GET")

	courseSubrouter.HandleFunc(
		"/{id}/subscribe",
		middlewares.Authenticate(
			http.HandlerFunc(
				v.ValidateCourseSubscription(
					http.HandlerFunc(courseController.Subscribe),
				),
			),
		),
	).Methods("POST")

	courseSubrouter.HandleFunc(
		"/{id}/modules/{moduleId}",
		middlewares.Authenticate(
			http.HandlerFunc(
				v.ValidateCourseExist(
					http.HandlerFunc(courseController.FetchSingleModule),
				),
			)),
	).Methods("GET")
}
