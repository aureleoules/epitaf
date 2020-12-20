package adminapi

import (
	"os"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/api/middleware"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

const version = "v1"

var router *gin.RouterGroup
var auth *jwt.GinJWTMiddleware

func createRouter() *gin.Engine {
	// Create router
	r := gin.Default()

	// Use CORS
	r.Use(middleware.Cors())

	// Default API route
	router = r.Group("/" + version)
	// JWT middleware
	auth = middleware.AuthMiddleware(authenticator)

	// Do not apply auth middleware here
	handleAuth()

	// Apply auth middleware on these routes
	router.Use(auth.MiddlewareFunc())

	handleUsers()
	handleTasks()
	handleGroups()

	return r
}

// Serve private api
func Serve() {
	r := createRouter()

	if err := r.Run(":" + os.Getenv("API_ADMIN_PORT")); err != nil {
		zap.S().Panic(err)
	}
}
