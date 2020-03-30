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

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewValidator ...
func NewValidator(c *mongo.Client, config *root.MongoConfig) *Validator {
	us := db.NewUserService(c, config)
	cs := db.NewService(c, config, "courses")
	css := db.NewCourseSubService(c, config)

	return &Validator{
		userService:               us,
		courseService:             cs,
		courseSubscriptionService: css,
		Errors:                    make(map[string]string),
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
			col := v.userService.Collection
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

		col := v.userService.Collection
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

// CheckCourseExist ...
func (v *Validator) CheckCourseExist(courseID primitive.ObjectID) (c *db.CourseModel, err error) {
	ctx := context.Background()
	col := v.courseService.Collection

	err = col.FindOne(ctx, bson.M{"_id": courseID}).Decode(&c)
	if err != nil {
		fmt.Println("[ERROR: ValidateCourseExist]: ", err)
		err = errors.New("something went wrong")

		if c == nil {
			err = errors.New("course not found")
		}

		return
	}

	return
}

// CheckSubscriptionExist ...
func (v *Validator) CheckSubscriptionExist(userID, courseID primitive.ObjectID) (csm *db.CourseSubscriptionModel, err error) {
	ctx := context.Background()
	col := v.courseSubscriptionService.Collection

	err = col.FindOne(ctx, bson.M{"userId": userID, "courseId": courseID}).Decode(&csm)
	if err != nil {
		fmt.Println("[ERROR: ValidateCourseExist]: ", err)
		err = errors.New("something went wrong")

		return
	}

	return
}

// ValidateSchedule ...
func (v *Validator) ValidateSchedule(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, r *http.Request) {
		v.Errors = make(map[string]string)
		e := &utils.ErrorWithStatusCode{}
		var cs root.CourseShedulePayload

		if r.Body == nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("no request body")

			utils.ErrorHandler(e, res)
			return
		}

		err := json.NewDecoder(r.Body).Decode(&cs)
		if err != nil {
			fmt.Println("[ERROR: ValidateUserLogin]", err)

			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("invalid request payload: all values must be in string format")

			utils.ErrorHandler(e, res)
			return
		}

		const dateFormat = "2006-01-02T15:04:05"
		ts, err := utils.Schedular(cs.Schedule, dateFormat, 3)

		if err != nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = err

			utils.ErrorHandler(e, res)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.ContextKey("verifiedSchedule"), ts)
		r = r.WithContext(ctx)
		next.ServeHTTP(res, r)
		return
	}
}

// ValidateCourseSubscription ...
func (v *Validator) ValidateCourseSubscription(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, r *http.Request) {
		v.Errors = make(map[string]string)
		e := &utils.ErrorWithStatusCode{}

		// validate request body
		if r.Body == nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("no request body")

			utils.ErrorHandler(e, res)
			return
		}

		// check if course exist
		params := mux.Vars(r)
		courseID, err := primitive.ObjectIDFromHex(params["id"])
		c, err := v.CheckCourseExist(courseID)

		if err != nil {
			fmt.Println("[ERROR: ValidateCourseSubscription]", err)

			e.StatusCode = http.StatusNotFound
			e.ErrorMessage = err

			utils.ErrorHandler(e, res)
			return
		}

		// check if subscription exist
		ctx := r.Context()
		u := ctx.Value(utils.ContextKey("claims")).(utils.JWTClaims)
		uid, _ := primitive.ObjectIDFromHex(u.ID)

		csm, _ := v.CheckSubscriptionExist(uid, courseID)

		if csm != nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("you have already subscribed to this course")

			utils.ErrorHandler(e, res)
			return
		}

		// Validate Schedule payload and create Schedule
		var cs root.CourseShedulePayload
		err = json.NewDecoder(r.Body).Decode(&cs)
		if err != nil {
			fmt.Println("[ERROR: ValidateUserLogin]", err)

			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("invalid request payload: all values must be in string format")

			utils.ErrorHandler(e, res)
			return
		}

		const dateFormat = "2006-01-02T15:04:05"
		ts, err := utils.Schedular(cs.Schedule, dateFormat, c.NumberOfModules)

		if err != nil {
			e.StatusCode = http.StatusBadRequest
			e.ErrorMessage = errors.New("invalid payload")

			if err.Error() == "invalid string" {
				v.Errors["schedule"] = "invalid schedule. valid schedule formats are: `every <number> days|weeks|months` or 'YYYY-MM-DDTHH:mm:ss,YYYY-MM-DDTHH:mm:ss'"
				e.Errors = v.Errors
			} else if err.Error() == "time has expired" {
				v.Errors["schedule"] = "invalid schedule. datetime must be greater than now"
			}

			e.Errors = v.Errors

			utils.ErrorHandler(e, res)
			return
		}

		ctx = context.WithValue(ctx, utils.ContextKey("verifiedCourse"), *c)
		ctx = context.WithValue(ctx, utils.ContextKey("verifiedSchedule"), ts)

		r = r.WithContext(ctx)
		next.ServeHTTP(res, r)
		return
	}
}
