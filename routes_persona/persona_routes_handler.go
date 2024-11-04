package routes_persona

import "BackEnd/router"

type Conversation struct {
	AuthID         string `json:"authID"`
	Message        string `json:"message"`
	ConversationID string `json:"conversationID"`
}

type StandardResult struct {
	Result int `json:"result"`
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
