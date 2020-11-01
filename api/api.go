package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
)

var api *gin.RouterGroup
var auth *jwt.GinJWTMiddleware

// Serve private api
func Serve() {
	// Create router
	r := gin.Default()

	// Use CORS
	r.Use(cors())

	// Default API route
	api = r.Group("/api")
	// JWT middleware
	auth = authMiddleware()

	// Do not apply auth middleware here
	handleAuth()

	// Apply auth middleware on these routes
	api.Use(auth.MiddlewareFunc())

	handleUsers()
	handleTasks()

	r.Run()
}
