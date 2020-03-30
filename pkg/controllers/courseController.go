package controllers

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

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
	courseSubscriptionService := db.NewCourseSubService(c, config)
	subscriptionScheduleService := db.NewService(c, config, "subscription_schedules")
	host := utils.EnvOrDefaultString("SMTP_HOST", "smtp.gmail.com")
	port := utils.EnvOrDefaultString("SMTP_PORT", "587")
	smptService := &utils.SMTPServer{Host: host, Port: port}

	return &CourseController{
		courseService,
		courseModuleService,
		courseSubscriptionService,
		subscriptionScheduleService,
		smptService,
	}
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

// FetchSingleCourse ....
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

	matchStage := bson.M{"$match": bson.M{"_id": courseID}}
	lookupStage := bson.M{"$lookup": bson.M{
		"from":         "course_modules",
		"localField":   "_id",
		"foreignField": "courseId",
		"as":           "modules",
	}}

	cur, err := col.Aggregate(ctx, []bson.M{matchStage, lookupStage})
	if err != nil {
		fmt.Println("[ERROR: FETCH_COURSES]: ", err)

		e.StatusCode = http.StatusInternalServerError
		e.ErrorMessage = errors.New("something went wrong")

		utils.ErrorHandler(e, res)
		return
	}

	var cwm []db.CourseWithModule
	if err = cur.All(ctx, &cwm); err != nil {
		fmt.Println("[ERROR: FETCH_COURSES]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	cur.Close(ctx)
	utils.JSONResponseHandler(res, http.StatusOK, &genericResponseWithData{"operation successful", &cwm[0]})
}

// Subscribe ...
func (cc *CourseController) Subscribe(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: errors.New("Something went wrong"),
	}

	ctx := r.Context()
	u := ctx.Value(utils.ContextKey("claims")).(utils.JWTClaims)
	c := ctx.Value(utils.ContextKey("verifiedCourse")).(db.CourseModel)
	cs := ctx.Value(utils.ContextKey("verifiedSchedule")).([]int64)

	uid, err := primitive.ObjectIDFromHex(u.ID)
	if err != nil {
		fmt.Println("[ERROR: COURSE_SUBSCRIPTION_HANDLER]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	newSubscription := db.CreateCourseSubscriptionModel(uid, c.ID, cs)

	// create subscription
	col := cc.courseSubscriptionService.Collection
	s, err := col.InsertOne(context.Background(), newSubscription)
	if err != nil {
		fmt.Println("[ERROR: COURSE_SUBSCRIPTION_HANDLER]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	schedules := []interface{}{}
	for _, schedule := range cs {
		ns := db.CreateSubscriptionScheduleModel(newSubscription.ID, schedule)
		schedules = append(schedules, *ns)
	}

	// create schedules
	ssCol := cc.subscriptionScheduleService.Collection
	_, err = ssCol.InsertMany(context.Background(), schedules)
	if err != nil {
		fmt.Println("[ERROR: COURSE_SUBSCRIPTION_HANDLER]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	result := &courseSubscription{
		ID:               s.InsertedID.(primitive.ObjectID).Hex(),
		UserID:           u.ID,
		CourseID:         c.ID.Hex(),
		ModulesCompleted: newSubscription.ModulesCompleted,
		Schedule:         cs,
		CreatedAt:        c.CreatedAt,
	}

	utils.JSONResponseHandler(res, http.StatusOK, &genericResponseWithData{"operation successful", result})
}

// FetchSubscriptions ...
func (cc *CourseController) FetchSubscriptions(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: errors.New("fetch failed"),
	}

	const oneMillisecInNanosec = 1e6
	const oneSecInMilliSec = 1000
	now := time.Now().UnixNano() / oneMillisecInNanosec
	sendMailIntervalMins, err := strconv.Atoi(utils.EnvOrDefaultString("SEND_MAIL_INTERVAL_MINS", "15"))
	if err != nil {
		sendMailIntervalMins = 15
	}
	interval := sendMailIntervalMins * 60 * oneSecInMilliSec

	ctx := context.Background()
	col := cc.subscriptionScheduleService.Collection

	addFieldStage := bson.M{"$addFields": bson.M{
		"timeDiff": bson.M{"$subtract": []interface{}{now, "$schedule"}},
		"isMailReady": bson.M{"$and": []bson.M{
			bson.M{"$gte": []interface{}{"$timeDiff", 0}},
			bson.M{"$lte": []interface{}{"$timeDiff", interval}},
		}},
	}}
	matchStage := bson.M{"$match": bson.M{"isMailReady": false, "completed": false}}
	lookupStage1 := bson.M{"$lookup": bson.M{
		"from":         "course_subscriptions",
		"localField":   "subscriptionId",
		"foreignField": "_id",
		"as":           "subscription",
	}}
	unwindStage1 := bson.M{"$unwind": bson.M{
		"path":                       "$subscription",
		"preserveNullAndEmptyArrays": true,
	}}
	lookupStage2 := bson.M{"$lookup": bson.M{
		"from":         "users",
		"localField":   "subscription.userId",
		"foreignField": "_id",
		"as":           "user",
	}}
	unwindStage2 := bson.M{"$unwind": bson.M{
		"path":                       "$user",
		"preserveNullAndEmptyArrays": true,
	}}
	lookupStage3 := bson.M{"$lookup": bson.M{
		"from":         "courses",
		"localField":   "subscription.courseId",
		"foreignField": "_id",
		"as":           "course",
	}}
	unwindStage3 := bson.M{"$unwind": bson.M{
		"path":                       "$course",
		"preserveNullAndEmptyArrays": true,
	}}
	lookupStage4 := bson.M{"$lookup": bson.M{
		"from":         "course_modules",
		"localField":   "subscription.courseId",
		"foreignField": "courseId",
		"as":           "courseModules",
	}}
	project := bson.M{"$project": bson.M{
		"subscriptionId":        0,
		"subscription.userId":   0,
		"subscription.courseId": 0,
		"user._id":              0,
		"user.password":         0,
		"user.createdAt":        0,
	}}

	cur, err := col.Aggregate(ctx, []bson.M{
		addFieldStage,
		matchStage,
		lookupStage1,
		unwindStage1,
		lookupStage2,
		unwindStage2,
		lookupStage3,
		unwindStage3,
		lookupStage4,
		project,
	})
	if err != nil {
		fmt.Println("[ERROR: FETCH_SUBS]: ", err)

		e.StatusCode = http.StatusInternalServerError
		e.ErrorMessage = errors.New("something went wrong")

		utils.ErrorHandler(e, res)
		return
	}

	var subs []fetchSchedulesResult
	if err = cur.All(ctx, &subs); err != nil {
		fmt.Println("[ERROR: FETCH_SUBS]: ", err)

		utils.ErrorHandler(e, res)
		return
	}

	cur.Close(ctx)
	utils.JSONResponseHandler(res, http.StatusOK, &genericResponseWithData{"operation successful", &subs})
}

// MailCourse ...
func (cc *CourseController) MailCourse(emails []string, messages [][]byte) {
	from := utils.EnvOrDefaultString("SMTP_EMAIL", "")
	password := utils.EnvOrDefaultString("SMTP_EMAIL_PASSWORD", "")

	for i, email := range emails {
		err := cc.smtpService.SendMail(from, password, []string{email}, messages[i])
		if err != nil {
			panic(err)
		}
	}
}
