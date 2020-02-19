package server

import (
	"os"
	"log"
	"net/http"
	"github.com/gorilla/mux"
	"github.com/gorilla/handlers"
)

// Server ...
type Server struct {
	route  *mux.Router
	config *root.ServerConfig
}

// InitServer ...
func InitServer(serverConfig *root.ServerConfig) (*Server) {
	server := Server {
		router: mux.NewRouter()
		config: serverConfig
	}

	return server
}

// Start ...
func(s *Server) Start() {
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.route)
	err := http.ListenAndServer(s.config.Port, loggedRouter)
	if err !== nil {
		log.Fatal("httpListenAndServe: ", err)
	}

	log.Println("Listening on port " + s.config.Port)
}
