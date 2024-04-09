// Package middleware provides HTTP middleware for the URL shortener application.
package middleware

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"

	"github.com/Azzonya/go-shortener/internal/auth"
	"github.com/Azzonya/go-shortener/internal/logger"
	"github.com/Azzonya/go-shortener/internal/session"
	"github.com/Azzonya/go-shortener/internal/user"
)

// AuthMiddleware is a Gin middleware function that handles authentication.
// It takes a JWT secret as input and returns a Gin middleware function.
// The middleware extracts the user information from the request's JWT cookie.
// If the user information is not found or an error occurs during extraction,
// a new user is created, a JWT cookie is set, and the user context is added to the request.
// If any errors occur during these processes, an internal server error response is sent.
func AuthMiddleware(jwtSecret string) gin.HandlerFunc {
	return func(c *gin.Context) {
		authorizer := auth.New(jwtSecret)
		u, err := authorizer.GetUserFromCookie(c)
		if err != nil {
			u, err = user.New()
			if err != nil {
				// we're helpless here
				logger.Log.Debug("cannot create new user", zap.Error(err))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			sessionCookie, errS := authorizer.CreateJWTCookie(u)
			if errS != nil {
				logger.Log.Debug("cannot create session cookie for user", zap.Error(errS))
				c.AbortWithStatus(http.StatusInternalServerError)
				return
			}

			c.Header("Set-Cookie", sessionCookie.String())
		}
		c.Request = c.Request.WithContext(session.SetUserContext(c.Request.Context(), u))
	}
}
