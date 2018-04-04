package utils

import (
	"testing"
	"time"
	"fmt"
)


func testSetup() []map[string]interface{}{

	var testVariables []map[string]interface{}

	loc, _ := time.LoadLocation("America/New_York")
	tn := time.Now().In(loc)

	trainRRSchedules := TrainRRSchedules{ TrainID:"TrainID",
	DocDate:"DocDate",Timestamp: tn,
	RRSchedules:nil}

	fmt.Printf("%v\n",
		trainRRSchedules.TrainID)


	return testVariables



}


func TestDefinitions(t *testing.T) {

	testSetup()

}