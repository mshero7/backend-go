package main

import (
	"backend-go/server"
	"github.com/gin-gonic/gin"
	_ "github.com/lib/pq"
	"log"
	"net/http"
)

func ginServerLoad() {
	r := gin.Default()
	r.GET("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	r.POST("/", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	//listen and serve on 0.0.0.0:8080 (for windows "localhost:8080")
	err := r.Run()
	if err != nil {
		return
	}
}

func testServerLoad() {
	srv := server.NewHTTPServer(":8080")
	log.Fatal(srv.ListenAndServe())
}

func main() {
	ginServerLoad()
	testServerLoad() // for log prototype test
}
