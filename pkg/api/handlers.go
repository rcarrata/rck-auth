package api

import (
	"encoding/json"
	"net/http"

	"github.com/sirupsen/logrus"
	db "github.com/rcarrata/rck-auth/pkg/database"


)

func SignUp(w http.ResponseWriter, r *http.Request) {
	connection := db.GetDatabase()
	defer db.Closedatabase(connection)

	user := db.User{}
	// Extract from the Body the Email/Password struct inputs and store into new memory address of new struct
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		var err Error
		err = SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Returns to the http response the Err struct in json format encoded
		json.NewEncoder(w).Encode(err)
		return
	}

	// Define a empty struct for the User
	loginuser := db.User{}

	// Retrieve the first matched record of a User (struct) in the database and compare it
	// with the user.Email that was sent in the POST request Body.
	// Use the Where GORM sentence for
	// Gorm Conditionals: https://gorm.io/docs/query.html#String-Conditions

	// Select first matched record email == user.Email() and store the result into the User{} struct
	// If the select is empty, the user/email is not present and you can create it
	connection.Where("email = ?", user.Email).First(&loginuser)

	// Check if the Email is already registered or not
	// If the output of the struct loginuser have NOT the Email empty after the Where clause
	// the email is already repeated
	if loginuser.Email != "" {
		// For debugging purposes
		logrus.Println("The User:", loginuser.Email, "with Password:", loginuser.Password)

		err := Error{}
		err = SetError(err, "Email already in use")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(err)
		return
	}

	logrus.Println("Generating the Password for user", user.Email)

	// Update to the user.Password param in struct the generated password
	user.Password, err = GeneratePass(user.Password)
	if err != nil {
		logrus.Fatalln("Error in password generation hash")
	}
	logrus.Println("Password: ", user.Password)

	// Create a new user with the struct of the User updated
	connection.Create(&user)
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	// Return the request with the User struct
	json.NewEncoder(w).Encode(user)

}

// SignIn function that checks if the user is present in the system, and check the key:values
// stored in the database. After compares the values from input and output and if its ok,
// generates a Golang JWT authentication
func SignIn(w http.ResponseWriter, r *http.Request) {
	// Connect to the Database
	connection := db.GetDatabase()

	// Defer the close of the Database to the end of the
	defer db.Closedatabase(connection)

	// Read from Request Body the auth input email and pass and store it in a Struct Authentication
	var authdetails Authentication
	err := json.NewDecoder(r.Body).Decode(&authdetails)

	// Raise an error if the Body is not well formatted or if have not the proper structure
	if err != nil {
		var err Error
		err = SetError(err, "Error in reading body")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		// Create a decoder with the request body and call the decode method with the pointer of the struct
		json.NewEncoder(w).Encode(err)
		return
	}

	// Authuser struct defined to store values from the DB
	authuser := db.User{}
	// Check the email that the User sends when sends the request, and it's stored in the Authentication struct defined before
	connection.Where("email = ?", authdetails.Email).First(&authuser)

	// If the User/Email is empty represents that the email introduced are not present into the database.
	if authuser.Email == "" {
		var err Error
		err = SetError(err, "Your email is not registered. Please first do the signup!")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(err)
		return
	}

	logrus.Info("Authdetails - User Request: ", authdetails.Password)
	logrus.Info("AuthUser - DB Stored: ", authuser.Password)
	// authdetails struct storing values from the User request to the API
	check := CheckPass(authdetails.Password, authuser.Password)

	// Check if the bool of the return err from the CheckPass is nil
	if !check {
		var err Error
		err = SetError(err, "Username or Password is incorrect. Please review them!")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(err)
		return
	}

	// Generate the JWT Token using the Email from the authuser Email and Roles stored into the DB
	validToken, err := GenerateJWT(authuser.Email, authuser.Role)
	if err != nil {
		var err Error
		err = SetError(err, "Failed to generate the token")
		w.Header().Set("Content-Type", "application/json; charset=UTF-8")
		json.NewEncoder(w).Encode(err)
		return
	}

	// Initialize a empty Token struct in a variable token
	token := Token{}
	// Define the email and role that is stored in the DB
	token.Email = authuser.Email
	token.Role = authuser.Role
	// Define the TokenString with the value of the Token generated
	token.TokenString = validToken
	logrus.Info("Generated JWT Token: ", token.TokenString)

	// Send the TokenString generated back to the user as response of the signin
	w.Header().Set("Content-Type", "application/json; charset=UTF-8")
	json.NewEncoder(w).Encode(token.TokenString)

}

func AdminIndex(w http.ResponseWriter, r *http.Request) {
	logrus.Println("Role -> ", r.Header.Get("Role"))
	if r.Header.Get("Role") != "admin" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("You are not authorized. Admin Only!"))
		return
	}
	w.Write([]byte("Welcome, Admin."))
}

func UserIndex(w http.ResponseWriter, r *http.Request) {
	logrus.Println("Role -> ", r.Header.Get("Role"))
	if r.Header.Get("Role") != "user" {
		w.WriteHeader(http.StatusForbidden)
		w.Write([]byte("Not authorized. User Only!!"))
		return
	}
	w.Write([]byte("Welcome, User."))
}
