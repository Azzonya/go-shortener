package server

import (
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"strings"
)

//func (o *Rest) HShortenerURL(w http.ResponseWriter, r *http.Request) {
//	if r.Method == http.MethodPost {
//		o.HShortener(w, r)
//	} else if r.Method == http.MethodGet {
//		o.HRedirect(w, r)
//	}
//}

func (o *Rest) HShortener(c *gin.Context) {
	body, err := io.ReadAll(c.Request.Body)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось прочитать тело запроса"})
		return
	}

	reqObj := strings.TrimSpace(string(body))

	shortURL := o.generateShortURL()

	o.urlMap[shortURL] = reqObj

	c.Header("Content-Type", "text/plain")
	c.String(http.StatusCreated, "http://localhost:8080/"+shortURL)
}

func (o *Rest) HRedirect(c *gin.Context) {
	shortURL, exist := c.Params.Get("id")
	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить ID"})
		return
	}

	URL, exist := o.urlMap[shortURL]

	if !exist {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Не удалось получить оригинальную ссылку"})
		return
	}

	//w.Header().Set("Content-Type", "text/plain")
	c.Header("Location", URL)
	c.Redirect(http.StatusTemporaryRedirect, URL)
}
