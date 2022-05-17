package main

import (
	"fmt"
	"github.com/donetkit/contrib/discovery"
	"github.com/donetkit/contrib/discovery/consul"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"time"
)

func main() {
	r := gin.New()
	consulClient, _ := consul.New(
		discovery.WithServiceRegisterAddr("127.0.0.1"),
		discovery.WithServiceRegisterPort(8500),
		discovery.WithCheckHTTP(func(url string) { r.GET(url, func(c *gin.Context) { c.String(200, "Healthy") }) }))
	// Example ping request.
	r.GET("/ping", func(c *gin.Context) {
		c.String(http.StatusOK, "pong "+fmt.Sprint(time.Now().Unix()))
	})
	r.GET("/register", func(c *gin.Context) {
		consulClient.Register()
		c.String(http.StatusOK, "ok")
	})
	r.GET("/deregister", func(c *gin.Context) {
		consulClient.Deregister()
		c.String(http.StatusOK, "ok")
	})
	// Listen and Server in 0.0.0.0:8080
	if err := r.Run(":8080"); err != nil {
		log.Fatal(err)
	}
}
