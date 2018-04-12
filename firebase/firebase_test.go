package firebase

import (
	"flag"
	"fmt"
	septa "github.com/mchirico/septa/utils"
	"github.com/stretchr/testify/assert"
	"google.golang.org/api/iterator"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"testing"
)

func TestMain(m *testing.M) {
	flag.Parse()
	os.Exit(m.Run())
}

func TestFlags(t *testing.T) {

	Flags()

	quiet := flag.Lookup("quiet").Value.(flag.Getter).Get().(bool)
	assert.False(t, quiet)

	time := flag.Lookup("time").Value.(flag.Getter).Get().(int)
	assert.Equal(t, 20, time)

	// Example of set
	flag.Set("time", "3")
	time = flag.Lookup("time").Value.(flag.Getter).Get().(int)
	assert.Equal(t, 3, time)

	token := flag.Lookup("token").Value.(flag.Getter).Get().(string)
	assert.Equal(t, "", token)

	fmt.Printf("time: %d", time)

}

// Very cool... see
//  https://talks.golang.org/2014/testing.slide#23
func TestFatalExit(t *testing.T) {
	if os.Getenv("BE_CRASHER") == "1" {
		fatalExit("help!")
		return
	}
	cmd := exec.Command(os.Args[0], "-test.run=TestFatalExit")
	cmd.Env = append(os.Environ(), "BE_CRASHER=1")
	err := cmd.Run()
	if e, ok := err.(*exec.ExitError); ok && !e.Success() {
		return
	}
	t.Fatalf("process ran with err %v, want exit status 1", err)
}

// Reference: https://npf.io/2015/06/testing-exec-command/

func TestTokenDirAndFile(t *testing.T) {

	before, _ := clientSecretFile()

	TokenDirAndFile = "/note/spud/.token.json"
	result, _ := clientSecretFile()
	assert.EqualValues(t, result, TokenDirAndFile)

	TokenDirAndFile = ""
	result, _ = clientSecretFile()
	assert.EqualValues(t, result, before)

}

func TestClientSecretFileNoToken(t *testing.T) {
	file, _ := clientSecretFile()
	assert.FileExistsf(t, file, "Token does not exit.")

}

func TestAddDeleteStation(t *testing.T) {
	station := "Elkins Park"
	//DeleteStation(station)
	AddStation(station)
}

// TestGetStationRecords -- live test and requires token.json
func TestGetStationRecords(t *testing.T) {
	station := "Elkins Park"
	expected := 1
	found := -1

	m := GetStationRecords(station, 1)
	tmp := strings.Split(m[0]["station"], ":")
	departureTime := tmp[1]

	fmt.Println("TMP", departureTime)
	fmt.Println("Live train_id:", m[0]["train_id"])

	i, err := strconv.Atoi(m[0]["train_id"])
	if err != nil {

	} else {
		if i > 100 {
			expected = i
			found = i
		}

	}
	assert.EqualValues(t, expected, found)
}

func TestGetStationRecordsWrapper(t *testing.T) {

	// TODO: Fix - Stations not all correct
	for i, station := range septa.ListStations() {
		//station := "Elkins Park"
		fmt.Printf("Station: %s\n", station)
		number := 3
		tmp := GetStationRecordsWrapper(station, number)
		assert.Contains(t, tmp[0]["station_rec_type"], station)
		assert.EqualValues(t, number*2, len(tmp), "Short on records")
		fmt.Printf(":%s", tmp[0]["time"])
		for _, rec := range tmp {
			fmt.Printf("\n:rec %v", rec)
		}

		if i >= 0 {
			break
		}

	}
}

func TestGetSingleDocument(t *testing.T) {
	SingleDocument()
}

func TestAllDocumentsSingleDocument(t *testing.T) {
	AllDocuments("trains", "train_id")
	AllDocuments("trainView", "TrainNo")
	//fmt.Println(m[0]["train_id"])
}

func TestDeleteDocument(t *testing.T) {
	DeleteDocument("trainView", "1507")
}

