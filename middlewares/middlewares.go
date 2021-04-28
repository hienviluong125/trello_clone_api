package middlewares

import (
	"encoding/json"
	"errors"
	"net/http"
	"os"
	"strings"

	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/context"
)

type JsonResponse map[string]interface{}

func Authenticate(next http.HandlerFunc) http.HandlerFunc {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		authorizationHeader := r.Header.Get("authorization")

		if authorizationHeader == "" {
			json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": "An authorization header is required"})
			return
		}

		bearerToken := strings.Split(authorizationHeader, " ")

		if len(bearerToken) != 2 {
			json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": "Invalid authorization structure"})
			return
		}

		token, err := jwt.Parse(bearerToken[1], func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, errors.New("there was an error")
			}
			return []byte(os.Getenv("API_SECRET")), nil
		})

		if err != nil {
			json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": "An authorization header is required"})
			return
		}

		if !token.Valid {
			json.NewEncoder(w).Encode(JsonResponse{"success": false, "message": "Invalid token"})
			return
		}

		context.Set(r, "jwt", token.Claims)
		next(w, r)
	})
}
