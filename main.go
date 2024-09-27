package main

import (
	"BackEnd/database"
	"BackEnd/router"
	"BackEnd/routes"
)

// https://pkg.go.dev/github.com/gin-gonic/gin#section-readme
// https://go.dev/doc/tutorial/web-service-gin

type message struct {
	Message string `json:"message"`
}

func main() {
	database.Create_Tables()

	//database.Randomize_auth_token("e5eb13a7-cea0-414b-9391-80627e6bb321/cded7614/H64bRKmgCPTyPaWZ1wR-Zg==")
	//database.Verify_User_Auth_Token("4df4bfb9-476c-4a05-a642-254c0b68b495/cded6d7f/7mBC5dHsv2SvklUpInkjng==")
	//return

	r := router.NewRouter("localhost", "4040")

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
		Route:     "/reset_auth_token",
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

	/*
		add_route(&router, Receiver{
			route:      "/login",
			routeType:  RoutePost,
			middleware: default_middleware,
			sender: func(gc *gin.Context, pool *pgxpool.Pool) {
				log.Println("login req")

				var loginReq LoginRequest

				err := gc.ShouldBindJSON(&loginReq)
				if err != nil {
					gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
					return
				}

				log.Println(loginReq)

				gc.JSON(http.StatusOK, "{\"auth\":\"thisisakey123\"}")
			},
		})
	*/

	router.RunRouter(r)
}
