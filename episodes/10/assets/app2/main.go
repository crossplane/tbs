package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
)

var db *gorm.DB
var err error

type TBSStream struct {
	Title   string `json:"title"`
	Host    string `json:"host"`
	Viewers int    `json:"viewers"`
}

func handleRequests() {
	log.Println("Starting TBSStreams server...")
	myRouter := mux.NewRouter().StrictSlash(true)
	myRouter.HandleFunc("/", home)
	myRouter.HandleFunc("/create", create).Methods("POST")
	myRouter.HandleFunc("/list", list)
	log.Fatal(http.ListenAndServe(":8080", myRouter))
}

func main() {
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	host := os.Getenv("DB_HOST")
	db, err = gorm.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:3306)/okteto?charset=utf8&parseTime=True", user, password, host))
	defer db.Close()

	if err != nil {
		log.Println("Connection Failed to Open")
	} else {
		log.Println("Connection Established")
	}

	db.AutoMigrate(&TBSStream{})
	handleRequests()
}

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Fprint(w, "Welcome to the 11th episode of The Binding Status!")
	fmt.Println("Request Received: /")
}

func create(w http.ResponseWriter, r *http.Request) {
	reqBody, _ := ioutil.ReadAll(r.Body)
	var stream TBSStream
	json.Unmarshal(reqBody, &stream)
	db.Create(&stream)
	fmt.Println("Request Received: /create")
	json.NewEncoder(w).Encode(stream)
}

func list(w http.ResponseWriter, r *http.Request) {
	streams := []TBSStream{}
	db.Find(&streams)
	fmt.Println("Request Received: /list")
	json.NewEncoder(w).Encode(streams)
}
