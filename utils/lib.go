package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"reflect"
	"strings"
)

func Parse(b []byte) string {

	var data0 map[string]interface{}
	if err := json.Unmarshal(b, &data0); err != nil {
		panic(err)
	}

	var buffer bytes.Buffer
	for key, value := range data0 {
		fmt.Println("Key:", key, "Value:", value)
		for _, v := range value.([]interface{}) {
			records := v.(map[string]interface{})["Northbound"]
			if records == nil {
				continue
			}
			for _, rec := range records.([]interface{}) {
				train := rec.(map[string]interface{})
				train_id := train["train_id"]
				depart_time := train["depart_time"]
				status := train["status"]

				tmp_string := fmt.Sprintf("train: %s, depart: %s, status %s",
					train_id, depart_time, status)
				buffer.WriteString(tmp_string + "\n")

			}
		}
	}
	return buffer.String()
}

func printStations() {
	for _, k := range StationList() {
		fmt.Println(reflect.TypeOf(k))
		fmt.Println(k.Station)
	}
}

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

func ListStations() []string {
	station_url := "http://www3.septa.org/hackathon/Arrivals/station_id_name.csv"
	b := GetData(station_url)

	t_array := []string{}

	s := string(b)
	lines := strings.Split(s,"\n")
	if len(lines) < 2 {
		return t_array
	}

	for _,v := range lines[1:] {
		v := strings.Split(v,",")
		t_array = append(t_array, v[1])
	}
	//fmt.Println(string(b))
	return t_array

}

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

func GetStationRecords(
	station string, number int,
	database []map[string]string) []map[string]string {

	url := fmt.Sprintf("https://www3.septa.org/hackathon/Arrivals/%s/%d/",
		station, number)

	data := GetData(url)

	var data0 map[string]interface{}
	if err := json.Unmarshal(data, &data0); err != nil {
		log.Printf("Bad data GetStationRecords: %s",string(data))
		return nil
	}
	return GetParseMap(data, database)

}
