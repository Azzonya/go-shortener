package service

import (
	"fmt"
	"github.com/Azzonya/go-shortener/internal/util"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

func (o *Rest) Shorten(c *gin.Context) {
	fmt.Println(1)
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.String(http.StatusBadRequest, "Не удалось прочитать тело запроса")
		return
	}

	reqObj := strings.TrimSpace(string(body))

	shortURL := util.GenerateShortURL()

	o.storage.Add(shortURL, reqObj)

	outputURL := fmt.Sprintf("%s/%s", o.baseURL, shortURL)

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, outputURL)
}

func (o *Rest) Redirect(c *gin.Context) {
	shortURL, exist := c.Params.Get("id")
	if !exist {
		c.String(http.StatusBadRequest, "Не удалось получить ID")
		return
	}

	URL, exist := o.storage.GetOne(shortURL)
	if !exist {
		c.String(http.StatusBadRequest, "Не удалось получить оригинальную ссылку")
		return
	}

	c.Header("Location", URL)
	c.Redirect(http.StatusTemporaryRedirect, URL)
}
