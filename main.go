package main

import (
	"BackEnd/database"
	"BackEnd/router"
	"BackEnd/routes"

	"github.com/joho/godotenv"
)

// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// https://go.dev/doc/tutorial/web-service-gin

func main() {
	godotenv.Load()

	database.Create_Tables()
	database.Initalize_Directories()

	r := router.NewRouter(":4040")

	router.AddRoute(&r, router.Receiver{
		Route:     "/ping",
		RouteType: router.RouteGet,
		// This assigns the user role middleware and requires the user value in the header to be "role"
		// Middleware: router.User_Role_Middleware("role"),
		Sender: routes.Ping,
	})

	routes.User_Routes(r)
	routes.AI_Routes(r)
	routes.Persona_Routes(r)

	router.RunRouter(r)
}
