package db

import (
	"context"
	"fmt"
	"log"

	root "course_syndicate_api/pkg"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"

	"gopkg.in/mgo.v2"
)

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
	clientOptions := options.Client().ApplyURI(config.URL)

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
	session, err := mgo.Dial(config.URL)
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
