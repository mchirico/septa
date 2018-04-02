package main

import (
	"github.com/mchirico/septa/firebase"
	"time"
)

func init() {
	firebase.Flags()
}

// allstations -- separate thread. This takes a long time.
func allstations(number int) {
	for {

		firebase.AddAllStations(number)
		time.Sleep(time.Duration(firebase.QueryTime) *
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
	go allstations(3)
	for {

		firebase.RefreshLiveView()

		time.Sleep(time.Duration(firebase.QueryTime) *
			1000 * time.Millisecond)
	}

}
