// Package auth provides functionalities for user authentication and JWT handling.
package auth

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"

	"github.com/Azzonya/go-shortener/internal/user"
)

const (
	sessionCookie              = "userID"
	defaultJWTCookieExpiration = 24 * time.Hour
)

// Claims represents the custom claims structure for JWT tokens.
type Claims struct {
	jwt.RegisteredClaims
	UID string
}

// Auth provides methods for user authentication and JWT handling.
type Auth struct {
	JwtSecret string // JwtSecret is the secret key used for signing JWT tokens.
}

// New creates a new instance of the Auth struct with the provided JWT secret key.
func New(jwtSecret string) *Auth {
	return &Auth{JwtSecret: jwtSecret}
}

// GetUserFromCookie retrieves the user information from the session cookie.
func (a *Auth) GetUserFromCookie(c *gin.Context) (*user.User, error) {
	userCookie, err := c.Cookie(sessionCookie)
	if err != nil {
		return nil, err
	}

	return a.GetUserFromJWT(userCookie)
}

// GetUserFromJWT retrieves the user information from the JWT token.
func (a *Auth) GetUserFromJWT(signedToken string) (*user.User, error) {
	token, err := jwt.ParseWithClaims(signedToken, &Claims{},
		func(t *jwt.Token) (interface{}, error) {
			if _, ok := t.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, fmt.Errorf("unexpected signing method: %v", t.Header["alg"])
			}
			return []byte(a.JwtSecret), nil
		})
	if err != nil {
		return nil, fmt.Errorf("token is not valid")
	}

	if !token.Valid {
		return nil, errors.New("token is not valid")
	}

	if claims, ok := token.Claims.(*Claims); !ok {
		return nil, errors.New("token does not contain user id")
	} else {
		return user.NewWithID(claims.UID), nil
	}
}

// CreateJWTCookie creates a new JWT cookie for the given user.
func (a *Auth) CreateJWTCookie(u *user.User) (*http.Cookie, error) {
	token, err := a.NewToken(u)
	if err != nil {
		return nil, fmt.Errorf("cannot create auth token: %w", err)
	}
	return &http.Cookie{
		Name:  sessionCookie,
		Value: token,
	}, nil
}

// NewToken creates a new JWT token for the given user.
func (a *Auth) NewToken(u *user.User) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(defaultJWTCookieExpiration)),
		},
		UID: u.ID,
	})

	signedToken, err := token.SignedString([]byte(a.JwtSecret))
	if err != nil {
		return "", fmt.Errorf("cannot sign jwt token: %w", err)
	}
	return signedToken, nil
}
