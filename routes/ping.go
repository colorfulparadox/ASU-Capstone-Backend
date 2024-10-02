type PingTest struct {
	Ping string `json:"ping"`
}

func Ping(){
	pingTest := PingTest{
		Ping: "pong",
	}

	gc.JSON(http.StatusOK, pingResponse)
}
