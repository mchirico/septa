package main

import (
	"github.com/mchirico/septa/firebase"
)

func main() {

	station := "Elkins Park"

	firebase.AddStation(station)

}
