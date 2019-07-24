package mongoConnetion

import (
	"context"
	"errors"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
	"time"
)

const (
	ShortUrlAlreadyExistsErrorMessage = "Shorten url version already exists. Please choose another one"
)

type Record struct {
	Url      string
	ShortUrl string
}

// return new MongoDb client for connection
func NewClient(user string, password string, url string, port string, database string, collection string) (*mongo.Collection, error) {
	connectionUrl := "mongodb://" + user + ":" + password + "@" + url + ":" + port
	clientOpts := options.Client().ApplyURI(connectionUrl)
	client, err := mongo.Connect(context.TODO(), clientOpts)
	if err != nil {
		return nil, err
	}
	c := client.Database(database).Collection(collection)
	return c, nil
}

// return pointer to list of all records
func GetAllRecords(client *mongo.Collection) *[]Record {
	cursor, err := client.Find(context.TODO(), bson.D{{}})
	if err != nil {
		log.Fatalln("Unable to find any records: %s", err)
	}

	var result []Record
	for cursor.Next(context.TODO()) {
		var elem Record
		err := cursor.Decode(&elem)
		if err != nil {
			log.Fatalln(err)
		}
		result = append(result, elem)
	}
	cursor.Close(context.TODO())
	return &result
}

func GetRecord(client *mongo.Collection, short string) string {
	records := GetAllRecords(client)
	for _, r := range *records {
		if r.ShortUrl == short {
			return r.Url
		}
	}
	return "/"
}

func CheckIfRecordExists(client *mongo.Collection, record Record) bool {
	records := GetAllRecords(client)
	for _, r := range *records {
		if record.ShortUrl == r.ShortUrl {
			return true
		}
	}
	return false
}

// insert one record, fail if the short name already exists
func InsertRecord(client *mongo.Collection, record Record) (interface{}, error) {
	exists := CheckIfRecordExists(client, record)
	if exists {
		return nil, errors.New(ShortUrlAlreadyExistsErrorMessage)
	}
	insertResult, err := client.InsertOne(context.TODO(), record)
	log.Printf("[%s] Inserting record for url: %s to with short url: %s", time.Now().Format(time.ANSIC), record.Url, record.ShortUrl)
	return insertResult, err

}
