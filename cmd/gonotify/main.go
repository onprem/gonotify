package main

import (
	"github.com/gin-gonic/gin"
	"github.com/prmsrswt/gonotify/pkg/api"
)

func main() {
	r := gin.Default()
	api.Register(r)

	r.Run()
}
