package utils

import (
	"encoding/json"
	"github.com/stretchr/testify/assert"
	"io/ioutil"
	"log"
	"strings"
	"testing"
)

func fixtureReader(file string) *json.Decoder {

	dat, err := ioutil.ReadFile(file)
	if err != nil {
		log.Fatal(err)
	}
	jsonStream := string(dat)

	return json.NewDecoder(strings.NewReader(jsonStream))

}

func TestDefinitions(t *testing.T) {

	fixtureReader("../fixtures/trainView.json")
	status := []bool{false, false, false}

	dat, err := ioutil.ReadFile("../fixtures/trainView.json")
	assert.Nil(t, err)
	jsonStream := string(dat)

	dec := json.NewDecoder(strings.NewReader(jsonStream))
	print(dec)

	_, err = dec.Token()
	assert.Nil(t, err)

	for dec.More() {
		var m LiveViewMessage
		err := dec.Decode(&m)
		assert.Nil(t, err)
		if m.Lat == "39.959606265" {
			status[0] = true
		}
		if m.Lon == "-75.063539333333" {
			status[1] = true
		}
		if m.Source == "Marcus Hook" {
			status[2] = true
		}

	}
	for _, v := range status {
		assert.True(t, v)
	}

	_, err = dec.Token()
	assert.Nil(t, err)

}

func TestRRSchedulesType(t *testing.T) {

}
