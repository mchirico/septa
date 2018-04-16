package firebase

import (
	"firebase.google.com/go"
	septa "github.com/mchirico/septa/utils"
	"golang.org/x/net/context"
	"google.golang.org/api/option"
	"log"
	"net/url"
	"os"
	"os/user"
	"path/filepath"
	"strings"

	"cloud.google.com/go/firestore"
	"flag"
	"fmt"
	"google.golang.org/api/iterator"
	"time"
)

// Version -- version of this program
var Version = "0.0.1"

// TokenDirAndFile -- use this for system router
var TokenDirAndFile = ""

// QueryTime -- seconds, used to update data
var QueryTime = 20

// QuietMode -- router no output if true
var QuietMode = false

// Doc -- struct for entering RRSchedules
type Doc struct {
	DateTrainID      string
	TrainRRSchedules septa.TrainRRSchedules
}

// Record --
type Record map[string]Doc

// Database --
type Database map[string]Record

// StationStop --
type StationStop struct {
	Station string
	ActTM   string
	EstTM   string
	SchedTM string
}

func fatalExit(msg string) {

	// Add future update to firebase
	log.Fatalf("Fatal Call")

}

// Flags -- used in routerfirebase
func Flags() {

	flag.BoolVar(&QuietMode, "quiet", false, "set to true"+
		" for no quiet output:  ./routefirebase -quiet=true")

	flag.IntVar(&QueryTime, "time", 20, "time in seconds")
	flag.StringVar(&TokenDirAndFile, "token",
		"",
		"directory and file of token.json\n"+
			"   ./routefirebase -token='/stuff/token.json'")
	flag.Parse()

}

func clientSecretFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}

	tokenCacheDir := filepath.Join(usr.HomeDir, ".google_firebase")

	if TokenDirAndFile != "" {
		return TokenDirAndFile, nil
	}

	file := filepath.Join(tokenCacheDir,
		url.QueryEscape("token.json"))

	// Can be used for travis-cli
	if _, err := os.Stat(file); os.IsNotExist(err) {
		tokenCacheDir = "../"
	}
	return filepath.Join(tokenCacheDir,
		url.QueryEscape("token.json")), err
}

//GetStationRecords Builds database
func GetStationRecords(station string, number int) []map[string]string {

	var database []map[string]string
	return septa.GetStationRecords(station, number, database)
}

//GetStationRecordsWrapper Adds timestamp to records
func GetStationRecordsWrapper(
	station string, number int) []map[string]string {

	records := GetStationRecords(station, number)
	tmp := strings.Split(records[0]["station"], ":")

	stationRecType := tmp[0]

	for k := range records {
		records[k]["timestamp"] = fmt.Sprintf("%s:%s", tmp[1],
			tmp[2])
		records[k]["station_rec_type"] = stationRecType
	}
	return records

}

// GetAllStationsRecordsWrapper -- look closely at this
func GetAllStationsRecordsWrapper(number int) []map[string]string {

	groupOfRecords := septa.GetAllStationRecords(number)
	for _, records := range groupOfRecords {

		tmp := strings.Split(records["station"], ":")

		stationRecType := tmp[0]

		records["timestamp"] = fmt.Sprintf("%s:%s", tmp[1],
			tmp[2])
		records["station_rec_type"] = stationRecType

	}
	return groupOfRecords
}

// DeleteDocument -- simple document delete
func DeleteDocument(collection string, document string) {

	ctx, client := OpenCtxClient()
	defer client.Close()

	_, err := client.Collection(collection).Doc(document).Delete(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// AllDocuments -- returns all documents in the collection
func AllDocuments(collection string, document string) []map[string]interface{} {

	ctx, client := OpenCtxClient()
	defer client.Close()

	var database []map[string]interface{}

	iter := client.Collection(collection).Documents(ctx)
	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			return database
		}
		database = append(database, doc.Data())
	}
	return database
}

// SingleDocument -- grap a single document in the collection
func SingleDocument() {

	ctx, client := OpenCtxClient()
	defer client.Close()

	dsnap, err := client.Collection("trains").Doc("3420").Get(ctx)
	if err != nil {
		return
	}
	m := dsnap.Data()
	fmt.Printf("Document data: %#v\n", m["station"])

}

// AddStation -- add station to Firestore - 3 records
func AddStation(station string) {

	ctx, client := OpenCtxClient()
	defer client.Close()

	number := 3
	records := GetStationRecordsWrapper(station, number)

	for _, rec := range records {
		_, err := client.Collection("trains").Doc(rec["train_id"]).Set(ctx, rec)
		if QuietMode == false {
			fmt.Printf("Train updated: %s\n", rec["train_id"])
		}

		if err != nil {
			fmt.Printf("error on insert collection")
		}

	}

}

