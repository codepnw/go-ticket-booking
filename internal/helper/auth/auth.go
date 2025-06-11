package auth

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/codepnw/go-ticket-booking/internal/domain"
	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
)

const UserCtxKey = "user"

type Auth struct {
	secret        string
	refreshSecret string
}

func SetupAuth(secret, refreshSecret string) Auth {
	return Auth{
		secret:        secret,
		refreshSecret: refreshSecret,
	}
}

func (a *Auth) GenerateAccessToken(id int64, email, role string) (string, error) {
	duration := time.Hour * 24
	return a.generateToken(id, email, role, duration, a.secret)
}

func (a *Auth) GenerateRefreshToken(id int64, email, role string) (string, time.Time, error) {
	duration := time.Hour * 24 * 7
	exp := time.Now().Add(duration)
	token, err := a.generateToken(id, email, role, duration, a.refreshSecret)
	return token, exp, err
}

func (a *Auth) VerifyAccessToken(token string) (*domain.User, error) {
	return a.verifyToken(token, a.secret)
}

func (a *Auth) VerifyRefreshToken(token string) (*domain.User, error) {
	return a.verifyToken(token, a.refreshSecret)
}

func (a *Auth) Authorize(ctx *fiber.Ctx) error {
	authHeader := ctx.GetReqHeaders()["Authorization"]

	if authHeader == nil {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "auth header is missing",
		})
	}

	user, err := a.VerifyAccessToken(authHeader[0])
	if err == nil && user.ID > 0 {
		ctx.Locals(UserCtxKey, user)
		return ctx.Next()
	} else {
		return ctx.Status(http.StatusUnauthorized).JSON(&fiber.Map{
			"message": "authorization failed",
			"error":   err,
		})
	}
}

func GetCurrentUser(ctx *fiber.Ctx) (*domain.User, bool) {
	user, ok := ctx.Locals(UserCtxKey).(*domain.User)
	return user, ok
}

// ----- private -----
func (a *Auth) generateToken(id int64, email, role string, duration time.Duration, key string) (string, error) {
	if id == 0 || email == "" {
		return "", errors.New("required input are missing")
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user_id": id,
		"email":   email,
		"role":    role,
		"exp":     time.Now().Add(duration).Unix(),
	})

	tokenStr, err := token.SignedString([]byte(key))
	if err != nil {
		return "", errors.New("signed token failed")
	}

	return tokenStr, nil
}

// Token : Verify
func (a *Auth) verifyToken(t string, key string) (*domain.User, error) {
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
		return []byte(key), nil
	})
	if err != nil {
		return nil, err
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
