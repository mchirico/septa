package firebase

import (
	"fmt"
	septa "github.com/mchirico/septa/utils"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"strings"
	"testing"
)

func TestAddDeleteStation(t *testing.T) {
	station := "Elkins Park"
	//DeleteStation(station)
	AddStation(station)
}

func TestGetStationRecords(t *testing.T) {
	station := "Elkins Park"

	m := GetStationRecords(station, 1)
	tmp := strings.Split(m[0]["station"], ":")
	departureTime := tmp[1]
	fmt.Println("TMP", departureTime)

	for k, v := range m[0] {
		fmt.Println(k, v)
	}
}
func TestGetStationRecordsWrapper(t *testing.T) {

	// TODO: Fix - Stations not all correct
	for i, station := range septa.ListStations() {
		//station := "Elkins Park"
		fmt.Printf("Station: %s\n", station)
		number := 3
		tmp := GetStationRecordsWrapper(station, number)
		assert.Contains(t, tmp[0]["station_rec_type"], station)
		assert.EqualValues(t, number, len(tmp), "Short on records")
		fmt.Printf(":%s", tmp[0]["time"])
		for _, rec := range tmp {
			fmt.Printf("\n:rec %v", rec)
		}

		if i >= 0 {
			break
		}

	}
}

func TestGetSingleDocument(t *testing.T) {
	SingleDocument()
}

func TestAllDocumentsSingleDocument(t *testing.T) {
	AllDocuments("trains")
	//fmt.Println(m[0]["train_id"])
}
