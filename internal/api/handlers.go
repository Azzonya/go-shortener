package api

import (
	"encoding/json"
	"fmt"
	"github.com/Azzonya/go-shortener/internal/util"
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

	shortURL := util.GenerateShortURL()

	err = o.storage.Add(shortURL, req.URL)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to add line to storage",
			"error":   err.Error(),
		})
		return
	}

	outputURL := fmt.Sprintf("%s/%s", o.baseURL, shortURL)

	resp.Result = outputURL

	resultJSON, err := json.Marshal(resp)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"message": "Failed to Marshall result struct",
			"error":   err.Error(),
		})
		return
	}

	c.Header("Content-Type", "application/json")
	c.Data(http.StatusCreated, "application/json", resultJSON)
}

func (o *Rest) Shorten(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to read request body")
		return
	}

	reqObj := strings.TrimSpace(string(body))

	shortURL := util.GenerateShortURL()

	err = o.storage.Add(shortURL, reqObj)
	if err != nil {
		c.String(http.StatusBadRequest, "Failed to add line to storage")
		return
	}

	outputURL := fmt.Sprintf("%s/%s", o.baseURL, shortURL)

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, outputURL)
}

func (o *Rest) Redirect(c *gin.Context) {
	shortURL, exist := c.Params.Get("id")
	if !exist {
		c.String(http.StatusBadRequest, "Failed to get ID")
		return
	}

	URL, exist := o.storage.GetOne(shortURL)
	if !exist {
		c.String(http.StatusBadRequest, "Failed to get original URL")
		return
	}

	c.Header("Location", URL)
	c.Redirect(http.StatusTemporaryRedirect, URL)
}
