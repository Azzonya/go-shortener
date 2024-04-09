package api

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"

	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/Azzonya/go-shortener/internal/session"
)

type Request struct {
	URL string `json:"url"`
}

type Response struct {
	Result string `json:"result"`
}

func (o *Rest) ShortenJSON(c *gin.Context) {
	var err error
	var exist bool
	var userID string

	req := &Request{}
	resp := Response{}

	err = c.BindJSON(req)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	user, ok := session.GetUserFromContext(c.Request.Context())
	if ok {
		userID = user.ID
	}

	resp.Result, err = o.shortener.ShortenAndSaveLink(req.URL, userID)
	if err != nil {
		resp.Result, exist = o.shortener.GetOneByOriginalURL(req.URL)
		if !exist {
			c.JSON(http.StatusBadRequest, gin.H{
				"message": "Failed to create short URL",
				"error":   err.Error(),
			})
			return
		}
	}

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Marshall result struct",
			"error":   err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")

	if exist {
		c.Data(http.StatusConflict, "application/json", resultJSON)

	} else {
		c.Data(http.StatusCreated, "application/json", resultJSON)
	}
}

func (o *Rest) Shorten(c *gin.Context) {
	var exist bool
	var userID string

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read request body")
		return
	}

	user, ok := session.GetUserFromContext(c.Request.Context())
	if ok {
		userID = user.ID
	}

	reqObj := strings.TrimSpace(string(body))

	outputURL, err := o.shortener.ShortenAndSaveLink(reqObj, userID)
	if err != nil {
		outputURL, exist = o.shortener.GetOneByOriginalURL(reqObj)
		if !exist {
			c.String(http.StatusBadRequest, "Failed to add line to inmemory")
			return
		}
	}

	c.Header("Content-Type", "text/plain")
	if exist {
		c.String(http.StatusConflict, outputURL)

	} else {
		c.String(http.StatusCreated, outputURL)
	}
}

func (o *Rest) Redirect(c *gin.Context) {
	shortURL, exist := c.Params.Get("id")
	if !exist {
		c.String(http.StatusBadRequest, "Failed to get ID")
		return
	}

	if o.shortener.IsDeleted(shortURL) {
		c.AbortWithStatus(http.StatusGone)
		return
	}

	URL, exist := o.shortener.GetOneByShortURL(shortURL)
	if !exist {
		c.String(http.StatusBadRequest, "Failed to get original URL")
		return
	}

	c.Header("Location", URL)
	c.Redirect(http.StatusTemporaryRedirect, URL)
}

func (o *Rest) ShortenURLs(c *gin.Context) {
	var userID string
	var URLs []*entities.ReqURL

	err := c.BindJSON(&URLs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	user, ok := session.GetUserFromContext(c.Request.Context())
	if ok {
		userID = user.ID
	}

	shortenedURLs, err := o.shortener.ShortenURLs(URLs, userID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Shorten URLs",
			"error":   err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusCreated, shortenedURLs)
}

func (o *Rest) ListAll(c *gin.Context) {
	var err error

	u, ok := session.GetUserFromContext(c.Request.Context())
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user",
			"error":   errors.New("failed to get user from context").Error(),
		})
		return
	}

	if u.IsNew() {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "no cookie",
			"error":   errors.New("no authorized").Error(),
		})
		return
	}

	userID := u.ID

	result, err := o.shortener.ListAll(userID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get urls",
			"error":   err.Error(),
		})
		return
	}

	if len(result) == 0 {
		c.JSON(http.StatusNoContent, nil)
		return
	}

	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusOK, result)
}

func (o *Rest) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}

func (o *Rest) DeleteURLs(c *gin.Context) {
	var err error
	var shortURLs []string
	err = c.BindJSON(&shortURLs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	u, err := session.GetUser(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"message": "Failed to get user",
			"error":   err,
		})
		return
	}

	userID := u.ID

	o.shortener.DeleteURLs(shortURLs, userID)

	c.AbortWithStatus(http.StatusAccepted)
}
