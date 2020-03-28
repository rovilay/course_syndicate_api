package middlewares

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"regexp"

	root "course_syndicate_api/pkg"
	"course_syndicate_api/pkg/db"
	"course_syndicate_api/pkg/utils"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewValidator ...
func NewValidator(c *mongo.Client, config *root.MongoConfig) *Validator {
	us := db.NewUserService(c, config)

	return &Validator{
		UserService: us,
		Errors:      make(map[string]string),
	}
}

// ValidateUserRegister ...
func (v *Validator) ValidateUserRegister(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, r *http.Request) {
		v.Errors = make(map[string]string)
		e := &utils.ErrorWithStatusCode{}
		var user root.User

		if r.Body == nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("no request body")

			utils.ErrorHandler(e, res)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&user)

		if err != nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("invalid request payload: all values must be in string format")

			utils.ErrorHandler(e, res)
			return
		}

		if user.FirstName == "" {
			v.Errors["firstname"] = "firstname is required"
		}

		if user.LastName == "" {
			v.Errors["lastname"] = "lastname is required"
		}

		if user.Email == "" {
			v.Errors["email"] = "email is required"
		} else {
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			if !re.MatchString(user.Email) {
				v.Errors["email"] = "email is invalid"
			}
		}

		if user.Password == "" {
			v.Errors["password"] = "password is required"
		} else if len(user.Password) < 7 {
			v.Errors["password"] = "password must be 7 or more"
		}

		if len(v.Errors) == 0 {
			// check if user already exist
			col := v.UserService.Collection
			var foundUser *db.UserModel

			col.FindOne(context.Background(), bson.M{"email": user.Email}).Decode(&foundUser)

			if foundUser != nil {
				e.StatusCode = http.StatusBadRequest
				e.ErrorMessage = errors.New("user already exist")

				utils.ErrorHandler(e, res)
				return
			}

			ctx := r.Context()
			ctx = context.WithValue(ctx, utils.ContextKey("user"), user)

			r = r.WithContext(ctx)
			next.ServeHTTP(res, r)
			return
		}

		e.StatusCode = http.StatusBadRequest
		e.ErrorMessage = errors.New("invalid request payload")
		e.Errors = v.Errors

		utils.ErrorHandler(e, res)
	}
}

// ValidateUserLogin ...
func (v *Validator) ValidateUserLogin(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, r *http.Request) {
		v.Errors = make(map[string]string)
		e := &utils.ErrorWithStatusCode{}
		var u root.User

		if r.Body == nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("no request body")

			utils.ErrorHandler(e, res)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&u)
		if err != nil {
			fmt.Println("[ERROR: ValidateUserLogin]", err)

			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("invalid request payload: all values must be in string format")

			utils.ErrorHandler(e, res)
			return
		}

		if u.Email == "" {
			v.Errors["email"] = "email is required"
		} else {
			re := regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

			if !re.MatchString(u.Email) {
				v.Errors["email"] = "email is invalid"
			}
		}

		if u.Password == "" {
			v.Errors["password"] = "password is required"
		}

		if len(v.Errors) == 0 {
			ctx := r.Context()
			ctx = context.WithValue(ctx, utils.ContextKey("user"), u)

			r = r.WithContext(ctx)
			next.ServeHTTP(res, r)
			return
		}

		e.StatusCode = http.StatusBadRequest
		e.ErrorMessage = errors.New("invalid request payload")
		e.Errors = v.Errors

		utils.ErrorHandler(e, res)
	}
}

// ValidateUserExist ...
func (v *Validator) ValidateUserExist(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, r *http.Request) {
		e := &utils.ErrorWithStatusCode{}

		ctx := r.Context()
		u := ctx.Value(utils.ContextKey("claims")).(*utils.JWTClaims)

		col := v.UserService.Collection
		var foundUser *db.UserModel

		objID, _ := primitive.ObjectIDFromHex(u.ID)
		err := col.FindOne(context.Background(), bson.M{"_id": objID}).Decode(&foundUser)

		if err != nil || foundUser == nil {
			fmt.Println("[ERROR: ValidateUserExist]: ", err)

			e.StatusCode = http.StatusUnauthorized
			e.ErrorMessage = errors.New("user does not exist")

			utils.ErrorHandler(e, res)
			return
		}

		ctx = context.WithValue(ctx, utils.ContextKey("authUser"), foundUser)
		r = r.WithContext(ctx)
		next.ServeHTTP(res, r)
		return
	}
}
