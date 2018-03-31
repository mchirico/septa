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
	"fmt"
	"google.golang.org/api/iterator"
)

func clientSecretFile() (string, error) {
	usr, err := user.Current()
	if err != nil {
		return "", err
	}
	tokenCacheDir := filepath.Join(usr.HomeDir, ".google_firebase")
	os.MkdirAll(tokenCacheDir, 0700)
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

	station_rec_type := tmp[0]

	for k := range records {
		records[k]["timestamp"] = fmt.Sprintf("%s:%s", tmp[1],
			tmp[2])
		records[k]["station_rec_type"] = station_rec_type
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

// DeleteStation -- simple delete test
func DeleteStation(station string) {
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

	_, err = client.Collection(station).Doc("*").Delete(ctx)
	if err != nil {
		log.Fatalln(err)
	}
}

// AllDocuments -- returns all documents in the collection
func AllDocuments(collection string) []map[string]interface{} {

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

	fmt.Println("All Trains:")
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
		fmt.Println(doc.Data()["train_id"])
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
