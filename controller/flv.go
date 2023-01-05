package controller

import (
	"github.com/labstack/echo/v4"
	"net/http"
	"rtsp-to-flv/echox"
	"rtsp-to-flv/lib"
)

type HttpFlvParams struct {
	RtspUrl  string `json:"rtsp_url" query:"rtsp_url"`
	HasAudio bool   `json:"has_audio" query:"has_audio"`
	Debug    bool   `json:"debug" query:"debug"`
}

func HttpFlv(ctx echo.Context) error {
	config := new(HttpFlvParams)
	if err := ctx.Bind(config); err != nil {
		return echox.Response{Code: http.StatusBadRequest, Message: err.Error()}.JSON(ctx)
	}
	return lib.RTSPToFlvWorker(config.RtspUrl, config.HasAudio, config.Debug, ctx.Response())
}
