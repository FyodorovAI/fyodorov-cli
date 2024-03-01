package authenticator

import "github.com/gin-gonic/gin"

var (
	endpoints = make(map[string]gin.HandlerFunc)
	router    = gin.Default()
)

func init() {
	auth, err := New()
	if err != nil {
		panic(err)
	}
	router.GET("/", IsAuthenticated)
	router.GET("/sign_in", LoginHandler(auth))
	router.GET("/callback", CallbackHandler(auth))
	router.GET("/logout", LogoutHandler)
}
