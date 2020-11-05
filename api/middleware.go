package api

import (
	"errors"
	"os"
	"time"

	jwt "github.com/appleboy/gin-jwt/v2"
	"github.com/aureleoules/epitaf/models"
	"github.com/gin-gonic/gin"
)

// AuthMiddleware handles JWT authentications
func AuthMiddleware() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "epitaf",
		Key:        []byte(os.Getenv("JWT_SECRET")),
		Timeout:    time.Hour * 48,
		MaxRefresh: time.Hour * 48,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			u := data.(*models.User)
			// What we put in the JWT claims
			return jwt.MapClaims{
				"uuid":      u.UUID.String(),
				"email":     u.Email,
				"name":      u.Name,
				"login":     u.Login,
				"promotion": u.Promotion,
				"class":     u.Class,
				"region":    u.Region,
				"semester":  u.Semester,
				"teacher":   u.Teacher,
			}
		},
		Authenticator: callbackHandler,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			_ = c.AbortWithError(code, errors.New(message))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	// Authmiddleware must be active
	if err != nil {
		panic(err)
	}

	return authMiddleware
}
