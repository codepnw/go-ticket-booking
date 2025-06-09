package helper

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

type Auth struct {
	secret string
}

func SetupAuth(secret string) Auth {
	return Auth{secret: secret}
}

func (a Auth) GenenrateHashPassword(password string) (string, error) {
	if len(password) < 6 {
		return "", errors.New("password length be at least 6 characters long")
	}

	hashed, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", errors.New("password hash failed")
	}

	return string(hashed), nil
}

func (a Auth) VerifyPassword(password, hashedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(password))
	if err != nil {
		return errors.New("password does not matchh")
	}

	return nil
}

// Token : Generate
func (a Auth) GenerateToken(id int64, email, role string) (string, error) {
	if id == 0 || email == "" {
		return "", errors.New("required input are missing")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(time.Hour * 24 * 30).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(a.secret))
	if err != nil {
		return "", errors.New("signed token failed")
	}

	return tokenStr, nil
}

// Token : Verify
func (a Auth) VerifyToken(t string) (*domain.User, error) {
	tokenArr := strings.Split(t, " ")
	if len(tokenArr) != 2 {
		return nil, errors.New("invalid token format")
	}

	if tokenArr[0] != "Bearer" {
		return nil, errors.New("invalid token format")
	}

	token, err := jwt.Parse(tokenArr[1], func(t *jwt.Token) (interface{}, error) {
		if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unknow signing method: %v", t.Header)
		}
		return []byte(a.secret), nil
	})
	if err != nil {
		return nil, errors.New("invalid signing method")
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		if float64(time.Now().Unix()) > claims["exp"].(float64) {
			return nil, errors.New("token is expired")
		}

		user := domain.User{}
		user.ID = int64(claims["user_id"].(float64))
		user.Email = claims["email"].(string)
		user.Role = claims["role"].(string)

		return &user, nil
	}

	return nil, errors.New("token verification failed")
}

func (a Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.GetReqHeaders()["Authorization"]

	if authHeader == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "auth header is missing",
		})
	}

	user, err := a.VerifyToken(authHeader[0])
	if err == nil && user.ID > 0 {
		ctx.Locals("user", user)
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"error":   err,
		})
	}
}
