package api

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
)

// isAuthOk returns a handler that executes some logic,
// and then calls the next handler.

func isAuthOk(handler http.HandlerFunc) http.HandlerFunc {

	// this middleware function uses an anonymous function to simplify
	return func(w http.ResponseWriter, r *http.Request) {

		if r.Header["Token"] == nil {
			var err Error
			err = SetError(err, "No Token Found")
			// Returns to the http response the Err struct in json format encoded
			json.NewEncoder(w).Encode(err)
		}

		// Define the SigningKey var and convert this secretkey into byte
		var newSigningKey = []byte(secretkey)

		// Received Token from the Header when the request is performed
		receivedToken := r.Header["Token"][0]
		logrus.Println(receivedToken)

		// Parsing and Validating the token received in the request using the HMAC signing method
		// https://pkg.go.dev/github.com/golang-jwt/jwt@v3.2.2+incompatible#Parse
		token, err := jwt.Parse(receivedToken, func(token *jwt.Token) (interface{}, error) {

			// Parse takes the token string and a function for looking up the key. The latter is especially
			// useful if you use multiple keys for your application.
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				// Validate the Token and return an error if the signing token is not the proper one
				// TODO: change fmt -> logrus
				return nil, fmt.Errorf("unexpected signing method: %v in token of type: %v", token.Header["alg"], token.Header["typ"])
			}

			// logrus.Println(token.Header["alg"])
			// logrus.Println(token.Header["typ"])
			return newSigningKey, nil
		})

		// If the token Parser have an error the Token is considered as Expired
		// TODO: Improve with the jwt.ValidationErrorMalformed
		if err != nil {
			var err Error
			err = SetError(err, "Your JWT Token is wrong or is expired")
			// Returns to the http response the Err struct in json format encoded
			w.WriteHeader(http.StatusUnauthorized)
			json.NewEncoder(w).Encode(err)
			return
		}

		// claims are actually a map[string]interface{}
		// Check if the token provided is Valid
		if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
			logrus.Println("Token is Valid")
			logrus.Println(claims["role"])
			if claims["role"] == "admin" {
				logrus.Println("Assigned Token role to Admin")
				r.Header.Set("Role", "admin")
				handler.ServeHTTP(w, r)
				return
			} else if claims["role"] == "user" {
				logrus.Println("Assigned Token role to User")
				r.Header.Set("Role", "user")
				handler.ServeHTTP(w, r)
				return
			} else {
				// If the role is not admin or user, send 403 status code Unauthorized
				var reserr Error
				reserr = SetError(reserr, "Role Not Authorized.")
				// Returns to the http response the Err struct in json format encoded
				w.WriteHeader(http.StatusForbidden)
				json.NewEncoder(w).Encode(reserr)
				return
			}
		}

	}
}
