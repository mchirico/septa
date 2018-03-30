package main

import (
	septa "github.com/mchirico/septa/utils"
	"fmt"
)

func main() {

	url := "https://www3.septa.org/hackathon/" +
		"Arrivals/Suburban%20Station/25/"
	fmt.Println(septa.Parse(septa.GetData(url)))

}
