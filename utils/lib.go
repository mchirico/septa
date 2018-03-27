package utils

import (
	"encoding/json"
	"fmt"
	"github.com/mchirico/bugGoProject/septa_fixtures"
	"reflect"
	"net/http"
	"log"
	"io/ioutil"
	"bytes"
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
					train_id,depart_time,status)
				buffer.WriteString(tmp_string+"\n")


			}
		}
	}
	return buffer.String()
}



func printStations(){
	for _, k := range septa_fixtures.StationList() {
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

func ListStations() {
	station_url := "http://www3.septa.org/hackathon/Arrivals/station_id_name.csv"
	b := GetData(station_url)

	fmt.Println(string(b))
}

