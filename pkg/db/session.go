package db

import (
	"context"
	"fmt"
	"log"

	root "github.com/rovilay/course_syndicate_api/pkg"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/mgo.v2"
)

// const (
// 	AuthDatabase   = "course_syndicate"
// 	AuthUserName   = "rovilay"
// 	AuthPassword   = "qwertyUp1."
// 	ReplicaSetName = "ireporter-cluster-shard-0"
// )

// Session ...
type Session struct {
	session *mgo.Session
}

// Client ...
type Client struct {
	client *mongo.Client
}

// NewClient ...
func NewClient(config *root.MongoConfig) (*Client, error) {
	clientOptions := options.Client().ApplyURI(config.Url)

	// Connect to MongoDB
	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		return nil, err
	}

	fmt.Println("Connected to MongoDB!")
	return &Client{client}, err
}

// NewSession ...
func NewSession(config *root.MongoConfig) (*Session, error) {
	// MongoDBHosts := []string{
	// 	"ireporter-cluster-shard-00-00-y4nzl.mongodb.net:27017",
	// 	"ireporter-cluster-shard-00-01-y4nzl.mongodb.net:27017",
	// 	"ireporter-cluster-shard-00-02-y4nzl.mongodb.net:27017",
	// }

	// mongoDBDialInfo := &mgo.DialInfo{
	// 	Addrs:          MongoDBHosts,
	// 	Timeout:        60 * time.Second,
	// 	Database:       AuthDatabase,
	// 	Username:       AuthUserName,
	// 	Password:       AuthPassword,
	// 	ReplicaSetName: ReplicaSetName,
	// }

	// mongoDBDialInfo.DialServer = func(addr *mgo.ServerAddr) (net.Conn, error) {
	// 	conn, err := tls.Dial("tcp", addr.String(), tlsConfig)
	// 	return conn, err
	// }

	// session, err := mgo.DialWithInfo(mongoDBDialInfo)

	session, err := mgo.Dial(config.Url)
	if err != nil {
		return nil, err
	}

	return &Session{session}, nil
}

// Copy ...
func (s *Session) Copy() *mgo.Session {
	return s.session.Copy()
}

// Close ...
func (s *Session) Close() {
	if s.session != nil {
		fmt.Println("CLOSING MONGO SESSION")
		s.session.Close()
	}
}

// Close ...
func (c *Client) Close() {
	if c.client != nil {
		err := c.client.Disconnect(context.Background())
		if err != nil {
			log.Fatal(err)
		}
		fmt.Println("Connection to MongoDB closed.")
	}
}

// Copy ...
func (c *Client) Copy() *mongo.Client {
	return c.client
}
