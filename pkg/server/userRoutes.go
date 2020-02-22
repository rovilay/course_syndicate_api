package server

import (
	"github.com/gorilla/mux"
	root "github.com/rovilay/course_syndicate_api/pkg"
	"github.com/rovilay/course_syndicate_api/pkg/controllers"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserRouter creates users routes
func NewUserRouter(r *mux.Router, c *mongo.Client, config *root.MongoConfig) {
	// create subrouter
	userSubrouter := r.PathPrefix("/api/v1/users").Subrouter()

	// attach handlers to subrouter
	userController := controllers.NewUserController(c, config)
	userSubrouter.HandleFunc("", userController.CreateUserHandler).Methods("POST")
	userSubrouter.HandleFunc("/", userController.CreateUserHandler).Methods("POST")
}
