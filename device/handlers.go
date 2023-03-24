package device

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/byuoitav/clevertouch-control/device/actions"
	"go.uber.org/zap"

	"github.com/gin-gonic/gin"
)

// communication needs to be buffered/slowed down in this package
var takeANumber int

const communicationFrequency int = 1300 //time in ms for polling

func init() {
	go func() {
		for {
			if takeANumber > 0 {
				time.Sleep(time.Duration(communicationFrequency) * time.Millisecond)
				takeANumber -= 1

			}
		}
	}()
}

func que() {
	takeANumber += 1
	totalWaitTime := time.Duration(communicationFrequency * (takeANumber - 1))
	time.Sleep(totalWaitTime * time.Millisecond)
}

func getDeviceType(context *gin.Context) (device string) {
	path := context.FullPath()
	paths := strings.Split(path, "/")
	device = paths[3]

	return device
}

func (d *DeviceManager) setInput(context *gin.Context) {
	que()

	d.Log.Debug("setting input", zap.String("input", context.Param("input")), zap.String("output", context.Param("output")), zap.String("address", context.Param("address")))

	input, err := actions.SetInput(context.Param("address"), context.Param("input"), context.Param("output"))
	if err != nil {
		d.Log.Warn("failed to set input", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	d.Log.Debug("successfully set input", zap.String("input", context.Param("input")), zap.String("output", context.Param("output")), zap.String("address", context.Param("address")))
	context.JSON(http.StatusOK, input)
}

func (d *DeviceManager) getInput(context *gin.Context) {
	que()

	d.Log.Debug("getting input status", zap.String("address", context.Param("address")), zap.String("output", context.Param("output")))

	input, err := actions.GetInput(context.Param("address"), context.Param("output"))
	if err != nil {
		d.Log.Warn("failed to get input status", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	d.Log.Debug("received input status", zap.String("input", input.Input), zap.String("address", context.Param("address")))
	context.JSON(http.StatusOK, input)
}

func (d *DeviceManager) setMute(context *gin.Context) {
	que()
	device := getDeviceType(context)

	d.Log.Debug("setting mute", zap.String("mute", context.Param("mute")), zap.String("address", context.Param("address")), zap.String("device type", device))

	mute, err := strconv.ParseBool(context.Param("mute"))
	if err != nil {
		d.Log.Warn("could not set mute. 'mute' parameter not a valid boolean value", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}
	state, err := actions.SetMute(context.Param("address"), context.Param("output"), mute, device)
	if err != nil {
		d.Log.Warn("failed to set mute", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	d.Log.Debug("successfully set mute", zap.String("mute", strconv.FormatBool(mute)), zap.String("address", context.Param("addres")))
	context.JSON(http.StatusOK, state)
}

func (d *DeviceManager) getMute(context *gin.Context) {
	que()
	device := getDeviceType(context)

	d.Log.Debug("getting mute status", zap.String("address", context.Param("address")))

	state, err := actions.GetMute(context.Param("address"), context.Param("output"), device)
	if err != nil {
		d.Log.Warn("failed to get mute status", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	d.Log.Debug("received mute status", zap.String("mute", strconv.FormatBool(state.Muted)))
	context.JSON(http.StatusOK, state)
}

func (d *DeviceManager) setVolume(context *gin.Context) {
	que()
	device := getDeviceType(context)

	d.Log.Debug("setting volume", zap.String("level", context.Param("level")), zap.String("address", context.Param("address")), zap.String("device type", device))

	volume, err := actions.SetVolume(context.Param("address"), context.Param("output"), context.Param("level"), device)
	if err != nil {
		d.Log.Warn("failed to set volume", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	d.Log.Debug("successfully set volume", zap.String("output", context.Param("output")), zap.String("level", context.Param("level")), zap.String("address", context.Param("address")))
	context.JSON(http.StatusOK, volume)
}

func (d *DeviceManager) getVolume(context *gin.Context) {
	que()
	device := getDeviceType(context)

	d.Log.Debug("getting volume", zap.String("address", context.Param("address")), zap.String("device type", device))

	volume, err := actions.GetVolume(context.Param("address"), context.Param("output"), device)
	if err != nil {
		d.Log.Warn("failed to get volume", zap.Error(err))
		context.JSON(http.StatusInternalServerError, err.Error())
		return
	}

	d.Log.Debug("received volume status", zap.String("volume", strconv.Itoa(volume.Volume)), zap.String("address", context.Param("address")))
	context.JSON(http.StatusOK, volume)
}
