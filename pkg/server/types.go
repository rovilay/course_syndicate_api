package server

import (
	root "course_syndicate_api/pkg"

	"github.com/gorilla/mux"
)

// Server ...
type Server struct {
	route  *mux.Router
	config *root.ServerConfig
}

// WelcomeResponse ...
type WelcomeResponse struct {
	Message string `json:"message" bson:"message"`
}
