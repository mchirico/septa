package main

import (
	"github.com/mchirico/septa/firebase"
	"time"
)

func main() {

	firebase.Flags()
	station := "Elkins Park"

	for {
		firebase.AddStation(station)
		firebase.RefreshLiveView()
		time.Sleep(time.Duration(firebase.QueryTime) *
			1000 * time.Millisecond)
	}

}
