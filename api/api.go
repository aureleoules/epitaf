package api

import (
	"net/http"

	jwt "github.com/appleboy/gin-jwt/v2"
	// Import GoSwagger
	_ "github.com/aureleoules/epitaf/docs"
	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	"github.com/swaggo/gin-swagger/swaggerFiles"
	"go.uber.org/zap"
)

const version = "v1"

var router *gin.RouterGroup
var auth *jwt.GinJWTMiddleware

func createRouter() *gin.Engine {
	// Create router
	r := gin.Default()

	// Use CORS
	r.Use(cors())

	// Swagger
	url := ginSwagger.URL("/swagger/doc.json")
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler, url))
	r.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, "/swagger/index.html")
	})

	// Default API route
	router = r.Group("/" + version)
	// JWT middleware
	auth = AuthMiddleware()

	// Do not apply auth middleware here
	handleAuth()
	handleZeus()

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
