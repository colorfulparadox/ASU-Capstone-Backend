// Basic file simply to hold commonly used data for other routes

package routes_user

import (
	"BackEnd/router"
)

type UserData struct {
	AuthID           string `json:"authID"`
	Edit_User        string `json:"edit_user"`
	Name             string `json:"name"`
	Username         string `json:"username"`
	Password         string `json:"password"`
	Average_Points   int    `json:"points"`
	Sentiment_Points int    `java:"sentiment_points"`
	Sales_Points     int    `java:"sales_points"`
	Knowledge_Points int    `java:"knowledge_points"`
	PermissionLevel  int    `json:"permission_level"`
	Email            string `json:"email"`
}

// StandardResult is a basic JSON format for returning one of the result enum types ()
type StandardResult struct {
	Result int `json:"result"`
}

func User_Routes(r router.Router) {
	router.AddRoute(&r, router.Receiver{
		Route:     "/login",
		RouteType: router.RoutePost,
		Sender:    Login,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/create_user",
		RouteType: router.RoutePost,
		Sender:    Create_User,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/update_user",
		RouteType: router.RoutePost,
		Sender:    Update_User,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/reset_auth_id",
		RouteType: router.RoutePost,
		Sender:    Reset_Auth_Token,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/delete_user",
		RouteType: router.RoutePost,
		Sender:    Delete_User,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/authenticate",
		RouteType: router.RoutePost,
		Sender:    Authenticate,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/modify_points",
		RouteType: router.RoutePost,
		Sender:    Modify_Points,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/user_list",
		RouteType: router.RoutePost,
		Sender:    User_List,
	})
}
