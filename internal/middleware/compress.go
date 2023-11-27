package middleware

import (
	"compress/gzip"
	"github.com/gin-gonic/gin"
	"io"
	"log"
	"strings"
)

type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

func CompressRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		if !(c.GetHeader("Content-type") == "application/json" || c.GetHeader("Content-type") == "text/html") {
			c.Next()
			return
		}

		acceptEncoding := c.GetHeader("Accept-Encoding")
		supportsGzip := strings.Contains(acceptEncoding, "gzip")

		if supportsGzip {
			cw := gzip.NewWriter(c.Writer)
			defer cw.Close()
			c.Header("Content-Encoding", "gzip")
			c.Writer = &gzipWriter{c.Writer, cw}
		}

		c.Next()
	}
}

func DecompressRequest() gin.HandlerFunc {
	return func(c *gin.Context) {
		contentEncoding := c.GetHeader("Content-Encoding")
		supportsGzip := strings.Contains(contentEncoding, "gzip")

		if supportsGzip {
			cr, err := gzip.NewReader(c.Request.Body)
			if err != nil {
				log.Fatalf("gzip new reader error - %d", err)
				return
			}

			defer cr.Close()

			c.Request.Body = cr
		}

		c.Next()
	}
}
