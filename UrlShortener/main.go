package main

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/casek14/GoPhercises/UrlShortener/mongoConnetion"
	"go.mongodb.org/mongo-driver/bson"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

const (
	jsonFilePath         = "./short.json"
	mongoConnetionString = "localhost"
	databaseName         = "test"
	collectionName       = "trainers"
	mongoUser            = "root"
	mongoPassword        = "r00tme"
	mongoPort            = "27017"
)

type UrlRecord struct {
	Url       string `json:"url"`
	ShortName string `json:"short"`
}

type UrlsCatalog struct {
	Catalog []UrlRecord `json:"urls"`
}

type MongoConnectionConfig struct {
	DbUrl string
	DbPort string
	DbName string
	DbCollectionName string
	DbUser string
	DbPassword string
}

// Load json file and parse long and short urls
// Returns a catalog
func LoadCatalog(f string) map[string]string {

	jsonFile, err := os.Open(f)
	if err != nil {
		log.Fatalf("Unable to open %s file. %s", jsonFile, err)
	}

	var catalog UrlsCatalog
	bytes, _ := ioutil.ReadAll(jsonFile)
	err = json.Unmarshal(bytes, &catalog)
	if err != nil {
		log.Fatalf("Unable to load json. %s", err)
	}
	log.Println("JSON FILE LOADED !!")
	mapa := make(map[string]string)
	for _, u := range catalog.Catalog {
		mapa[u.ShortName] = u.Url
	}
	return mapa
}

func MapHandler(mapa map[string]string, fallback http.Handler) http.HandlerFunc {
	return func(writer http.ResponseWriter, request *http.Request) {
		path := request.URL.Path
		log.Printf("Handling redirect for %s url!", path)
		if destination, err := mapa[path]; err {
			log.Printf("------------ HANDLING PATH: %s\n",path	)
			http.Redirect(writer, request, destination, http.StatusFound)
			return
		}
		fallback.ServeHTTP(writer, request)
	}
}

func defaultMux() *http.ServeMux {
	mux := http.NewServeMux()
	mux.HandleFunc("/", hello)
	mux.HandleFunc("/register", register)
	mux.HandleFunc("/list", listRecords)
	mux.HandleFunc("/short/", redirect)
	return mux
}

func hello(w http.ResponseWriter, r *http.Request) {
	log.Println("Handling default Hello page !")
}

func listRecords(w http.ResponseWriter, r *http.Request) {
	config := initDbConnectionConfig()
	client, err := mongoConnetion.NewClient(config.DbUser, config.DbPassword, config.DbUrl, config.DbPort, config.DbName, config.DbCollectionName)

	if err != nil {
		log.Fatalln("Failed to get mongo client")
	}
	records := mongoConnetion.GetAllRecords(client)
	for _, r := range *records {
		fmt.Fprintf(w, "Short url: %s\n", r.ShortUrl)
		fmt.Fprintf(w, "Url: %s\n", r.Url)
		fmt.Fprintln(w, "=========================================")
	}
}

func register(w http.ResponseWriter, r *http.Request) {
	if r.Method == http.MethodGet {
		t, _ := template.ParseFiles("register.gtpl")
		t.Execute(w, nil)
	} else {
		r.ParseForm()
		config := initDbConnectionConfig()
		client, err := mongoConnetion.NewClient(config.DbUser, config.DbPassword, config.DbUrl, config.DbPort, config.DbName, config.DbCollectionName)

		if err != nil {
			log.Fatalln("Failed to get mongo client")
		}
		if r.FormValue("url") != "" && r.FormValue("short") != "" {
			_, err = mongoConnetion.InsertRecord(client, mongoConnetion.Record{Url: r.FormValue("url"), ShortUrl: r.FormValue("short")})
		}
		http.Redirect(w, r, "/list", http.StatusFound)
	}
}

func redirect(w http.ResponseWriter, r *http.Request) {
	newUrl := r.URL.Path[7:]
	fmt.Println(newUrl)
	config := initDbConnectionConfig()
	client, err := mongoConnetion.NewClient(config.DbUser, config.DbPassword, config.DbUrl, config.DbPort, config.DbName, config.DbCollectionName)

	if err != nil {
		log.Fatalln("Failed to get mongo client")
	}
	redirectUrl := mongoConnetion.GetRecord(client, newUrl)
	http.Redirect(w, r, redirectUrl, http.StatusFound)
}
// Create Mongo connection config from env variables
// DbUrl - url to mongo db (DBURL)
// DbPort - mongo port (DBPORT)
// DbUser - user used for db connection (DBUSER)
// DbPassword - password used for db connection (DBPASSWORD)
// DbName - name of database to use (DBNAME)
// DbCollectionName - name of collection to use (DBCOLLECTIONNAME)
func initDbConnectionConfig() *MongoConnectionConfig{
	var config MongoConnectionConfig
	envVars := map[string]string{"DBURL":"","DBPORT":"","DBUSER":"","DBPASSWORD":"","DBNAME":"","DBCOLLECTIONNAME":""}
	for env, _ := range envVars{
		value, ok := os.LookupEnv(env)
		if ok {
			envVars[env] = value
		}
	}
	config.DbUrl = envVars["DBURL"]
	config.DbPort = envVars["DBPORT"]
	config.DbUser = envVars["DBUSER"]
	config.DbPassword = envVars["DBPASSWORD"]
	config.DbName = envVars["DBNAME"]
	config.DbCollectionName = envVars["DBCOLLECTIONNAME"]
	log.Printf("MONGO CONNECTION STRING= %+v\n",config)
	return &config
}


func main() {
	mux := defaultMux()
	catalog := LoadCatalog(jsonFilePath)
	config := initDbConnectionConfig()
	client, err := mongoConnetion.NewClient(config.DbUser, config.DbPassword, config.DbUrl, config.DbPort, config.DbName, config.DbCollectionName)
	if err != nil {
		log.Fatalln("Failed to get mongo client")
	}
	//Clean database
	_, _ = client.DeleteMany(context.TODO(), bson.D{{}})
	fmt.Println("Starting server on port :8080")
	http.ListenAndServe(":8080", MapHandler(catalog, mux))
}
