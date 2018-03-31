package utils

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"testing"
)

func TestListtations(t *testing.T) {
	tmp := ListStations()
	println(tmp[0])

}

func TestGetParseMap(t *testing.T) {
	var database []map[string]string
	m := GetParseMap(testGetParseMapData, database)
	assert.EqualValues(t, 3, len(m))
	assert.EqualValues(t, "436", m[0]["train_id"])
	assert.EqualValues(t, "438", m[1]["train_id"])
	assert.EqualValues(t, "440", m[2]["train_id"])

}

func TestGetStationRecords(t *testing.T) {

	station := "Elkins Park"

	database := []map[string]string{}

	m := GetStationRecords(station, 3, database)
	for _, v := range m {
		fmt.Println(v["train_id"])
	}
	fmt.Println(m[0]["station"])
}
