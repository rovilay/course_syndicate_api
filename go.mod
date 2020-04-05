module course_syndicate_api

go 1.14

require (
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/gorilla/handlers v1.4.2
	github.com/gorilla/mux v1.7.4
	github.com/joho/godotenv v1.3.0
	github.com/stretchr/testify v1.5.1 // indirect
	go.mongodb.org/mongo-driver v1.3.1
	golang.org/x/crypto v0.0.0-20200323165209-0ec3e9974c59
	gopkg.in/mgo.v2 v2.0.0-20190816093944-a6b53ec6cb22
	gopkg.in/yaml.v2 v2.2.8 // indirect
)

// +heroku goVersion go1.14
// +heroku install -o bin/course-course_syndicate_api ./cmd/...
