package main

import (
	"BackEnd/database"
	"BackEnd/router"
	"BackEnd/routes"
)

const (
	Temp_Images = "/temp"
	Thumbnails  = "/thumbnails"
)

// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// https://go.dev/doc/tutorial/web-service-gin

func main() {
	database.Create_Tables()

	r := router.NewRouter(":4040")

	router.AddRoute(&r, router.Receiver{
		Route:     "/ping",
		RouteType: router.RouteGet,
		// This assigns the user role middleware and requires the user value in the header to be "role"
		// Middleware: router.User_Role_Middleware("role"),
		Sender: routes.Ping,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/login",
		RouteType: router.RoutePost,
		Sender:    routes.Login,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/create_user",
		RouteType: router.RoutePost,
		Sender:    routes.Create_User,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/update_user",
		RouteType: router.RoutePost,
		Sender:    routes.Update_User,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/reset_auth_id",
		RouteType: router.RoutePost,
		Sender:    routes.Reset_Auth_Token,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/delete_user",
		RouteType: router.RoutePost,
		Sender:    routes.Delete_User,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/authenticate",
		RouteType: router.RoutePost,
		Sender:    routes.Authenticate,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/modify_points",
		RouteType: router.RoutePost,
		Sender:    routes.ModifyPoints,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/user_list",
		RouteType: router.RoutePost,
		Sender:    routes.UserList,
	})

	router.RunRouter(r)
}
