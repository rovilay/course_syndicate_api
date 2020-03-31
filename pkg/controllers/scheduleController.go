package controllers

import (
	"context"
	root "course_syndicate_api/pkg"
	"course_syndicate_api/pkg/db"
	"course_syndicate_api/pkg/utils"
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"sync"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

// NewScheduleController ...
func NewScheduleController(c *mongo.Client, config *root.MongoConfig) *ScheduleController {
	courseSubscriptionService := db.NewCourseSubService(c, config)
	subscriptionScheduleService := db.NewService(c, config, "subscription_schedules")
	host := utils.EnvOrDefaultString("SMTP_HOST", "smtp.gmail.com")
	port := utils.EnvOrDefaultString("SMTP_PORT", "587")
	smptService := &utils.SMTPServer{Host: host, Port: port}

	return &ScheduleController{
		courseSubscriptionService,
		subscriptionScheduleService,
		smptService,
	}
}

// FetchSchedules ...
func (sc *ScheduleController) FetchSchedules() (subs []fetchSchedulesResult, err error) {
	const oneMillisecInNanosec = 1e6
	const oneSecInMilliSec = 1000
	now := time.Now().UnixNano() / oneMillisecInNanosec
	sendMailIntervalMins, err := strconv.Atoi(utils.EnvOrDefaultString("SEND_MAIL_INTERVAL_MINS", "15"))
	if err != nil {
		sendMailIntervalMins = 15
	}
	interval := sendMailIntervalMins * 60 * oneSecInMilliSec

	ctx := context.Background()
	col := sc.subscriptionScheduleService.Collection

	addFieldStage := bson.M{"$addFields": bson.M{
		"timeDiff": bson.M{"$subtract": []interface{}{now, "$schedule"}},
		"isMailReady": bson.M{"$and": []bson.M{
			bson.M{"$gte": []interface{}{"$timeDiff", 0}},
			bson.M{"$lte": []interface{}{"$timeDiff", interval}},
		}},
	}}
	matchStage := bson.M{"$match": bson.M{"isMailReady": true, "completed": false}}
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

		return
	}

	if err = cur.All(ctx, &subs); err != nil {
		fmt.Println("[ERROR: FETCH_SUBS]: ", err)

		return
	}

	cur.Close(ctx)
	return
}

// reconcileSchedule updates the schedule completed status and
// subscription's modulesCompleted count
func (sc *ScheduleController) reconcileSchedule(sID, subID primitive.ObjectID) (err error) {
	ctx := context.Background()
	sCol := sc.subscriptionScheduleService.Collection
	subCol := sc.courseSubscriptionService.Collection

	// update schedule
	update := bson.M{"$set": bson.M{"completed": true}}
	_, err = sCol.UpdateOne(ctx, bson.M{"_id": sID}, update)
	if err != nil {
		return
	}

	// update subscription
	update = bson.M{"$inc": bson.M{"modulesCompleted": 1}}
	_, err = subCol.UpdateOne(ctx, bson.M{"_id": subID}, update)

	return
}

func (sc *ScheduleController) processSchedule(sender, password string, sch fetchSchedulesResult) {
	// because there is a possiblity that the courseModules can be reduced
	// then it gives a scenerio where the modulesCompleted can be greater or equals to the courseModules
	if sch.Subscription.ModulesCompleted >= len(sch.CourseModules) {
		return
	}

	module := sch.CourseModules[sch.Subscription.ModulesCompleted]
	ml := fmt.Sprintf("http://localhost:4444/api/v1/courses/%s/modules/%s", sch.Course.ID.Hex(), module.ID.Hex())
	sbj := "Course Syndicate"

	msg, err := utils.GenerateMailTemplate("template.html", &utils.MailTemplateData{
		Username:    fmt.Sprintf("%s %s", sch.User.FirstName, sch.User.LastName),
		CourseTitle: sch.Course.Title,
		ModuleTitle: module.Title,
		ModuleLink:  ml,
	})
	if err != nil {
		panic(err)
	}

	// send mail
	err = sc.smtpService.SendEmail(sender, password, sbj, msg, []string{sch.User.Email})
	if err != nil {
		panic(err)
	}

	// reconcile schedule
	err = sc.reconcileSchedule(sch.ID, sch.Subscription.ID)
	if err != nil {
		panic(err)
	}
}

func (sc *ScheduleController) psWorker(roc <-chan fetchSchedulesResult, wg *sync.WaitGroup, from, password string) {
	defer wg.Done()
	for sh := range roc {
		sc.processSchedule(from, password, sh)
	}
}

// SyncSchedules ...
func (sc *ScheduleController) SyncSchedules(res http.ResponseWriter, r *http.Request) {
	e := &utils.ErrorWithStatusCode{
		StatusCode:   http.StatusInternalServerError,
		ErrorMessage: errors.New("fetch failed"),
	}

	schdls, err := sc.FetchSchedules()
	if err != nil {
		utils.ErrorHandler(e, res)
		return
	}

	from := utils.EnvOrDefaultString("SMTP_EMAIL", "")
	password := utils.EnvOrDefaultString("SMTP_EMAIL_PASSWORD", "")

	if from == "" || password == "" {
		e.ErrorMessage = errors.New("SMTP_EMAIL and SMTP_EMAIL_PASSWORD not provided")

		utils.ErrorHandler(e, res)
		return
	}

	wg := sync.WaitGroup{}
	schdlChan := make(chan fetchSchedulesResult, len(schdls))
	wg.Add(2) // add total number of goroutines spawn

	go sc.psWorker(schdlChan, &wg, from, password)

	go func(woc chan<- fetchSchedulesResult) {
		defer wg.Done()

		for _, sh := range schdls {
			woc <- sh
		}

		close(woc)
	}(schdlChan)

	wg.Wait()

	utils.JSONResponseHandler(res, http.StatusOK, &genericResponse{"operation successful"})
	return
}
