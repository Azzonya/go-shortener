package auth

import (
	"github.com/Azzonya/go-shortener/internal/user"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestAuth_CreateJWTCookie(t *testing.T) {
	auth := New("testSecret")
	user := &user.User{ID: "testUserID"}

	cookie, err := auth.CreateJWTCookie(user)

	assert.NoError(t, err)
	assert.Equal(t, sessionCookie, cookie.Name)
	assert.NotEmpty(t, cookie.Value)
}

func TestAuth_GetUserFromJWT(t *testing.T) {
	auth := New("testSecret")
	user := &user.User{ID: "testUserID"}
	token, _ := auth.NewToken(user)

	resultUser, err := auth.GetUserFromJWT(token)

	assert.NoError(t, err)
	assert.Equal(t, user, resultUser)
}

func TestAuth_NewToken(t *testing.T) {
	auth := New("testSecret")
	user := &user.User{ID: "testUserID"}

	token, err := auth.NewToken(user)

	assert.NoError(t, err)
	assert.NotEmpty(t, token)
}
