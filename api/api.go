package api

import (
	jwt "github.com/appleboy/gin-jwt/v2"
	// Import GoSwagger
	_ "github.com/aureleoules/epitaf/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
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
	url := ginSwagger.URL("http://localhost:8080/swagger/doc.json") // The url pointing to API definition
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))

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
