package utils

import (
	"fmt"
	"reflect"
	"net/http"
	"log"
	"io/ioutil"
	"strings"
	"encoding/json"
)

func printStations() {
	for _, k := range StationList() {
		fmt.Println(reflect.TypeOf(k))
		fmt.Println(k.Station)
	}
}


// GetData --  basic http GET
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


// GetParseMap -- this comment required
func GetParseMap(
	b []byte,
	database []map[string]string) []map[string]string {

	var data0 map[string]interface{}
	if err := json.Unmarshal(b, &data0); err != nil {
		panic(err)
	}

	for key, value := range data0 {
		for _, v := range value.([]interface{}) {
			station := key
			records := v.(map[string]interface{})["Northbound"]
			if records == nil {
				continue
			}
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

		}
	}
	return database
}


// GetStationsRecords -- this will probably have to run this throug
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
