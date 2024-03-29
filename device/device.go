package device

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

type DeviceManager struct {
	Log *zap.Logger
}

func (d *DeviceManager) RunHTTPServer(router *gin.Engine, port string) error {
	d.Log.Info("registering http endpoints")

	//AT-UHD-SW-52ED
	// action endpoints
	at52 := router.Group("/api/v1/AT-UHD-SW-52ED")
	at52.GET("/:address/output/:output/input/:input", d.setInput)  //change input
	at52.GET("/:address/block/:output/volume/:level", d.setVolume) //set volume
	at52.GET("/:address/block/:output/muted/:mute", d.setMute)     //set mute true/false

	// status endpoints
	at52.GET("/:address/output/:output/input", d.getInput)  //get input
	at52.GET("/:address/block/:output/volume", d.getVolume) //get volume
	at52.GET("/:address/block/:output/muted", d.getMute)    //get mute state

	//AT-OME-PS62
	// action endpoints
	at62 := router.Group("/api/v1/AT-OME-PS62")
	at62.GET("/:address/output/:output/input/:input", d.setInput)  //change input
	at62.GET("/:address/block/:output/volume/:level", d.setVolume) //set volume
	at62.GET("/:address/block/:output/muted/:mute", d.setMute)     //set mute true/false

	// status endpoints
	at62.GET("/:address/output/:output/input", d.getInput)  //get input
	at62.GET("/:address/block/:output/volume", d.getVolume) //get volume
	at62.GET("/:address/block/:output/muted", d.getMute)    //get mute state

	//AT-GAIN-60
	// action endpoints
	atGain60 := router.Group("/api/v1/AT-GAIN-60")
	atGain60.GET("/:address/block/:output/volume/:level", d.setVolume) //set volume
	atGain60.GET("/:address/block/:output/muted/:mute", d.setMute)     //set mute true/false

	// status endpoints
	atGain60.GET("/:address/block/:output/volume", d.getVolume) //get volume
	atGain60.GET("/:address/block/:output/muted", d.getMute)    //get mute state

	server := &http.Server{
		Addr:           port,
		MaxHeaderBytes: 1021 * 10,
	}

	d.Log.Info("running http server")
	router.Run(server.Addr)

	return fmt.Errorf("http server stopped")
}
