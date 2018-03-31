package main

import (
	"github.com/mchirico/septa/firebase"
	"time"
)

func init() {
	firebase.Flags()
}

func main() {

	stations := []string{
		"Elkins Park",
		"30th Street Station",
		"Suburban Station",
		"Airport Terminal A"}

	for {

		for _, station := range stations {
			firebase.AddStation(station)
			firebase.AddStation(station)
		}
		firebase.RefreshLiveView()

		time.Sleep(time.Duration(firebase.QueryTime) *
			1000 * time.Millisecond)
	}

}
