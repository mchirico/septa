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

func help() {
	os.Exit(1)
}

// Flags -- used in routerfirebase
func Flags() {

	helpVar := false
	flag.BoolVar(&helpVar, "help", false, "help listing")

	flag.BoolVar(&QuietMode, "quiet", false, "set to true"+
		" for no quiet output:  ./routefirebase -quiet=true")

	flag.IntVar(&QueryTime, "time", 20, "time in seconds")
	flag.StringVar(&TokenDirAndFile, "token",
		"",
		"directory and file of token.json\n"+
			"   ./routefirebase -token='/stuff/token.json'")
	flag.Parse()

	if helpVar {
		help()
	}

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

func deleteCollection(ctx context.Context, client *firestore.Client,
	ref *firestore.CollectionRef, batchSize int) error {

	for {
		// Get a batch of documents
		iter := ref.Limit(batchSize).Documents(ctx)
		numDeleted := 0

		// Iterate through the documents, adding
		// a delete operation for each one to a
		// WriteBatch.
		batch := client.Batch()
		for {
			doc, err := iter.Next()
			if err == iterator.Done {
				break
			}
			if err != nil {
				return err
			}

			batch.Delete(doc.Ref)
			numDeleted++
		}

		// If there are no documents to delete,
		// the process is over.
		if numDeleted == 0 {
			return nil
		}

		_, err := batch.Commit(ctx)
		if err != nil {
			return err
		}
	}
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

// AddRRSchedules -- mergeAll
func AddRRSchedules() error {

	ctx, client := OpenCtxClient()
	defer client.Close()

	records := septa.GetLiveViewRecords()

	type Doc struct {
		DateTrainID      string
		TrainRRSchedules septa.TrainRRSchedules
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

// AddAllStations -- add all stations to Firestore
func AddAllStations(number int) {

	ctx, client := OpenCtxClient()
	defer client.Close()

	records := GetAllStationsRecordsWrapper(number)

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

// RefreshLiveView -- need to clean this up
func RefreshLiveView() {

	ctx, client := OpenCtxClient()
	defer client.Close()

	records := septa.GetLiveViewRecords()

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

func QueryRRSchedules(trainNo string, docDate string) (firestore.Query, context.Context,
	*firestore.Client) {

	dateTrainID := fmt.Sprintf("%s:%s", trainNo, docDate)

	ctx, client := OpenCtxClient()
	defer client.Close()

	//iter := client.Collection("rrSchedules").
	//	Where("3236.DateTrainID", "==", "3236:2018-04-03").Documents(ctx)

	iter := client.Collection("rrSchedules").
		Where("434.DateTrainID", "==", "432:2018-04-03").Documents(ctx)

	//client.Collection("rrSchedules").Doc("2018-04-03").Delete(ctx)
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

	query := client.Collection("rrSchedules").
		Where("DateTrainID", "==", dateTrainID)

	return query, ctx, client
}

// DateTimeParse -- takes are variety of expected dates
func DateTimeParse(s string) (time.Time, error) {
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

	s = strings.Join(strings.Fields(s), " ")
	fmt.Printf("-->%s\n", s)

	for _, l := range layout[:len(layout)-1] {
		t, err := time.Parse(l, s)
		if err == nil {
			return t, err
		}

	}

	return time.Parse(layout[len(layout)-1], s)

}
