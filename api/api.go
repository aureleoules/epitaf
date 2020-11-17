package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

var router *gin.RouterGroup
var auth *jwt.GinJWTMiddleware

func createRouter() *gin.Engine {
	// Create router
	r := gin.Default()

	// Use CORS
	r.Use(cors())

	// Default API route
	router = r.Group("/api")
	// JWT middleware
	auth = AuthMiddleware()

	// Do not apply auth middleware here
	handleAuth()

	// Apply auth middleware on these routes
	router.Use(auth.MiddlewareFunc())

	handleUsers()
	handleTasks()
	handleClasses()

	return r
}

// Serve private api
func Serve() {
	r := createRouter()

	if err := r.Run(); err != nil {
		zap.S().Panic(err)
	}
}
