package controllers

import (
	"context"
	"errors"
	"fmt"
	"net/http"

	root "course_syndicate_api/pkg"
	"course_syndicate_api/pkg/db"
	"course_syndicate_api/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewUserController ...
func NewUserController(c *mongo.Client, config *root.MongoConfig) *UserController {
	userService := db.NewUserService(c, config)

	return &UserController{userService}
}

// CreateUserHandler ...
func (uc *UserController) CreateUserHandler(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: errors.New("Something went wrong"),
	}

	ctx := r.Context()
	user := ctx.Value(utils.ContextKey("user")).(root.User)

	newUser, err := db.CreateUserModel(&user)

	if err != nil {
		fmt.Println("[ERROR: CREATE_USER_HANDLER]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	col := uc.userService.Collection
	result, err := col.InsertOne(context.Background(), newUser)

	if err != nil {
		fmt.Println("[ERROR: CREATE_USER_HANDLER]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	tp := &utils.TokenPayload{
		ID:    result.InsertedID.(primitive.ObjectID).Hex(),
		Email: user.Email,
	}

	token, err := utils.GenerateToken(tp)

	if err != nil {
		fmt.Println("[ERROR: CREATE_USER_HANDLER]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	utils.JSONResponseHandler(res, http.StatusCreated, authResponse{token})
}

// LoginUserHandler ...
func (uc *UserController) LoginUserHandler(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{}
	ctx := r.Context()
	u := ctx.Value(utils.ContextKey("user")).(root.User)
	var result *db.UserModel

	col := uc.userService.Collection
	err := col.FindOne(context.Background(), bson.M{"email": u.Email}).Decode(&result)

	if err != nil {
		fmt.Println("[ERROR: LoginUserHandler]: ", err)

		e.StatusCode = http.StatusInternalServerError
		e.ErrorMessage = errors.New("something went wrong")

		if result == nil {
			e.StatusCode = http.StatusNotFound
			e.ErrorMessage = errors.New("user not found")
		}

		utils.ErrorHandler(e, res)
		return
	}

	if !result.ComparePasswordHash(u.Password) {
		e.StatusCode = http.StatusBadRequest
		e.ErrorMessage = errors.New("wrong password")

		utils.ErrorHandler(e, res)
		return
	}

	tp := &utils.TokenPayload{
		ID:    result.ID.Hex(),
		Email: result.Email,
	}

	token, err := utils.GenerateToken(tp)

	if err != nil {
		fmt.Println("[ERROR: CREATE_USER_HANDLER]: ", err)
		e.StatusCode = http.StatusInternalServerError
		e.ErrorMessage = errors.New("something went wrong")

		utils.ErrorHandler(e, res)
		return
	}

	utils.JSONResponseHandler(res, http.StatusOK, authResponse{token})
}

// DummyController ...
func (uc *UserController) DummyController(res http.ResponseWriter, r *http.Request) {
	utils.JSONResponseHandler(res, http.StatusOK, &root.User{})
}
