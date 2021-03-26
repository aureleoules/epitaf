package adminapi

import (
	"net/http"
	"os"
	"time"

	"github.com/aureleoules/epitaf/models"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo/v4"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func handleAuth() {
	router.POST("/auth/login", loginHandler)
}

func generateToken(id models.UUID) (string, error) {
	token := jwt.New(jwt.SigningMethodHS256)
	claims := token.Claims.(jwt.MapClaims)
	claims["id"] = id.String()
	claims["exp"] = time.Now().Add(time.Hour * 48).Unix()

	t, err := token.SignedString([]byte(os.Getenv("JWT_SECRET")))
	if err != nil {
		return "", err
	}
	return t, nil
}

func loginHandler(c echo.Context) error {
	req := struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}{}

	err := c.Bind(&req)
	if err != nil {
		zap.S().Warn(err)
		return c.JSON(http.StatusNotAcceptable, resp{"error": "not acceptable"})
	}

	u, err := models.GetAdminByEmail(req.Email)
	if err != nil {
		return c.JSON(http.StatusUnauthorized, resp{"error": "wrong email or password"})
	}

	err = bcrypt.CompareHashAndPassword([]byte(u.Password), []byte(req.Password))
	if err != nil {
		return c.JSON(http.StatusUnauthorized, resp{"error": "wrong email or password"})
	}

	t, err := generateToken(u.ID)
	if err != nil {
		return c.JSON(http.StatusInternalServerError, resp{"error": err.Error()})
	}

	return c.JSON(http.StatusOK, resp{
		"token": t,
	})
}
