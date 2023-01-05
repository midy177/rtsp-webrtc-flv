package lib

import (
	"github.com/deepch/vdk/format/rtspv2"
	webrtc "github.com/deepch/vdk/format/webrtcv3"
	"log"
	"time"
)

func RTSPToWebrtcWorker(url, spd string, hasAudio, Debug bool) (string, error) {
	client, err := rtspv2.Dial(rtspv2.RTSPClientOptions{URL: url, DisableAudio: !hasAudio, DialTimeout: 3 * time.Second, ReadWriteTimeout: 3 * time.Second, Debug: Debug})
	if err != nil {
		return "", err
	}
	muxerWebRTC := webrtc.NewMuxer(webrtc.Options{
		ICEServers:    []string{"stun:stun.l.google.com:19302"},
		ICEUsername:   "",
		ICECredential: "",
		PortMin:       0,
		PortMax:       0})
	answer, err := muxerWebRTC.WriteHeader(client.CodecData, spd)
	if err != nil {
		log.Println("WriteHeader", err)
		return "", err
	}
	go func(rtspCli *rtspv2.RTSPClient, rtcCli *webrtc.Muxer) {
		defer rtspCli.Close()
		var AudioOnly bool
		if len(client.CodecData) == 1 && client.CodecData[0].Type().IsAudio() {
			AudioOnly = true
		}
		keyTest := time.NewTimer(20 * time.Second)
	startWebrtc:
		for {
			select {
			case <-keyTest.C:
				log.Printf("SignalStreamRTPStop -> %s", ErrorStreamExitNoVideoOnStream)
				break startWebrtc
			case signals := <-rtspCli.Signals:
				switch signals {
				case rtspv2.SignalCodecUpdate:
					//Config.coAd(name, client.CodecData)
				case rtspv2.SignalStreamRTPStop:
					log.Printf("SignalStreamRTPStop -> %s", ErrorStreamExitRtspDisconnect)
					break startWebrtc
				}
			case packetAV := <-rtspCli.OutgoingPacketQueue:
				if AudioOnly || packetAV.IsKeyFrame {
					keyTest.Reset(20 * time.Second)
				}
				err = rtcCli.WritePacket(*packetAV)
				if err != nil {
					log.Printf("writePacket -> %s", err)
					break startWebrtc
				}
			}
		}
	}(client, muxerWebRTC)
	return answer, err
}
