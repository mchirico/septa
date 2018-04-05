package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
	"time"
)

// GetData --  basic http GET returning raw data
func GetData(url string) []byte {
	// http://www3.septa.org/hackathon/TrainView/
	res, err := http.Get(url)
	if err != nil {
		log.Fatal(err)
	}
	data, err := ioutil.ReadAll(res.Body)
	res.Body.Close()
	if err != nil {
		log.Fatal(err)
	}
	return data

}

// ListStations -- list all stations
func ListStations() []string {
	stationURL := "http://www3.septa.org/hackathon/Arrivals/station_id_name.csv"

	b := GetData(stationURL)
	tArray := []string{}

	s := string(b)
	lines := strings.Split(s, "\n")
	if len(lines) < 2 {
		return tArray
	}

	for _, v := range lines[1:] {
		v := strings.Split(v, ",")
		tArray = append(tArray, v[1])
	}
	//fmt.Println(string(b))
	return tArray

}

// insertIndividualRecord -- local to GetParseMap
func insertIndividualRecord(station string,
	database []map[string]string,
	records interface{},
) []map[string]string {

	for _, rec := range records.([]interface{}) {
		m := rec.(map[string]interface{})
		recordmap := map[string]string{}
		recordmap["station"] = station
		for k, v := range m {
			if v != nil {
				recordmap[string(k)] = v.(string)
			}
		}
		database = append(database, recordmap)
	}
	return database
}

// GetParseMap -- specific to station arrivals
func GetParseMap(
	b []byte,
	database []map[string]string) []map[string]string {

	var data0 map[string]interface{}
	if err := json.Unmarshal(b, &data0); err != nil {
		panic(err)
	}

	for key, value := range data0 {
		valueType := fmt.Sprintf("%v", reflect.TypeOf(value))
		if valueType != "[]interface {}" {
			continue
		}
		iValue := value.([]interface{})
		if iValue == nil {
			continue
		}
		for _, v := range iValue {
			station := key
			rType := fmt.Sprintf("%v", reflect.TypeOf(v))
			if rType != "map[string]interface {}" {
				return database
			}

			records := v.(map[string]interface{})["Northbound"]

			if records != nil {
				database = insertIndividualRecord(station,
					database,
					records)
			}
			records = v.(map[string]interface{})["Southbound"]
			if records != nil {
				database = insertIndividualRecord(station,
					database, records)
			}

		}
	}
	return database
}

// GetStationRecords -- gets records from arrivals
func GetStationRecords(
	station string, number int,
	database []map[string]string) []map[string]string {

	url := fmt.Sprintf("https://www3.septa.org/hackathon/Arrivals/%s/%d/",
		station, number)

	data := GetData(url)

	var data0 map[string]interface{}
	if err := json.Unmarshal(data, &data0); err != nil {
		log.Printf("Bad data GetStationRecords: %s", string(data))
		return nil
	}
	return GetParseMap(data, database)

}

// GetAllStationRecords -- this gets everything
func GetAllStationRecords(number int) []map[string]string {
	database := []map[string]string{}
	for _, station := range ListStations() {
		database = GetStationRecords(station, number, database)
		if database == nil {
			fmt.Printf("Bad station:%s", station)
		}

	}
	return database
}

// ParseLiveView -- specific to live view
//    Reference: https://play.golang.org/p/XxsmA8a7YPj
func ParseLiveView(jsonStream []byte) []LiveViewMessage {

	dec := json.NewDecoder(strings.NewReader(string(jsonStream)))

	database := []LiveViewMessage{}

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	for dec.More() {
		var m LiveViewMessage
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		database = append(database, m)
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	return database
}

// GetLiveViewRecords -- gets live view
func GetLiveViewRecords() []LiveViewMessage {

	url := "http://www3.septa.org/hackathon/TrainView/"
	jsonStream := GetData(url)
	return ParseLiveView(jsonStream)

}

// GetRRSchedules --
func GetRRSchedules(train string) TrainRRSchedules {

	url := fmt.Sprintf(
		"http://www3.septa.org/hackathon/RRSchedules/%s", train)
	jsonStream := GetData(url)
	return ParseRRSchedules(jsonStream, train)

}

// GetAlerts - url for get alerts
func GetAlerts() string {

	url := "http://www3.septa.org/hackathon/Alerts/"
	return url
}

// ParseRRSchedules -- schedule times of trains
func ParseRRSchedules(jsonStream []byte, train string) TrainRRSchedules {

	dec := json.NewDecoder(strings.NewReader(string(jsonStream)))

	trainRRSchedules := TrainRRSchedules{}

	records := []RRSchedules{}

	// read open bracket
	_, err := dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	for dec.More() {
		var m RRSchedules
		err := dec.Decode(&m)
		if err != nil {
			log.Fatal(err)
		}
		records = append(records, m)
	}

	// read closing bracket
	_, err = dec.Token()
	if err != nil {
		log.Fatal(err)
	}

	loc, _ := time.LoadLocation("America/New_York")
	t := time.Now().In(loc)
	trainRRSchedules.RRSchedules = records
	trainRRSchedules.TrainID = train
	trainRRSchedules.Timestamp = t

	trainRRSchedules.DocDate = t.Format("2006-01-02")

	return trainRRSchedules
}
