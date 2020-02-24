package server

import (
	"net/http"

	"github.com/gorilla/mux"
	root "github.com/rovilay/course_syndicate_api/pkg"
	"github.com/rovilay/course_syndicate_api/pkg/controllers"
	middlewares "github.com/rovilay/course_syndicate_api/pkg/middlewares"
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

	userSubrouter.HandleFunc(
		"/auth-test",
		middlewares.Authenticate(
			http.HandlerFunc(
				v.ValidateUserExist(http.HandlerFunc(userController.DummyController)),
			),
		),
	).Methods("GET")
}
