package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gorilla/mux"
)

var path = "lurls.txt"

func TagRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	match := false
	tag := params["tag"]
	log.Print(fmt.Sprintf("Request for %s using %s from %s", tag, r.Proto, r.RemoteAddr))
	var file, err = os.OpenFile(path, os.O_RDWR, 0644)
	if err != nil {
		// URL file not found or can't be opened
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("URL file not found or can't be opened"))
		log.Print("Received request but URL file was missing or could not be opened")
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := strings.Fields(scanner.Text())
		if v[0] == tag {
			log.Print(fmt.Sprintf("Match found for tag %s. Redirecting to %s", tag, v[1]))
			match = true
			http.Redirect(w, r, v[1], http.StatusFound)
			break
		}
	}
	if match == false {
		log.Print(fmt.Sprintf("No match found for tag %s", tag))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Tag not found"))
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{tag}", TagRequest).Methods("GET")

	fmt.Println("Server at 8080")
	log.Fatal(http.ListenAndServe(":8080", router))

}
