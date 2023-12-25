package api

import (
	"encoding/json"
	"github.com/Azzonya/go-shortener/internal/entities"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
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

	o.shortener.UserID, err = c.Cookie("userID")
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"message": "Failed to get cookie - userID",
	//		"error":   err.Error(),
	//	})
	//	return
	//}

	resp.Result, err = o.shortener.ShortenAndSaveLink(req.URL)
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

	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read request body")
		return
	}

	o.shortener.UserID, err = c.Cookie("userID")
	//if err != nil {
	//	c.String(http.StatusUnauthorized, "Failed to get cookie - userID")
	//	return
	//}

	reqObj := strings.TrimSpace(string(body))

	outputURL, err := o.shortener.ShortenAndSaveLink(reqObj)
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

	URL, exist := o.shortener.GetOneByShortURL(shortURL)
	if !exist {
		c.String(http.StatusBadRequest, "Failed to get original URL")
		return
	}

	c.Header("Location", URL)
	c.Redirect(http.StatusTemporaryRedirect, URL)
}

func (o *Rest) ShortenURLs(c *gin.Context) {

	var URLs []*entities.ReqURL

	err := c.BindJSON(&URLs)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to read body",
			"error":   err.Error(),
		})
		return
	}

	o.shortener.UserID, err = c.Cookie("userID")
	//if err != nil {
	//	c.JSON(http.StatusUnauthorized, gin.H{
	//		"message": "Failed to get cookie - userID",
	//		"error":   err.Error(),
	//	})
	//	return
	//}

	shortenedURLs, err := o.shortener.ShortenURLs(URLs)
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

	o.shortener.UserID, err = c.Cookie("userID")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{
			"message": "Failed to get cookie - userID",
			"error":   err.Error(),
		})
		return
	}

	result, err := o.shortener.ListAll()
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to get cookie - userID",
			"error":   err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	if len(result) == 0 {
		c.JSON(http.StatusNoContent, result)
	} else {
		c.JSON(http.StatusCreated, result)
	}
}

func (o *Rest) Ping(c *gin.Context) {
	c.JSON(http.StatusOK, "")
}
