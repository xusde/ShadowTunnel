package speed

import (
	"time"
)

var UploadChannel chan int64 = make(chan int64, 20)
var DownloadChannel chan int64 = make(chan int64, 20)
var uploadSpeed, downloadSpeed float64 = 0, 0
var uploadTraffic, downloadTraffic int64 = 0, 0

func updateTraffic() {
	for {
		select {
		case u := <-UploadChannel:
			uploadTraffic += u
		case d := <-DownloadChannel:
			downloadTraffic += d
		}
	}
}

func updateSpeed(span time.Duration) {
	for {
		uploadTrafficT0 := uploadTraffic
		downloadTrafficT0 := downloadTraffic
		time.Sleep(span)
		// speed = (trafficT1 - trafficT0) / span
		// speed is in kilobytes per second
		uploadSpeed = float64(uploadTraffic-uploadTrafficT0) / float64(span/time.Second) / float64(1024)
		downloadSpeed = float64(downloadTraffic-downloadTrafficT0) / float64(span/time.Second) / float64(1024)
	}
}

func GetSpeed() (float64, float64) {
	return uploadSpeed, downloadSpeed
}

func GetTotalTraffic() (int64, int64) {
	return uploadTraffic, downloadTraffic
}

func StartSpeedMonitor(span time.Duration) {
	go updateTraffic()
	go updateSpeed(span)
}
