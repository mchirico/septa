package main

import (
	septa "github.com/mchirico/septa/utils"
	"fmt"
)

func main() {

	database := []map[string]string{}
	station := "Elkins Park"
	database = septa.GetStationRecords(station, 12, database)

	for _, v := range database {
		fmt.Printf("train_id:%4s, status:%11s, "+
			"sched_time: %6s\n",
			v["train_id"], v["status"], v["sched_time"])
	}

}
