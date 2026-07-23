package main

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()
	r.GET("/health", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "ok", "message": "qeko api is running successfully"})
	})

	err := http.ListenAndServe(":8080", r)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Server is running on port 8080")
	}
}
