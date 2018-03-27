package utils

import (
	_ "github.com/stretchr/testify/mock"
	"testing"
	"github.com/stretchr/testify/assert"
)


func TestGetSurban(t *testing.T) {

	// fmt.Printf("%s", data)
	url := "https://www3.septa.org/hackathon/Arrivals/Suburban%20Station/25/"
	output:=Parse(GetData(url))
	assert.Contains(t,output,"train: ")



}

