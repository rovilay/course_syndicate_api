package server

import (
	"github.com/gorilla/mux"
	root "github.com/rovilay/course_syndicate_api/pkg"
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
