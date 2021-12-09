package api

import (
	"time"

	"github.com/golang-jwt/jwt"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
)

var (
	secretkey string = "secretkeyjwt"
)

func SetError(err Error, message string) Error {
	err.IsError = true
	err.Message = message
	return err
}

// Generate Password from a Hash
func GeneratePass(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	// logrus.Println(string(bytes))
	return string(bytes), err
}

func CheckPass(password string, hash string) bool {
	// CompareHashAndPassword compares a bcrypt hashed password with its possible plaintext equivalent.
	// Returns nil on success, or an error on failure.
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	// If it's OK, set the error to bool nil / empty
	return err == nil
}

// JWT are divided in three separated elements:
// - Header: consists in two parts (JWT + Signign Algorithm) in json format, encoded in base64url
// - Payload: contains the Claims (usually the user) and other additional data
// - Signature: result of Header + Payload encoded, a secret, the signing algorithm and signing the Header + Payload

// Generate JWT Token based in the email and in the role as input. Creates a token by the algorithm signing method (HS256) and adds authorized email,
// role, and exp into claims.
// Claims are pieces of info added into the tokens.
func GenerateJWT(email string, role string) (string, error) {

	// Add the signingkey and convert it to an array of bytes
	signingKey := []byte(secretkey)

	// Generate a token with the HS256 as the Signign Method
	token := jwt.New(jwt.SigningMethodHS256)
	// logrus.Info("JWT Token: ", token) // Debug purposes

	// jwt library defines a struct with the MapClaims for define the different claims
	// to include in our token payload content in key-value format
	claims := token.Claims.(jwt.MapClaims)

	// TODO: Explore the token.jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims("key": "value")}
	// https://pkg.go.dev/github.com/golang-jwt/jwt@v3.2.2+incompatible#example-New-Hmac

	// Adding to the claims Map, authorized, the email, role and exp
	claims["authorized"] = true
	claims["email"] = email
	claims["role"] = role
	claims["exp"] = time.Now().Add(time.Minute * 30).Unix()

	// To Debug Claims
	// logrus.Println("claims auth", claims["authorized"])
	// logrus.Println("claims email", claims["email"])
	// logrus.Println("claims role", claims["role"])
	// logrus.Println("claims time", claims["exp"])

	// Sign the token with the signingkey defined in the step before
	tokenStr, err := token.SignedString(signingKey)
	if err != nil {
		logrus.Fatalln("Error during the Signing Token:", err.Error())
		return "", err
	}
	// For debugging purposes
	logrus.Println("Token Signed: ", tokenStr)

	return tokenStr, err

	// TODO: add Parser Token to increase the security purposes
}
