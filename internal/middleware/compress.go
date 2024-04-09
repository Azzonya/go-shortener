// Package middleware provides HTTP middleware for the URL shortener application.
package middleware

import (
	"compress/gzip"
	"io"
	"log"
	"strings"

	"github.com/gin-gonic/gin"
)

// gzipWriter is a custom response writer that supports gzip compression.
type gzipWriter struct {
	gin.ResponseWriter
	Writer io.Writer
}

// Write writes the compressed data to the underlying writer.
func (w gzipWriter) Write(b []byte) (int, error) {
	return w.Writer.Write(b)
}

// CompressRequest is a Gin middleware function that compresses HTTP request bodies
// using gzip compression if the request's content type is "application/json" or "text/html"
// and the client supports gzip encoding.
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

// DecompressRequest is a Gin middleware function that decompresses HTTP request bodies
// encoded with gzip compression.
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
