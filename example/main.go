package main

import (
	"fmt"
	"net/http"
	"os"
	"os/signal"

	"github.com/gin-gonic/gin"
)

func fn(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{
		"message": "pong",
	})
}
func fn2() {
	fmt.Println("aaaaa")
}
func main() {
	fn2()
	r := gin.Default()
	r.GET("/ping", fn)
	go r.Run() // listen and serve on http://localhost:8080/ping
	sig := make(chan os.Signal, 1)
	signal.Notify(sig, os.Interrupt, os.Kill)

	select {
	case <-sig:
		fmt.Println("exit!")
	}
}
