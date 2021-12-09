package main

import (
	api "github.com/rcarrata/rck-auth/pkg/api"
	db "github.com/rcarrata/rck-auth/pkg/database"
)

func main() {
	db.InitialMigration()
	api.CreateRouter()
	api.InitializeRoute()
	api.StartServer()
}