func TestRefreshLiveView(t *testing.T) {
	RefreshLiveView()
}

func TestGetAllStationsRecordsWrapper(t *testing.T) {

	success := false
	m := GetAllStationsRecordsWrapper(3)
	if len(m[0]["station_rec_type"]) > 3 {
		fmt.Printf("%s\n", m[0]["station_rec_type"])
		success = true
	}
	assert.Equal(t, true, success, "Network?")

}

func TestInsertUpdateDelete(t *testing.T) {
	insertUpdateDelete()
}

func TestNode(t *testing.T) {

}

func TestDateTimeParse(t *testing.T) {
	s := " April 2, 2018, 6:45 pm"
	tt, err := DateTimeParse(s).getTimeLoc()
	assert.Nil(t, err)

	t3, _ := DateTimeParse(s).getTimeLocSquish()
	fmt.Printf("Time: %v \n",
		t3)

	s = " Apr 2, 2018, 6:45 am"
	tt, err = DateTimeParse(s).getTime()
	assert.Nil(t, err)

	s = " Apr 2, 18, 6:45 am"
	tt, err = DateTimeParse(s).getTime()
	assert.Nil(t, err)

	fmt.Println(tt.Unix())
}

func TestQueryRRSchedules(t *testing.T) {

	QueryRRSchedules("452", "2018-04-03")
}

func TestDevExperimentingWithCollections(t *testing.T) {

	devExperimentingWithCollections()
}

func TestQuery(t *testing.T) {

	f, ctx, client := testQuery()
	defer client.Close()

	// Get the first 25 cities, ordered by population.
	//iter := f.OrderBy("population", firestore.Asc).Limit(25).Documents(ctx)

	iter := f.Documents(ctx)
	for {
		doc, err := iter.Next()

		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		fmt.Printf("HERE: created:%v, %v", doc.CreateTime, doc.Data())
	}
	//docs, err := firstPage.GetAll()

	fmt.Printf("\n ----------  Next -------------\n\n")
	iter = client.Collection("cities").Where("Capital", "==", true).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return
		}
		fmt.Println(doc.Data())
	}

	dsnap, err := client.Collection("trains").Doc("1054").Get(ctx)
	if err != nil {
		return
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m)

}

func TestAddRRSchedules(t *testing.T) {
	assert.Nil(t, AddRRSchedules(), "Running?")

}

func TestQueryRRSchedulesByDate(t *testing.T) {

	database := QueryRRSchedulesByDate("2018-04-10")
	for k, _ := range database {

		docDate := database[k][k].TrainRRSchedules.DocDate
		allStops := database[k][k].TrainRRSchedules.RRSchedules
		fmt.Printf("\n\n%v: %s \n", docDate, k)
		for _, v := range allStops {
			fmt.Printf("%s: %s,%s,%s\n",
				v.Station,
				v.SchedTM,
				v.EstTM,
				v.ActTM)
		}

	}

	// We know these values are correct
	testStation := database["330"]["330"].TrainRRSchedules.RRSchedules[0].Station
	assert.Equal(t, "Elwyn Station", testStation)
}

func TestAddStations(t *testing.T) {

	records, err := septa.GetLiveViewRecords()
	if err != nil {
		t.Fatal("Should not return nil")
	}
	if len(records) < 3 {
		t.Fatal("Not enough records to run test")
	}

	err = AddStations(records[0].TrainNo)
	if err != nil {
		t.Fatal("AddStations returned an error")
	}

	err = AddStations(records[1].TrainNo)
	if err != nil {
		t.Fatal("AddStations returned an error")
	}

}

func TestAddStationsByTime(t *testing.T) {

	//AddStationsByTime("412")
	//AddStationsByTime("426")
}

// This is too long
func TestAllStationsByTime(t *testing.T) {

	AllStationsByTime()

}

func TestGetTimeLocHRminS(t *testing.T) {

	r, _ := DateTimeParse("2018-04-10 3:01 pm").getTimeLocHRminS()

	assert.Equal(t, "15:01", r)

}