// AllStationsByTime --
func AllStationsByTime() {

	records, err := septa.GetLiveViewRecords()
	if err != nil {
		return
	}

	for _, train := range records {
		go AddStationsByTime(train.TrainNo)
	}
}

// AddStationsByTime -- add station data
func AddStationsByTime(trainNo string) error {
	ctx, client := OpenCtxClient()
	defer client.Close()
	//trainNo := "412"
	rrSchedules := septa.GetRRSchedules(trainNo)

	if QuietMode == false {
		fmt.Printf("Station updated for trainNo: %s\n", trainNo)
	}

	for i := range rrSchedules.RRSchedules {

		m := rrSchedules.RRSchedules[i]
		timeSquish, _ := DateTimeParse(
			fmt.Sprintf("%s %s",
				rrSchedules.DocDate,
				m.SchedTM)).getTimeLocSquish()

		r := map[string]string{"Station": m.Station,
			"ActTM": m.ActTM, "EstTM": m.EstTM,
			"SchedTM": m.SchedTM,
			"TrainID": rrSchedules.TrainID}

		SchedTM, _ := DateTimeParse(
			fmt.Sprintf("%s %s",
				rrSchedules.DocDate,
				m.SchedTM)).getTimeLocHRminS()

		_, err := client.Collection("StationsByTime").
			Doc(rrSchedules.DocDate).Collection(rrSchedules.RRSchedules[i].Station).
			Doc(timeSquish).Collection(rrSchedules.TrainID).Doc(SchedTM).Set(ctx, r, firestore.MergeAll)

		if err != nil {
			fmt.Printf("error on insert")
		}
	}
	return nil
}

// AddStations -- add station data
func AddStations(trainNo string) error {
	ctx, client := OpenCtxClient()
	defer client.Close()
	//trainNo := "412"
	rrSchedules := septa.GetRRSchedules(trainNo)

	fmt.Printf("%v\n", rrSchedules.DocDate)
	fmt.Printf("%v\n", rrSchedules.RRSchedules[0].Station)

	fmt.Printf("%v\n", rrSchedules.RRSchedules[0].ActTM)
	fmt.Printf("%v\n", rrSchedules.RRSchedules[0].ActTM)

	for i := range rrSchedules.RRSchedules {

		m := rrSchedules.RRSchedules[i]

		r := map[string]string{"Station": m.Station,
			"ActTM": m.ActTM, "EstTM": m.EstTM,
			"SchedTM": m.SchedTM,
			"trainNo": trainNo}

		_, err := client.Collection("Stations").
			Doc(rrSchedules.DocDate).Collection(rrSchedules.RRSchedules[i].Station).
			Doc(trainNo).Set(ctx, r, firestore.MergeAll)

		if err != nil {
			fmt.Printf("error on insert")
			return err
		}
	}
	return nil
}

// AddRRSchedules -- mergeAll
func AddRRSchedules() error {

	ctx, client := OpenCtxClient()
	defer client.Close()

	records, err := septa.GetLiveViewRecords()
	if err != nil {
		return err
	}

	for _, train := range records {

		rrSchedules := septa.GetRRSchedules(train.TrainNo)

		doc := Doc{}
		doc.DateTrainID = fmt.Sprintf("%s:%s", train.TrainNo, rrSchedules.DocDate)
		doc.TrainRRSchedules = rrSchedules

		m := map[string]Doc{}
		m[train.TrainNo] = doc
		_, err := client.Collection("rrSchedules").
			Doc(rrSchedules.DocDate).Collection(train.TrainNo).Doc(rrSchedules.DocDate).Set(ctx, m, firestore.MergeAll)

		if err != nil {
			fmt.Printf("error on insert collection")
			return err
		}

	}

	return nil
}

// RefreshLiveView -- need to clean this up
func RefreshLiveView() {

	ctx, client := OpenCtxClient()
	defer client.Close()

	records, err := septa.GetLiveViewRecords()
	if err != nil {
		return
	}

	oldRecords := AllDocuments("trainView", "TrainNo")
	trainMap := map[string]int{}

	for _, rec := range records {
		_, err := client.Collection("trainView").Doc(rec.TrainNo).Set(ctx, rec)
		if err != nil {
			fmt.Printf("error on insert collection")
		} else {
			trainMap[rec.TrainNo] = 1
		}

	}

	for _, v := range oldRecords {
		trainNo := v["TrainNo"].(string)
		if _, ok := trainMap[trainNo]; ok {

		} else {
			DeleteDocument("trainView", trainNo)
		}
	}

}

