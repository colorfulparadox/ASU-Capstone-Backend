package routes

import (
	"BackEnd/database"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/jackc/pgx/v5/pgxpool"
)

type LoginRequest struct {
	User string `json:"user"`
	Pass string `json:"pass"`
}

// type LoginToken struct {
// 	AuthID     string `json:"authID"`
// 	DateIssued int64  `json:"dateIssued"`
// 	Expires    int64  `json:"expires"`
// }

// Holds the auth_token
type LoginToken struct {
	AuthID string `json:"authID"`
}

func Login(gc *gin.Context, pool *pgxpool.Pool) {
	var loginReq LoginRequest

	// Parses JSON received from client
	err := gc.ShouldBindJSON(&loginReq)
	if err != nil {
		gc.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Gets the auth_token for the specifc user
	auth_token := database.Verify_User_Login(loginReq.User, loginReq.Pass)

	// Checks if user is valid
	if auth_token == "" {
		log.Println("username or password incorrect")
		gc.JSON(http.StatusForbidden, "{}")
	}

	// Puts auth_token into JSON object
	loginToken := LoginToken{
		AuthID: auth_token,
	}

	// Returns loginToken
	gc.JSON(http.StatusOK, loginToken)

	// var username string = ""
	// var password string

	// pool.QueryRow(
	// 	context.Background(),
	// 	"SELECT username, password FROM users WHERE username = $1",
	// 	loginReq.User,
	// ).Scan(&username, &password)

	// if loginReq.Pass != password || username == "" {
	// 	fmt.Println("invalid password")
	// 	gc.JSON(http.StatusForbidden, "{}")
	// }

	// dateIssued := time.Now().Unix()
	// expires := time.Now().Add(2 * 24 * time.Hour).Unix()

	// //ignore this stuff for now
	// key := uuid.New().String() + "/" + strconv.FormatInt(((int64)(rand.Intn(9000))*expires)+dateIssued^(int64)(rand.Intn(9999)), 16)
	// hash := sha256.New()
	// hash.Write([]byte(key))
	// key64 := base64.URLEncoding.EncodeToString(hash.Sum(nil))

	// fmt.Println(key)
	// fmt.Println(key64)

	// loginToken := LoginToken{
	// 	AuthID:     key64,
	// 	DateIssued: dateIssued,
	// 	Expires:    expires,
	// }

	//gc.JSON(http.StatusOK, loginToken)
}
