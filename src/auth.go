package main

import (
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt/v4"
	"github.com/google/uuid"
)

const jwtSecret = "super-secret-key"

func authorizationMiddleware(handlerFunc http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		tokenString := r.Header.Get("x-jwt-token")

		if tokenString == "" {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "token header is not present"})
			return
		}

		token, err := validateJWT(tokenString)
		if err != nil || !token.Valid {
			WriteJSON(w, http.StatusForbidden, APIError{Error: "invalid token"})
			return
		}

		claims := token.Claims.(jwt.MapClaims)
		r.Header.Add("x-account-id", claims["accountId"].(string))

		handlerFunc(w, r)
	}
}

func validateJWT(token string) (*jwt.Token, error) {
	return jwt.Parse(token, func(token *jwt.Token) (interface{}, error) {
		// Don't forget to validate the alg is what you expect:
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		// hmacSampleSecret is a []byte containing your secret, e.g. []byte("my_secret_key")
		return []byte(jwtSecret), nil
	})
}

func createJWT(accountId uuid.UUID) (string, error) {
	claims := &jwt.MapClaims{
		"expiresAt": 15000,
		"accountId": accountId,
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(jwtSecret))
}
