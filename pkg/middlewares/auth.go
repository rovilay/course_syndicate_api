package middlewares

import (
	"context"
	"errors"
	"net/http"

	"course_syndicate_api/pkg/utils"
)

// Authenticate ...
func Authenticate(next http.Handler) func(http.ResponseWriter, *http.Request) {
	return func(res http.ResponseWriter, r *http.Request) {
		e := &utils.ErrorWithStatusCode{}
		token := r.Header.Get("Authorization")

		if token == "" {
			e.StatusCode = http.StatusUnauthorized
			e.ErrorMessage = errors.New("Authorization token not provided")

			utils.ErrorHandler(e, res)
			return
		}

		claims, err := utils.DecodeToken(token)
		if err != nil {
			e.StatusCode = http.StatusUnauthorized
			e.ErrorMessage = errors.New("invalid/expired token")

			utils.ErrorHandler(e, res)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, utils.ContextKey("claims"), claims)

		r = r.WithContext(ctx)
		next.ServeHTTP(res, r)
		return
	}
}
