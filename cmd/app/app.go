package main

import (
	"log"

	root "github.com/rovilay/course_syndicate_api/pkg"
	"github.com/rovilay/course_syndicate_api/pkg/db"
	"github.com/rovilay/course_syndicate_api/pkg/server"
	"github.com/rovilay/course_syndicate_api/pkg/utils"
)

// App ...
type App struct {
	server *server.Server
	client *db.Client
	config *root.Config
}

// Initialize ...
func (a *App) Initialize() {
	a.config = &root.Config{
		MongoConfig: &root.MongoConfig{
			Url:    utils.EnvOrDefaultString("DB_URL", ""),
			DBName: utils.EnvOrDefaultString("DB_NAME", "course_syndicate"),
		},
		ServerConfig: &root.ServerConfig{
			Port: utils.EnvOrDefaultString("PORT", ":4444"),
		},
	}

	var err error
	a.client, err = db.NewClient(a.config.MongoConfig)

	if err != nil {
		log.Fatalln("[ERROR: CREATE MONGO CLIENT]", err)
	}

	a.server = server.InitServer(a.config, a.client)
}

// Start ...
func (a *App) Start() {
	defer a.client.Close()
	a.server.Start()
}