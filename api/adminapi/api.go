package adminapi

import (
	"os"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"go.uber.org/zap"
)

const version = "v1"

var router *echo.Group

// Response type
type resp map[string]interface{}

func createRouter() *echo.Echo {
	// Create router
	e := echo.New()

	// Use CORS
	e.Use(middleware.CORS())

	// Default API route
	router = e.Group("/" + version)

	// Do not apply auth middleware here
	handleAuth()
	jwtConfig := middleware.DefaultJWTConfig
	jwtConfig.TokenLookup = "header:" + echo.HeaderAuthorization
	jwtConfig.ContextKey = "user"
	jwtConfig.SigningKey = []byte(os.Getenv("JWT_SECRET"))
	jwtConfig.AuthScheme = "Bearer"
	jwtConfig.SigningMethod = "HS256"
	// Apply auth middleware
	router.Use(middleware.JWTWithConfig(jwtConfig))

	handleAdmins()
	handleUsers()
	handleTasks()
	handleGroups()
	handleRealms()

	return e
}

// Serve private api
func Serve() {
	r := createRouter()

	if err := r.Start(":" + os.Getenv("API_ADMIN_PORT")); err != nil {
		zap.S().Panic(err)
	}
}
