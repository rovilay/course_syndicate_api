package utils

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"html/template"
	"log"
	"net/http"
	"net/smtp"
	"os"
	"regexp"
	"strconv"
	"strings"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/joho/godotenv"
)

// EnvOrDefaultString ...
func EnvOrDefaultString(envVar string, defaultValue string) string {
	// load .env file
	err := godotenv.Load(".env")

	if err != nil {
		log.Fatalf("Error loading .env file")
	}

	value := os.Getenv(envVar)

	if value == "" {
		return defaultValue
	}

	return value
}

// ErrorHandler handles https error responses
func ErrorHandler(err *ErrorWithStatusCode, res http.ResponseWriter) {
	// log.Fatal(err.ErrorMessage.Error())
	res.Header().Set("Content-Type", "application/json")

	e, _ := json.Marshal(ErrorResponse{
		Message: err.ErrorMessage.Error(),
		Errors:  err.Errors,
	})

	res.WriteHeader(err.StatusCode)
	res.Write(e)
}

// JSONResponseHandler handles http response in json
func JSONResponseHandler(res http.ResponseWriter, statusCode int, payload interface{}) {
	response, _ := json.Marshal(payload)

	res.Header().Set("Content-Type", "application/json")
	res.WriteHeader(statusCode)
	res.Write(response)
}

// GenerateToken ...
func GenerateToken(tp *TokenPayload) (string, error) {
	secret := EnvOrDefaultString("SECRET", "")
	et := EnvOrDefaultString("TOKEN_EXPIRATION_IN_HOURS", "24")
	i, err := strconv.Atoi(et)

	if err != nil {
		i = 24
	}

	expirationTime := time.Now().Add(time.Duration(i) * time.Hour)
	claims := &JWTClaims{
		ID:    tp.ID,
		Email: tp.Email,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: expirationTime.Unix(),
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	// Create the JWT string
	tokenString, err := token.SignedString([]byte(secret))

	return tokenString, err
}

// DecodeToken ...
func DecodeToken(token string) (*JWTClaims, error) {
	secret := EnvOrDefaultString("SECRET", "")
	claims := &JWTClaims{}

	tkn, err := jwt.ParseWithClaims(token, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(secret), nil
	})

	if err != nil {
		return claims, err
	}

	if !tkn.Valid {
		return claims, errors.New("Invalid token")
	}

	return claims, err
}

// Schedular ...
func Schedular(s, dateFormat string, numberOfSchedule int) (ts []int64, err error) {
	re := regexp.MustCompile("^(every)[ ](0?[1-9]([0-9]{1,})?)[ ]((day|week|month)s?)$")

	if re.MatchString(s) {
		return TimeframeSchedular(s, numberOfSchedule)
	}

	return DateTimeStringSchedular(s, dateFormat)
}

// TimeframeSchedular ...
func TimeframeSchedular(s string, numberOfSchedule int) ([]int64, error) {
	var ts []int64
	re := regexp.MustCompile("^(every)[ ](0?[1-9]([0-9]{1,})?)[ ]((day|week|month)s?)$")

	if !re.MatchString(s) {
		err := errors.New("invalid string")
		return ts, err
	}

	sArr := strings.Split(strings.ToLower(s), " ")

	scheduleNumber, err := strconv.Atoi(sArr[1])
	if err != nil {
		return ts, err
	}

	const oneMillisecInNanosec = 1e6
	now := time.Now()
	for num := 0; num < numberOfSchedule; num++ {
		if num == 0 {
			// convert date to milliseconds
			tMilliSec := now.UnixNano() / oneMillisecInNanosec
			ts = append(ts, tMilliSec)
			continue
		}

		if sArr[2] == "month" || sArr[2] == "months" {
			// convert date to milliseconds
			tMilliSec := now.AddDate(0, scheduleNumber, 0).UnixNano() / oneMillisecInNanosec
			ts = append(ts, tMilliSec)
			continue
		}

		var hours string
		numOfHoursInADay := 24

		if sArr[2] == "day" || sArr[2] == "days" {
			hours = fmt.Sprintf("%dh", numOfHoursInADay*scheduleNumber*num)
		} else if sArr[2] == "week" || sArr[2] == "weeks" {
			numOfDaysInAWeek := 7
			hours = fmt.Sprintf("%dh", numOfHoursInADay*numOfDaysInAWeek*scheduleNumber*num)
		}

		duration, err := time.ParseDuration(hours)
		if err != nil {
			return ts, err
		}

		// convert date to milliseconds
		tMilliSec := now.Add(duration).UnixNano() / oneMillisecInNanosec
		ts = append(ts, tMilliSec)
	}

	return ts, err
}

// DateTimeStringSchedular ...
func DateTimeStringSchedular(s, dateFormat string) ([]int64, error) {
	var ts []int64

	sArr := strings.Split(s, ",")

	if len(sArr) < 1 {
		return ts, fmt.Errorf("invalid string")
	}

	for _, val := range sArr {
		t, err := time.Parse(dateFormat, val)
		if err != nil {
			return ts, fmt.Errorf("invalid string")
		}

		const oneMillisecInNanosec = 1e6
		// convert date to milliseconds
		tMillisec := t.UnixNano() / oneMillisecInNanosec
		nowMillisec := time.Now().UnixNano() / oneMillisecInNanosec

		// check if time has expired
		if (tMillisec - nowMillisec) < 0 {
			return ts, errors.New("time has expired")
		}

		ts = append(ts, tMillisec)
	}

	return ts, nil
}

// Address URI to smtp server
func (s *SMTPServer) Address() string {
	return fmt.Sprintf("%s:%s", s.Host, s.Port)
}

// SendEmail ...
func (s *SMTPServer) SendEmail(from, password, subject, message string, to []string) error {
	mime := "MIME-version: 1.0;\nContent-Type: text/plain; charset=\"UTF-8\";\n\n"
	sbj := fmt.Sprintf("Subject: %s \n", subject)
	tpl := fmt.Sprintf("%s%s\n%s", sbj, mime, message)
	msg := []byte(tpl)

	// Authentication.
	auth := smtp.PlainAuth("", from, password, s.Host)
	// Sending email.
	return smtp.SendMail(s.Address(), auth, from, to, msg)
}

// mailTemplate ...
// const mailTemplate = `
// 	<!DOCTYPE html>
// 	<html>
// 		<head>
// 			<meta charset="UTF-8">
// 			<title>{{.Title}}</title>
// 		</head>
// 		<body>
// 			<h1>{{.CourseTitle}}</h1>
// 			<h3>{{.ModuleTitle}}</h3>
// 			Click <a href={{.ModuleLink}}>here</a> to access the module.
// 		</body>
// 	</html>
// `

// GenerateMailTemplate ...
func GenerateMailTemplate(templateFileName string, data *MailTemplateData) (tpl string, err error) {
	t, err := template.ParseFiles(templateFileName)
	if err != nil {
		return
	}

	buf := new(bytes.Buffer)
	if err = t.Execute(buf, data); err != nil {
		return
	}

	tpl = buf.String()
	return
}
