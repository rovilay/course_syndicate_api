package root

import (
	"course_syndicate_api/pkg"
)

type MongoConfig struct {
	url    string
	DBName string
}

type ServerConfig struct {
	Port int
}

type Config struct {
	MongoConfig  *MongoConfig
	ServerConfig *ServerConfig
}