// OpenCtxClient -- open a connection
func OpenCtxClient() (context.Context, *firestore.Client) {
	ctx := context.Background()
	file, _ := clientSecretFile()
	sa := option.WithCredentialsFile(file)
	app, err := firebase.NewApp(ctx, nil, sa)
	if err != nil {
		log.Fatalln(err)
	}
	client, err := app.Firestore(ctx)
	if err != nil {
		log.Fatalln(err)
	}
	// defer client.Close()  //You must do this
	return ctx, client
}

// insertUpdateDelete -- testing only
//  Used for testing only
func insertUpdateDelete() {

	ctx, client := OpenCtxClient()
	defer client.Close()

	type H struct {
		A string
	}

	type City struct {
		Name       string
		State      string
		Country    string
		Capital    bool
		Population int64
		Timestamp  time.Time
		H          H
	}

	cities := []struct {
		id string
		c  City
	}{
		{id: "SF", c: City{Name: "San Francisco", State: "CA", Country: "USA", Capital: false, Population: 860000}},
		{id: "LA", c: City{Name: "Los Angeles", State: "CA", Country: "USA", Capital: false, Population: 3900000}},
		{id: "DC", c: City{Name: "Washington D.C.", Country: "USA", Capital: false, Population: 680000}},
		{id: "TOK", c: City{Name: "Tokyo", Country: "Japan", Capital: true, Population: 9000000}},
		{id: "BJ", c: City{Name: "Beijing", Country: "China", Capital: false, Population: 21500000}},
	}
	for _, c := range cities {
		_, err := client.Collection("cities").Doc(c.id).Set(ctx, c.c)
		if err != nil {
			return
		}
	}
}

// devExperimentingWithCollections -- used to do tests on collections and documents
func devExperimentingWithCollections() (firestore.Query, context.Context, *firestore.Client) {
	ctx, client := OpenCtxClient()

	type H struct {
		Z map[string]string
	}

	type Hb struct {
		A H
		M H
	}

	z := H{}
	hb := Hb{}
	m := map[string]string{}
	m["test"] = "s3"
	m["bob"] = "34"

	z.Z = m
	hb.A = z

	fmt.Printf("h.b: %v", hb)

	client.Collection("testTrain").Doc("2018-04-03").
		Collection("345").Doc("A").Collection("1").Doc("id").Set(ctx, hb)

	m["Bob"] = "Bounce"
	client.Collection("testTrain").Doc("2018-04-03").
		Collection("345").Doc("B").Collection("1").Doc("id").Set(ctx, hb)

	iter := client.Collection("testTrain").Doc("2018-04-03").
		Collection("345").Doc("A").Collection("1").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			print("iter done\n\n")
			break
		}
		if err != nil {
			fmt.Println("error")
		}
		fmt.Println("doc.Data():", doc.Data())
	}

	query := client.Collection("testTrain").Where("A", "==", map[string]string{"test": "ss3"})
	return query, ctx, client

}

// testQuery -- for experiments
func testQuery() (firestore.Query, context.Context, *firestore.Client) {
	ctx, client := OpenCtxClient()

	query := client.Collection("cities").Where("State", "==", "CA")

	iter := client.Collection("cities").Where("State", "==", "CA").Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			print("iter done\n\n")
			break
		}
		if err != nil {
			fmt.Println("error")
		}
		fmt.Println("doc.Data():", doc.Data())
	}

	return query, ctx, client
}

// QueryRRSchedulesByDate --
func QueryRRSchedulesByDate(docDate string) Database {

	ctx, client := OpenCtxClient()
	defer client.Close()

	database := Database{}

	iter := client.Collection("rrSchedules").
		Doc(docDate).Collections(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			break
		}
		if err != nil {
			fmt.Println("error")
		}

		iter2 := doc.Documents(ctx)
		for {
			doc, err := iter2.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				fmt.Println("error")
			}

			record := Record{}
			doc.DataTo(&record)

			for _, r := range record {
				//fmt.Printf("value = %v\n",v)
				database[r.TrainRRSchedules.TrainID] = record
				//fmt.Printf("TrainID:%v\n", r.DateTrainID)
				//fmt.Printf("v.TrainRRSchedules.RRSchedules: %v\n", r.TrainRRSchedules.RRSchedules)
				//schedules := r.TrainRRSchedules.RRSchedules[0]
				//
				//fmt.Printf("Station: %v, %v\n", schedules.Station, schedules.ActTM)
			}

		}

	}

	return database
}

