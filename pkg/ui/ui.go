package ui

import (
	"mime"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/log/level"
)

// BaseUI represent the React UI
type BaseUI struct {
	logger log.Logger
	router *gin.Engine
}

// NewBaseUI is used to create a instance of BaseUI
func NewBaseUI(router *gin.Engine, logger log.Logger) (*BaseUI, error) {
	return &BaseUI{
		logger: log.With(logger, "component", "ui"),
		router: router,
	}, nil
}

// Register registers all handlers for UI
func (bu *BaseUI) Register() {
	bu.router.NoRoute(bu.serveStatic())
}

func (bu *BaseUI) serveStatic() gin.HandlerFunc {
	return func(c *gin.Context) {
		method := c.Request.Method

		if method != "GET" {
			c.Next()
			return
		}

		path := c.Request.URL.Path
		path = strings.TrimPrefix(path, "/")

		data, err := Asset(path)
		if err != nil {
			level.Debug(bu.logger).Log("msg", "serving index", "path", path)
			path = "index.html"
			data, err = Asset("index.html")
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "some error occured"})
				level.Error(bu.logger).Log("err", err)
				return
			}
		}
		mimeType := getMime(path)

		c.Status(http.StatusOK)
		c.Writer.Header().Set("Content-Type", mimeType)
		c.Writer.Write(data)
	}
}

func getMime(path string) string {
	exts := strings.Split(path, ".")
	ext := "." + exts[len(exts)-1]

	mimeType := mime.TypeByExtension(ext)

	if mimeType == "" {
		return "text/plain"
	}
	return mimeType
}
