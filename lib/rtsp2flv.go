package lib

import (
	"errors"
	"github.com/deepch/vdk/format/flv"
	"net/http"
	"time"

	"github.com/deepch/vdk/format/rtspv2"
)

var (
	ErrorStreamExitNoVideoOnStream = errors.New("stream Exit No Video On Stream")
	ErrorStreamExitRtspDisconnect  = errors.New("stream Exit Rtsp Disconnect")
	ErrorStreamExitNoViewer        = errors.New("stream Exit On Demand No Viewer")
)

func RTSPToFlvWorker(url string, hasAudio, Debug bool, w http.ResponseWriter) error {
	keyTest := time.NewTimer(20 * time.Second)
	//add next TimeOut
	client, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: url, DisableAudio: !hasAudio, DialTimeout: 3 * time.Second, ReadWriteTimeout: 3 * time.Second, Debug: Debug})
	if err != nil {
		return err
	}
	defer client.Close()
	var AudioOnly bool
	if len(client.CodecData) == 1 && client.CodecData[0].Type().IsAudio() {
		AudioOnly = true
	}
	w.Header().Set("Content-Type", "video/x-flv")
	w.Header().Set("Transfer-Encoding", "chunked")
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.WriteHeader(200)
	fMux := flv.NewMuxer(w)
	if client.CodecData != nil {
		fMux.WriteHeader(client.CodecData)
	}
	for {
		select {
		case <-keyTest.C:
			return ErrorStreamExitNoVideoOnStream
		case signals := <-client.Signals:
			switch signals {
			case rtspv2.SignalCodecUpdate:
				//Config.coAd(name, client.CodecData)
			case rtspv2.SignalStreamRTPStop:
				return ErrorStreamExitRtspDisconnect
			}
		case packetAV := <-client.OutgoingPacketQueue:
			if AudioOnly || packetAV.IsKeyFrame {
				keyTest.Reset(20 * time.Second)
			}
			err = fMux.WritePacket(*packetAV)
			if err != nil {
				return err
			}
		}
	}
}