// QueryRRSchedules  -- will be used for querying.
func QueryRRSchedules(trainNo string, docDate string) (firestore.Query, context.Context,
	*firestore.Client) {

	dateTrainID := fmt.Sprintf("%s:%s", trainNo, docDate)

	ctx, client := OpenCtxClient()
	defer client.Close()

	iter := client.Collection("rrSchedules").Doc(docDate).
		Collection(trainNo).Documents(ctx)

	for {
		doc, err := iter.Next()
		if err == iterator.Done {
			print("\n\niter done\n\n")
			break
		}
		if err != nil {
			fmt.Println("error")
		}
		fmt.Println("doc.Data()[2357]:", doc.Data()[trainNo])

		type Doc struct {
			DateTrainID      string
			TrainRRSchedules septa.TrainRRSchedules
		}

		k2 := map[string]Doc{}
		doc.DataTo(&k2)

		fmt.Printf("k2[trainNo].DateTrainID:  %v\n",
			k2[trainNo].DateTrainID)
		fmt.Printf("k2[trainNo].TrainRRSchedules:  %v\n",
			k2[trainNo].TrainRRSchedules)
		fmt.Printf("k2[trainNo].TrainRRSchedules.Timestamp: %v\n",
			k2[trainNo].TrainRRSchedules.Timestamp)
		fmt.Printf("k2[trainNo].TrainRRSchedules.DocDate: %v\n",
			k2[trainNo].TrainRRSchedules.DocDate)
		fmt.Printf("k2[trainNo].TrainRRSchedules.RRSchedules[0].Station: %v\n",
			k2[trainNo].TrainRRSchedules.RRSchedules[0].Station)

		fmt.Printf("\n ---------------------------:\n")
		for _, v := range k2[trainNo].TrainRRSchedules.RRSchedules {

			fmt.Printf("%s: %s %s %s\n", v.Station, v.ActTM, v.EstTM, v.SchedTM)

		}

	}

	query := client.Collection("rrSchedules").
		Where("DateTrainID", "==", dateTrainID)

	return query, ctx, client
}

// DateTimeParse -- takes are variety of expected dates
type DateTimeParse string

// getTime --
func (s DateTimeParse) getTime() (time.Time, error) {
	layout := []string{
		"January 2, 2006, 3:04 pm",
		"January 2, 2006, 3:04pm",
		"January 2, 2006, 03:04 pm",
		"January 2 2006, 03:04 pm",
		"January 2 2006 03:04 pm",

		"January 2, 2006, 3:04 pm",
		"January 2 2006, 3:04 pm",
		"January 2 2006 3:04 pm",
		"Jan 2, 2006, 03:04 pm",
		"Jan 2 2006, 03:04 pm",
		"Jan 2, 2006, 3:04 pm",
		"Jan 2, 06, 3:04 pm",
		"2006-01-02 3:04 pm",
		"2006-01-02 3:04pm",
		"2006-01-02 3:04 PM",
		"2006-01-02 3:04PM",
		"2006-01-02 15:04",
	}

	st := strings.Join(strings.Fields(string(s)), " ")
	//fmt.Printf("-->%s\n", st)

	for _, l := range layout[:len(layout)-1] {
		t, err := time.Parse(l, st)
		if err == nil {
			return t, err
		}

	}

	return time.Parse(layout[len(layout)-1], st)

}

// getTimeLoc --
func (s DateTimeParse) getTimeLoc() (time.Time, error) {

	tt, err := DateTimeParse(s).getTime()
	if err != nil {
		return tt, err
	}
	loc, err := time.LoadLocation("America/New_York")

	return tt.In(loc), err

}

func (s DateTimeParse) getTimeLocSquish() (string, error) {

	tt, err := DateTimeParse(s).getTime()
	if err != nil {
		return "", err
	}
	squishMin := int(tt.Minute()/10) * 10
	ret := fmt.Sprintf("%02d:%02d", tt.Hour(), squishMin)
	return ret, err

}

func (s DateTimeParse) getTimeLocHRminS() (string, error) {

	tt, err := DateTimeParse(s).getTime()
	if err != nil {
		return "", err
	}
	ret := fmt.Sprintf("%02d:%02d", tt.Hour(), tt.Minute())
	return ret, err

}
