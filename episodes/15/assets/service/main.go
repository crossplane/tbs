package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
)

var db *gorm.DB
var err error

const titlepage = `
<html>
<h1>{{ range $i := .}}{{$i.Name}}, {{$i.Location}}, <a href="https://twitter.com/{{$i.Twitter}}">@{{$i.Twitter}}</a><br/>
{{end}}</h1>
</html>`

type TBSStream struct {
	Name     string `json:"name"`
	Location string `json:"location"`
	Twitter  string `json:"twitter"`
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
	b, err := ioutil.ReadFile("/var/bindings/sql/connectionstr")
	if err != nil || len(b) == 0 {
		user := os.Getenv("DB_USER")
		password := os.Getenv("DB_PASSWORD")
		host := os.Getenv("DB_HOST")
		db, err = gorm.Open("postgres", fmt.Sprintf("postgres://%s:%s@%s:5432/postgres", user, password, host))
	} else {
		db, err = gorm.Open("postgres", string(b))
	}

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
	fmt.Fprint(w, "Welcome to the 16th episode of The Binding Status!")
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
	t := template.Must(template.New("List").Parse(titlepage))
	t.Execute(w, streams)
}
