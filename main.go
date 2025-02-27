package main

import (
	"BackEnd/database"
	"BackEnd/router"
	"BackEnd/routes"
	"BackEnd/routes_persona"
	"BackEnd/routes_user"

	"github.com/joho/godotenv"
)

// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// https://go.dev/doc/tutorial/web-service-gin

func main() {
	godotenv.Load()

	database.Create_Tables()
	database.Initalize_Directories()

	//database.Add_Menu("a4abd78d-828e-460c-b6ab-7474d6b490b4/ce592ee1/3R1aF9Tkf616mulTLxG9hA==", "default", `{"test": "here"}`)

	r := router.NewRouter(":4040")

	router.AddRoute(&r, router.Receiver{
		Route:     "/ping",
		RouteType: router.RouteGet,
		// This assigns the user role middleware and requires the user value in the header to be "role"
		// Middleware: router.User_Role_Middleware("role"),
		Sender: routes.Ping,
	})

	routes_user.User_Routes(r)
	// routes_ai.AI_Routes(r)
	routes_persona.Persona_Routes(r)

	router.RunRouter(r)
}
