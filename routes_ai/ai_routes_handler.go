package routes_ai

import "BackEnd/router"

type Conversation struct {
	AuthID         string `json:"authID"`
	Message        string `json:"message"`
	ConversationID string `json:"conversationID"`
}

type StandardResult struct {
	Result int `json:"result"`
}

func AI_Routes(r router.Router) {
	router.AddRoute(&r, router.Receiver{
		Route:     "/add_menu",
		RouteType: router.RoutePost,
		Sender:    Add_Menu,
	})
}
