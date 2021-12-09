package api

import (
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/sirupsen/logrus"
)

var (
	router *mux.Router
)

// Authentication is for login data
type Authentication struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

// Token
type Token struct {
	Role        string `json:"role"`
	Email       string `json:"email"`
	TokenString string `json:"token"`
}

type Error struct {
	IsError bool   `json:"isError"`
	Message string `json:"message"`
}

// CreateRouter generates a new instance of Mux Router
func CreateRouter() {
	router = mux.NewRouter()
}

// Home Page Handler (No Auth Required)
func homeHandler(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to the Home Page [No Auth Required]\n"))
}

// Time Page Handler (No Auth Required)
func timeHandler(w http.ResponseWriter, r *http.Request) {
	tm := time.Now().Format(time.RFC1123)
	w.Write([]byte("The time is: " + tm))
}

// InitializeRoute creates handlers for the mux Router to handle
func InitializeRoute() {
	router.HandleFunc("/", homeHandler).Methods("GET")
	router.HandleFunc("/time", timeHandler).Methods("GET")
	router.HandleFunc("/signup", SignUp).Methods("POST")
	router.HandleFunc("/signin", SignIn).Methods("POST")
	router.HandleFunc("/admin", isAuthOk(AdminIndex)).Methods("GET")
	router.HandleFunc("/user", isAuthOk(UserIndex)).Methods("GET")

	// Option Methods - Enabling CORS
	router.Methods("OPTIONS").HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "")
		w.Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
		w.Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, Access-Control-Request-Headers, Access-Control-Request-Method, Connection, Host, Origin, User-Agent, Referer, Cache-Control, X-header")
	})
}

func StartServer() {
	port := ":" + os.Getenv("PORT")
	logrus.Info("Server running in port ", port)
	//err := http.ListenAndServe(port, handlers.CORS(handlers.AllowCredentials))

	err := http.ListenAndServe(port, handlers.CORS(handlers.AllowedHeaders([]string{"X-Requested-With", "Access-Control-Allow-Origin", "Content-Type", "Authorization"}), handlers.AllowedMethods([]string{"GET", "POST", "PUT", "DELETE", "HEAD", "OPTIONS"}), handlers.AllowedOrigins([]string{"*"}))(router))
	if err != nil {
		log.Fatal(err)
	}
}
