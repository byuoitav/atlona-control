package main

import (
	"net/http"

	"github.com/byuoitav/clevertouch-control/device"
	"github.com/gin-gonic/gin"

	"github.com/spf13/pflag"
)

// func init() {
// 	fmt.Println("Main Package Init")
// }

func main() {
	var (
		port     string
		logLevel string
	)
	pflag.StringVarP(&port, "port", "p", "8040", "port for microservice to av-api communication")
	pflag.StringVarP(&logLevel, "log", "l", "Debug", "Initial log level") //Change debug to Info
	pflag.Parse()

	port = ":" + port

	manager := device.DeviceManager{
		Log: buildLogger(logLevel),
	}

	router := gin.Default()
	router.GET("/ping", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "pong",
		})
	})

	router.GET("/status", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"message": "good",
		})
	})

	manager.RunHTTPServer(router, port)
}
