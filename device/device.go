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

	//AT-UHD-SW-52ED - should work with all Atlona 1 output video switchers (do not use volume if the switcher does not have volume control)
	// action endpoints
	at52 := router.Group("/api/v1/AT-UHD-SW-52ED")
	at52.GET("/:address/output/:output/input/:input", d.setInput)  //change input
	at52.GET("/:address/block/:output/volume/:level", d.setVolume) //set volume
	at52.GET("/:address/block/:output/muted/true", d.setMute)      //set mute true
	at52.GET("/:address/block/:output/muted/false", d.setUnMute)   //set mute false

	// status endpoints
	at52.GET("/:address/status/output/:output/input", d.getInput)  //get input
	at52.GET("/:address/status/block/:output/volume", d.getVolume) //get volume
	at52.GET("/:address/status/block/:output/muted", d.getMute)    //get mute state

	//AT-OME-PS62 - should work with all Atlona multi-output video switchers
	// action endpoints
	at62 := router.Group("/api/v1/AT-OME-PS62")
	at62.GET("/:address/output/:output/input/:input", d.setInput)  //change input
	at62.GET("/:address/block/:output/volume/:level", d.setVolume) //set volume
	at62.GET("/:address/block/:output/muted/true", d.setMute)      //set mute true
	at62.GET("/:address/block/:output/muted/false", d.setUnMute)   //set mute false

	// status endpoints
	at62.GET("/:address/status/output/:output/input", d.getInput)  //get input
	at62.GET("/:address/status/block/:output/volume", d.getVolume) //get volume
	at62.GET("/:address/status/block/:output/muted", d.getMute)    //get mute state

	//AT-GAIN-60 - should work with all Atlona multi-output video switchers
	// action endpoints
	atGain60 := router.Group("/api/v1/AT-GAIN-60")
	//atGain60.GET("/:address/output/:output/input/:input", d.changeInput) //change input
	atGain60.GET("/:address/block/:output/volume/:level", d.setVolumeGain60) //set volume
	atGain60.GET("/:address/block/:output/muted/true", d.setMuteGain60)      //set mute true
	atGain60.GET("/:address/block/:output/muted/false", d.setUnMuteGain60)   //set mute false

	// status endpoints
	//atGain60.GET("/:address/output/:port/input", d.getInput)   //get input
	atGain60.GET("/:address/status/block/:input/volume", d.getVolumeGain60) //get volume
	atGain60.GET("/:address/status/block/:input/muted", d.getMuteGain60)    //get mute state

	server := &http.Server{
		Addr:           port,
		MaxHeaderBytes: 1021 * 10,
	}

	d.Log.Info("running http server")
	router.Run(server.Addr)

	return fmt.Errorf("http server stopped")
}
