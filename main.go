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

func TagRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	match := false
	tag := params["tag"]

	var lurl_path string
	var err error
	var path_override, path_set = os.LookupEnv("LURLS")
	if !path_set {
		lurl_path, err = os.Getwd()
		if err != nil {
			w.WriteHeader(http.StatusInternalServerError)
			w.Write([]byte("Could not determine working directory"))
			log.Print("Could not determine working directory")
			return
		}
	} else {
		lurl_path = path_override
	}

	var file *os.File
	file, err = os.OpenFile(lurl_path+"/lurls.txt", os.O_RDONLY, 0644)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("URL file not found or can't be opened"))
		log.Print("Received request but URL file was missing or could not be opened")
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := strings.Fields(scanner.Text())
		if v[0] == tag {
			log.Print(fmt.Sprintf("Request for %s (%s) using %s from %s", tag, v[1], r.Proto, r.RemoteAddr))
			match = true
			http.Redirect(w, r, v[1], http.StatusFound)
			break
		}
	}
	if match == false {
		log.Print(fmt.Sprintf("Request for unknown tag %s using %s from %s", tag, r.Proto, r.RemoteAddr))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Tag not found"))
	}
}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/{tag}", TagRequest).Methods("GET")

	var port, port_set = os.LookupEnv("PORT")
	if !port_set {
		port = "8080"
		log.Print(fmt.Sprintf("Using default port: %s", port))
	} else {
		log.Print(fmt.Sprintf("Port override found: %s", port))
	}

	log.Print(fmt.Sprintf("Server listening at %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
