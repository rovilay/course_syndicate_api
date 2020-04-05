package main

import (
	"log"

	root "course_syndicate_api/pkg"
	"course_syndicate_api/pkg/controllers"
	"course_syndicate_api/pkg/db"
	"course_syndicate_api/pkg/server"
	"course_syndicate_api/pkg/utils"
)

// App ...
type App struct {
	server  *server.Server
	client  *db.Client
	config  *root.Config
	sc      *controllers.ScheduleController
	workers int
}

// Initialize ...
func (a *App) initialize(workers int) {
	a.config = &root.Config{
		MongoConfig: &root.MongoConfig{
			URL:    utils.EnvOrDefaultString("DB_URL", ""),
			DBName: utils.EnvOrDefaultString("DB_NAME", "course_syndicate"),
		},
		ServerConfig: &root.ServerConfig{
			Port: utils.EnvOrDefaultString("PORT", "4444"),
		},
	}

	var err error
	a.client, err = db.NewClient(a.config.MongoConfig)

	if err != nil {
		log.Fatalln("[ERROR: CREATE MONGO CLIENT]", err)
	}

	a.workers = workers
	a.sc = controllers.NewScheduleController(a.client.Copy(), a.config.MongoConfig)
	a.server = server.InitServer(a.config, a.client)
}

// Start ...
func (a *App) start() {
	defer a.client.Close()
	a.server.Start()
}
