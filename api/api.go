package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
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
	auth = AuthMiddleware()

	// Do not apply auth middleware here
	handleAuth()

	// Apply auth middleware on these routes
	api.Use(auth.MiddlewareFunc())

	handleUsers()
	handleTasks()
	handleClasses()

	if err := r.Run(); err != nil {
		zap.S().Panic(err)
	}
}
