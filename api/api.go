package api

import (
	"errors"
	"os"
	"time"

	"github.com/aureleoules/epitaf/models"

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

	handleUsers()

	// api.Use(auth.MiddlewareFunc())

	r.Run()
}

func authMiddleware() *jwt.GinJWTMiddleware {
	// the jwt middleware
	authMiddleware, err := jwt.New(&jwt.GinJWTMiddleware{
		Realm:      "epitaf",
		Key:        []byte(os.Getenv("JWT_SECRET")),
		Timeout:    time.Hour * 48,
		MaxRefresh: time.Hour * 48,
		PayloadFunc: func(data interface{}) jwt.MapClaims {
			u := data.(models.User)
			return jwt.MapClaims{
				"uuid":      u.UUID.String(),
				"email":     u.Email,
				"name":      u.Name,
				"promotion": u.Promotion,
				"class":     u.Class,
			}
		},
		Authenticator: callbackHandler,
		Authorizator: func(data interface{}, c *gin.Context) bool {
			return true
		},
		Unauthorized: func(c *gin.Context, code int, message string) {
			c.AbortWithError(code, errors.New(message))
		},
		TokenLookup:   "header: Authorization, query: token, cookie: jwt",
		TokenHeadName: "Bearer",
		TimeFunc:      time.Now,
	})

	if err != nil {
		panic(err)
	}

	return authMiddleware
}
