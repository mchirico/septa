package utils

import (
	"github.com/stretchr/testify/assert"
	_ "github.com/stretchr/testify/mock"
	"testing"
)

func TestGetSurban(t *testing.T) {

	// fmt.Printf("%s", data)
	url := "https://www3.septa.org/hackathon/Arrivals/Suburban%20Station/25/"
	output := Parse(GetData(url))
	assert.Contains(t, output, "train: ")

}

func TestPrintStations(t *testing.T) {
	printStations()
	ListStations()
}
