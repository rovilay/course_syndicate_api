package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	root "github.com/rovilay/course_syndicate_api/pkg"
	"github.com/rovilay/course_syndicate_api/pkg/db"
	"github.com/rovilay/course_syndicate_api/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// InitServer ...
func InitServer(config *root.Config, c *db.Client) *Server {
	s := Server{
		route:  mux.NewRouter(),
		config: config.ServerConfig,
	}

	// Initialize routers
	s.InitRouters(c.Copy(), config.MongoConfig)

	return &s
}

// InitRouters ...
func (s *Server) InitRouters(client *mongo.Client, config *root.MongoConfig) {
	s.route.HandleFunc("/", welcomeHandler)
	NewUserRouter(s.route, client, config)
}

func welcomeHandler(res http.ResponseWriter, r *http.Request) {
	res.Header().Set("Content-Type", "application/json")

	message := WelcomeResponse{"welcome to course_syndicate"}
	m, err := json.Marshal(message)

	if err != nil {
		e := &utils.ErrorWithStatusCode{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: err,
		}

		utils.ErrorHandler(e, res)
	}

	res.WriteHeader(http.StatusOK)
	res.Write(m)
}

// Start ...
func (s *Server) Start() {
	fmt.Println("Listening on port " + s.config.Port)
	loggedRouter := handlers.LoggingHandler(os.Stdout, s.route)
	err := http.ListenAndServe(s.config.Port, loggedRouter)
	if err != nil {
		log.Fatalln("httpListenAndServe: ", err)
	}
}
