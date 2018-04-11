package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"os"
	"strconv"
	"testing"
	"time"
)

func TestListtations(t *testing.T) {

	tmp := ListStations()
	// Number of stations should be at least 19
	success := false
	if len(tmp[0]) >= 19 {
		success = true
	}
	assert.Equal(t, true, success, "Stations not"+
		" found, can you access the web?"+
		" http://www3.septa.org/hackathon/Arrivals/station_id_name.csv")

}

func TestGetParseMap(t *testing.T) {
	var database []map[string]string
	m := GetParseMap(testGetParseMapData, database)
	assert.EqualValues(t, 6, len(m))
	assert.EqualValues(t, "436", m[0]["train_id"])
	assert.EqualValues(t, "438", m[1]["train_id"])
	assert.EqualValues(t, "440", m[2]["train_id"])

}

func TestGetStationRecords(t *testing.T) {

	station := "Elkins Park"

	database := []map[string]string{}

	GetAllStationRecords(3)

	m := GetStationRecords(station, 3, database)
	for _, v := range m {
		fmt.Println(v["train_id"])
	}
	fmt.Println(m[0]["station"])
}

func TestGetLiveView(t *testing.T) {

	m := GetLiveViewRecords()

	flag := false
	for _, v := range m {
		lat, err := strconv.ParseFloat(v.Lat, 8)
		if err != nil {
			flag = false
			fmt.Printf("***** Maybe Network Issue ********")
			fmt.Printf("%v %v\n", v, err)
			break
		}
		if lat > 30 {
			flag = true
		}
	}
	assert.EqualValues(t, true, flag, "Are trains running?")

}

func TestGetRRSchedules(t *testing.T) {

	m := GetLiveViewRecords()
	r := GetRRSchedules(m[0].TrainNo)
	fmt.Printf("%v\n", r.TrainID)
	for _, v := range r.RRSchedules {
		fmt.Printf("%v\n", v.Station)
	}
	loc, _ := time.LoadLocation("America/New_York")
	assert.Equal(t, r.DocDate, time.Now().In(loc).Format("2006-01-02"))

}

// You should not have to run this. This was used to create
// Fixture data.
func createFixtureData() {
	url := "http://www3.septa.org/hackathon/TrainView/"
	jsonStream := GetData(url)
	f, err := os.Create("../fixtures/FullTrainView.json")
	if err != nil {
		fmt.Printf("error opening file")
	}
	f.Write(jsonStream)
	f.Close()

}

func getSampleFullTrainView() []byte {
	f, err := os.Open("../fixtures/FullTrainView.json")
	if err != nil {
		fmt.Printf("error opening file")
	}
	data := make([]byte, 105210)
	n, err := f.Read(data)
	fmt.Printf("number of bytes read: %v\n", n)

	return data
}

func TestFixtureData(t *testing.T) {

	database := ParseLiveView(getSampleFullTrainView())

	assert.Equal(t, "39.95514", database[0].Lat)
	assert.Equal(t, 7, database[34].Late)

}

func BenchmarkStuff(b *testing.B) {

	jsonData := getSampleFullTrainView()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		ParseLiveView(jsonData)
	}
	b.StopTimer()
	b.ReportAllocs()
	fmt.Printf("N:%v\n", b.N)

}
