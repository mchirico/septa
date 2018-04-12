package main

import (
	"github.com/mchirico/septa/firebase"
	"time"
)

func init() {
	firebase.Flags()
}

// allstationsByTime -- separate thread. This takes a long time.
func allstationsByTime() {
	for {

		firebase.AllStationsByTime()

		time.Sleep(2 *
			1000 * time.Millisecond)
	}
}

func rrSchedules() {
	for {
		firebase.AddRRSchedules()
		time.Sleep(time.Duration(firebase.QueryTime) *
			1000 * time.Millisecond)

	}
}

func main() {

	go rrSchedules()
	go allstationsByTime()
	for {

		firebase.RefreshLiveView()

		time.Sleep(time.Duration(firebase.QueryTime) *
			1000 * time.Millisecond)
	}

}
