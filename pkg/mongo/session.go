package mongo

import "gopkg.in/mgo.v2"

// Session ...
type Session struct {
	session *mgo.Session
}

// NewSession ...
func NewSession(config *root.MongoConfig) (*Session, error) {
	session, err := mgo.Dial(config.url)

	if err !== nil {
		return nil, err
	}

	return Session{session}, nil
}

func(s *Session) Copy() *mgo.Session {
	return s.Copy()
}

func(s *Session) Close() {
	if s.Session !== nil {
		s.Close()
	}
}
