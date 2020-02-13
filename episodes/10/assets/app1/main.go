package main

import (
	"fmt"
	"net/http"
	"os"
)

func main() {
	fmt.Println("Starting hello-world server...")
	http.HandleFunc("/", helloServer)
	if err := http.ListenAndServe(":8080", nil); err != nil {
		panic(err)
	}
}

func helloServer(w http.ResponseWriter, r *http.Request) {
	message := os.Getenv("TBS_MESSAGE")
	fmt.Fprint(w, message)
}
