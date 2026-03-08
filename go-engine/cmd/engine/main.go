package main

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	//db := db.ConectDb()

	route := gin.Default()

	route.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"mensagem": "pong"})
	})

}
