package main

import (
	"github.com/mchirico/septa/firebase"
	septa "github.com/mchirico/septa/utils"
	"time"
)

func init() {
	firebase.Flags()
}

// allstationsByTime -- separate thread. This takes a long time.
func allstationsByTime() {
	for {

		firebase.AllStationsByTime()
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

func allStations() {

	records := septa.GetLiveViewRecords()

	for _, train := range records {
		go firebase.AddStationsByTime(train.TrainNo)
	}

}

func allStationsWrapper() {
	for {
		allStations()
		time.Sleep(time.Duration(firebase.QueryTime) *
			1000 * time.Millisecond)

	}
}

func main() {

	go allStationsWrapper()
	go rrSchedules()
	go allstationsByTime()
	for {

		firebase.RefreshLiveView()

		time.Sleep(time.Duration(firebase.QueryTime) *
			1000 * time.Millisecond)
	}

}
