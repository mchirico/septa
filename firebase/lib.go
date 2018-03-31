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
)

// Version -- version of this program
var Version = "0.0.1"

// TokenDirAndFile -- use this for system router
var TokenDirAndFile = ""

// QueryTime -- seconds, used to update data
var QueryTime = 20

func help() {
	os.Exit(1)
}

// Flags -- used in routerfirebase
func Flags() {

	helpVar := false
	flag.BoolVar(&helpVar, "help", false, "help listing")
	flag.IntVar(&QueryTime, "time", 20, "time in seconds")
	flag.StringVar(&TokenDirAndFile, "token",
		"",
		"directory and file of token.json /stuff/token.json")
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
	defer client.Close()

	_, err = client.Collection(collection).Doc(document).Delete(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// AllDocuments -- returns all documents in the collection
func AllDocuments(collection string, document string) []map[string]interface{} {

	var database []map[string]interface{}

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
	defer client.Close()

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
	defer client.Close()

	number := 3
	records := GetStationRecordsWrapper(station, number)

	for _, rec := range records {
		_, err := client.Collection("trains").Doc(rec["train_id"]).Set(ctx, rec)
		fmt.Printf("Train updated: %s\n", rec["train_id"])
		if err != nil {
			fmt.Printf("error on insert collection")
		}

	}

}

// RefreshLiveView -- need to clean this up
func RefreshLiveView() {
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
	defer client.Close()

	records := septa.GetLiveViewRecords()

	oldRecords := AllDocuments("trainView", "TrainNo")
	trainMap := map[string]int{}

	for _, rec := range records {
		_, err := client.Collection("trainView").Doc(rec.TrainNo).Set(ctx, rec)
		trainMap[rec.TrainNo] = 1
		if err != nil {
			fmt.Printf("error on insert collection")
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
