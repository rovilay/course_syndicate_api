package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"

	root "course_syndicate_api/pkg"
	"course_syndicate_api/pkg/db"
	"course_syndicate_api/pkg/utils"

	"github.com/gorilla/mux"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

// NewCourseController ...
func NewCourseController(c *mongo.Client, config *root.MongoConfig) *CourseController {
	courseService := db.NewService(c, config, "courses")
	courseModuleService := db.NewService(c, config, "course_modules")

	return &CourseController{courseService, courseModuleService}
}

// SeedCoursesData ...
func (cc *CourseController) SeedCoursesData(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: errors.New("course seed failed"),
	}

	col := cc.courseService.Collection
	cmCol := cc.courseModuleService.Collection

	courses := []interface{}{}
	modulesA := []interface{}{}
	modulesB := []interface{}{}

	for _, c := range db.Courses {
		nc := db.CreateCourseModel(c)
		courses = append(courses, *nc)
	}

	results, err := col.InsertMany(context.Background(), courses)
	if err != nil {
		fmt.Println("[ERROR: SEED_COURSES]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	for _, m := range db.CourseAModules {
		ncm := db.CreateCourseModuleModel(results.InsertedIDs[0].(primitive.ObjectID), m)
		modulesA = append(modulesA, *ncm)
	}

	for _, m := range db.CourseBModules {
		ncm := db.CreateCourseModuleModel(results.InsertedIDs[1].(primitive.ObjectID), m)
		modulesB = append(modulesB, *ncm)
	}

	_, err = cmCol.InsertMany(context.Background(), append(modulesA, modulesB...))
	if err != nil {
		fmt.Println("[ERROR: SEED_COURSES]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	utils.JSONResponseHandler(res, http.StatusCreated, genericResponse{"course seed successful"})
}

// FetchCourses ...
func (cc *CourseController) FetchCourses(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: errors.New("fetch failed"),
	}

	col := cc.courseService.Collection

	findOptions := options.Find()
	findOptions.SetLimit(10)

	var results []*db.CourseModel
	ctx := context.Background()
	cur, err := col.Find(ctx, bson.M{}, findOptions)

	if err != nil {
		fmt.Println("[ERROR: FETCH_COURSES]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var c db.CourseModel
		err := cur.Decode(&c)
		if err != nil {
			utils.ErrorHandler(e, res)
			log.Fatalln("[ERROR: FETCH_COURSES]: ", err)
		}

		results = append(results, &c)
	}

	if err := cur.Err(); err != nil {
		fmt.Println("[ERROR: FETCH_COURSES]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	// Close the cursor once finished
	cur.Close(ctx)
	utils.JSONResponseHandler(res, http.StatusOK, &genericResponseWithData{"operation successful", results})
}

// FetchSingleCourse ...
func (cc *CourseController) FetchSingleCourse(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: errors.New("fetch failed"),
	}

	params := mux.Vars(r)
	courseID, err := primitive.ObjectIDFromHex(params["id"])

	if err != nil {
		fmt.Println("[ERROR: FETCH_COURSES]: ", err)

		e.StatusCode = http.StatusBadRequest
		e.ErrorMessage = errors.New("invalid course id")

		utils.ErrorHandler(e, res)
		return
	}

	ctx := context.Background()
	col := cc.courseService.Collection
	var c *db.CourseModel

	err = col.FindOne(ctx, bson.M{"_id": courseID}).Decode(&c)
	if err != nil {
		fmt.Println("[ERROR: FETCH_COURSES]: ", err)

		e.StatusCode = http.StatusInternalServerError
		e.ErrorMessage = errors.New("something went wrong")

		if c == nil {
			e.StatusCode = http.StatusNotFound
			e.ErrorMessage = errors.New("course not found")
		}

		utils.ErrorHandler(e, res)
		return
	}

	// find course modules
	mCol := cc.courseModuleService.Collection
	var modules []*db.CourseModuleModel

	cur, err := mCol.Find(ctx, bson.M{"courseId": c.ID})
	if err := cur.Err(); err != nil {
		fmt.Println("[ERROR: FETCH_COURSES]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	for cur.Next(ctx) {
		// create a value into which the single document can be decoded
		var cm db.CourseModuleModel
		err := cur.Decode(&cm)
		if err != nil {
			e.StatusCode = http.StatusInternalServerError
			e.ErrorMessage = errors.New("something went wrong")

			utils.ErrorHandler(e, res)
			log.Fatalln("[ERROR: FETCH_COURSES]: ", err)
		}

		modules = append(modules, &cm)
	}

	result := &courseWithModule{
		ID:              c.ID,
		Title:           c.Title,
		NumberOfModules: c.NumberOfModules,
		Modules:         modules,
		CreatedAt:       c.CreatedAt,
	}

	// Close the cursor once finished
	cur.Close(ctx)
	utils.JSONResponseHandler(res, http.StatusOK, &genericResponseWithData{"operation successful", &result})
}
