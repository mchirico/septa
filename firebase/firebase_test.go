package firebase

import (
	"fmt"
	septa "github.com/mchirico/septa/utils"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"strconv"
	"strings"
	"testing"
)

func TestTokenDirAndFile(t *testing.T) {

	before, _ := clientSecretFile()

	TokenDirAndFile = "/note/spud/.token.json"
	result, _ := clientSecretFile()
	assert.EqualValues(t, result, TokenDirAndFile)

	TokenDirAndFile = ""
	result, _ = clientSecretFile()
	assert.EqualValues(t, result, before)

}

func TestClientSecretFileNoToken(t *testing.T) {
	file, _ := clientSecretFile()
	assert.FileExistsf(t, file, "Token does not exit.")

}

func TestAddDeleteStation(t *testing.T) {
	station := "Elkins Park"
	//DeleteStation(station)
	AddStation(station)
}

// TestGetStationRecords -- live test and requires token.json
func TestGetStationRecords(t *testing.T) {
	station := "Elkins Park"
	expected := 1
	found := -1

	m := GetStationRecords(station, 1)
	tmp := strings.Split(m[0]["station"], ":")
	departureTime := tmp[1]

	fmt.Println("TMP", departureTime)
	fmt.Println("Live train_id:", m[0]["train_id"])

	i, err := strconv.Atoi(m[0]["train_id"])
	if err != nil {

	} else {
		if i > 100 {
			expected = i
			found = i
		}

	}
	assert.EqualValues(t, expected, found)
}

func TestGetStationRecordsWrapper(t *testing.T) {

	// TODO: Fix - Stations not all correct
	for i, station := range septa.ListStations() {
		//station := "Elkins Park"
		fmt.Printf("Station: %s\n", station)
		number := 3
		tmp := GetStationRecordsWrapper(station, number)
		assert.Contains(t, tmp[0]["station_rec_type"], station)
		assert.EqualValues(t, number*2, len(tmp), "Short on records")
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
	AllDocuments("trains", "train_id")
	AllDocuments("trainView", "TrainNo")
	//fmt.Println(m[0]["train_id"])
}

func TestDeleteDocument(t *testing.T) {
	DeleteDocument("trainView", "1507")
}

func TestRefreshLiveView(t *testing.T) {
	RefreshLiveView()
}

func TestGetAllStationsRecordsWrapper(t *testing.T) {

	success := false
	m := GetAllStationsRecordsWrapper(3)
	if len(m[0]["station_rec_type"]) > 3 {
		fmt.Printf("%s\n", m[0]["station_rec_type"])
		success = true
	}
	assert.Equal(t, true, success, "Network?")

}

func TestAddAllStations(t *testing.T) {

	//AddAllStations(3)
	fmt.Printf("This should only be run interactively")

}

func TestInsertUpdateDelete(t *testing.T) {

	insertUpdateDelete()

}
