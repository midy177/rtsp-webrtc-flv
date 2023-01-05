package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rtsp-to-flv/echox"
	"rtsp-to-flv/lib"
)

type HttpWebrtcParams struct {
	RtspUrl  string `json:"rtsp_url" query:"rtsp_url"`
	SPD      string `json:"spd" query:"spd"`
	HasAudio bool   `json:"has_audio" query:"has_audio"`
	Debug    bool   `json:"debug" query:"debug"`
}

func HttpWebrtc(ctx echo.Context) error {
	config := new(HttpWebrtcParams)
	if err := ctx.Bind(config); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err.Error()}.JSON(ctx)
	}
	data, err := lib.RTSPToWebrtcWorker(config.RtspUrl, config.SPD, config.HasAudio, config.Debug)
	if err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err.Error()}.JSON(ctx)
	}
	return echox.Response{Code: http.StatusOK, Data: data}.JSON(ctx)
}
