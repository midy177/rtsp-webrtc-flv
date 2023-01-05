package main

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/net/http2"
	"log"
	"rtsp-to-flv/controller"
	"time"
)

func main() {
	e := echo.New()
	e.Use(middleware.Recover())
	e.Use(middleware.CORS())
	e.Use(middleware.BodyLimit("50M"))
	e.GET("/flv", controller.HttpFlv)
	e.GET("/webrtc", controller.HttpWebrtc)
	e.Server.ReadTimeout = time.Second * 15
	e.Server.WriteTimeout = time.Minute * 36000
	e.ListenerNetwork = "tcp4"
	e.DisableHTTP2 = false
	h2s := &http2.Server{}
	err := e.StartH2CServer("0.0.0.0:44444", h2s)
	if err != nil {
		log.Println(err.Error())
	}
}
