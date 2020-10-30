package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var api *gin.RouterGroup
var auth *jwt.GinJWTMiddleware

// Serve private api
func Serve() {
	r := gin.Default()

	r.Use(cors())

	api = r.Group("/api")
	auth = authMiddleware()

	// Do not apply auth middleware here
	handleAuth()

	// Apply auth middleware on these routes
	api.Use(auth.MiddlewareFunc())
	handleUsers()
	handleTasks()

	r.Run()
}
