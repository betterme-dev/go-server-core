package auth

import (
	"context"
	"encoding/json"
	"net/http"
	"strings"

	log "github.com/sirupsen/logrus"
)

const JSONContentType = "application/json; charset=utf-8"

type ContextKey string

var CtxKey ContextKey = "auth"

func Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authToken := extractToken(r)
		log.Debugf("authToken: %s", authToken)
		if authToken == "" {
			log.Debug("Auth token is empty")
			unauthorized(w)

			return
		}

		authService := NewService()
		if !authService.AuthByBearerToken(authToken) {
			unauthorized(w)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, CtxKey, &authService)

		// Calls the next handler, which can be another middleware in the chain, or the final handler.
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func ServiceFromCtx(ctx context.Context) *Service {
	auth, ok := ctx.Value(CtxKey).(*Service)
	if !ok {
		return nil
	}

	return auth
}

func unauthorized(res http.ResponseWriter) {
	log.Debug("Unable to authorize request")
	type Unauthorized struct {
		Message string `json:"message"`
	}
	data := Unauthorized{Message: "Your request was made with invalid credentials."}
	jsonData, _ := json.Marshal(data)
	writeContentType(res, []string{JSONContentType})
	res.WriteHeader(http.StatusUnauthorized)
	if _, err := res.Write(jsonData); err != nil {
		log.Error(err)
	}
}

func writeContentType(res http.ResponseWriter, value []string) {
	header := res.Header()
	if val := header["Content-Type"]; len(val) == 0 {
		header["Content-Type"] = value
	}
}

// Extract Bearer token from a request
func extractToken(r *http.Request) string {
	bearToken := r.Header.Get("Authorization")
	strArr := strings.Split(bearToken, " ")
	if len(strArr) == 2 {
		return strArr[1]
	}

	return ""
}
