package root

// MongoConfig ...
type MongoConfig struct {
	Url    string
	DBName string
}

// ServerConfig ...
type ServerConfig struct {
	Port string
}

// Config ...
type Config struct {
	MongoConfig  *MongoConfig
	ServerConfig *ServerConfig
}
