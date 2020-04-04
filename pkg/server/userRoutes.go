package server

import (
	"net/http"

	root "course_syndicate_api/pkg"
	"course_syndicate_api/pkg/controllers"
	middlewares "course_syndicate_api/pkg/middlewares"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserRouter creates users routes
func NewUserRouter(r *mux.Router, c *mongo.Client, config *root.MongoConfig) {
	// create subrouter
	userSubrouter := r.PathPrefix("/api/v1").Subrouter()

	// attach handlers to subrouter
	userController := controllers.NewUserController(c, config)
	v := middlewares.NewValidator(c, config)

	userSubrouter.HandleFunc(
		"/register",
		v.ValidateUserRegister(http.HandlerFunc(userController.CreateUserHandler)),
	).Methods("POST")

	userSubrouter.HandleFunc(
		"/login",
		v.ValidateUserLogin(http.HandlerFunc(userController.LoginUserHandler)),
	).Methods("POST")
}
