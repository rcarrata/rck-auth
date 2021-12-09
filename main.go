package main

import (
	"log"
	"net/http"
	"os"

	"github.com/gorilla/handlers"
	"github.com/rcarrata/rck-auth/pkg/utils"
	"github.com/rcarrata/rck-auth/pkg/database"
)

func main() {
	InitialMigration()
	CreateRouter()
	InitializeRoute()
	StartServer()
}
