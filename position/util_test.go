package position

import (
	"flag"
	"os"
	"testing"
	"github.com/mchirico/septa/utils"
	"fmt"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestFlags(t *testing.T) {
	m, _ := utils.GetLiveViewRecords()

	t.Logf("Test: %v\n",m[0].TrainNo)

	for i,rec := range m {

		if i==2 {
			fmt.Printf("train:",rec)
		}

	}

}