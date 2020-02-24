package controllers

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"

	root "github.com/rovilay/course_syndicate_api/pkg"
	"github.com/rovilay/course_syndicate_api/pkg/db"
	"github.com/rovilay/course_syndicate_api/pkg/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserController ...
func NewUserController(c *mongo.Client, config *root.MongoConfig) *UserController {
	userService := db.NewUserService(c, config)

	return &UserController{userService}
}

// CreateUserHandler ...
func (uc *UserController) CreateUserHandler(res http.ResponseWriter, r *http.Request) {
	var user root.User

	err := json.NewDecoder(r.Body).Decode(&user)
	newUser, err := db.CreateUserModel(&user)

	if err != nil {
		fmt.Println("[ERROR: CREATE_USER_HANDLER]: ", err)
		e := &utils.ErrorWithStatusCode{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: errors.New("Something went wrong"),
		}

		utils.ErrorHandler(e, res)
	}

	col := uc.userService.Collection
	result, err := col.InsertOne(context.Background(), newUser)

	if err != nil {
		e := &utils.ErrorWithStatusCode{
			StatusCode:   http.StatusInternalServerError,
			ErrorMessage: errors.New("Something went wrong"),
		}

		utils.ErrorHandler(e, res)
	}

	utils.JSONResponseHandler(res, http.StatusOK, result)
}
