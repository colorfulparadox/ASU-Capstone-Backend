package routes

import (
	"BackEnd/router"
)

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

func AI_Routes(r router.Router) {
	router.AddRoute(&r, router.Receiver{
		Route:     "/add_menu",
		RouteType: router.RoutePost,
		Sender:    Add_Menu,
	})
}

func Persona_Routes(r router.Router) {
	router.AddRoute(&r, router.Receiver{
		Route:     "/start_conversation",
		RouteType: router.RoutePost,
		Sender:    Start_Conversation,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/continue_conversation",
		RouteType: router.RoutePost,
		Sender:    Continue_Conversation,
	})

	router.AddRoute(&r, router.Receiver{
		Route:     "/end_conversation",
		RouteType: router.RoutePost,
		Sender:    End_Conversation,
	})
}
