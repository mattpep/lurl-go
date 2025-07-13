package main

import (
	"bufio"
	"fmt"
	"log"
	"net/http"
	"os"
	"runtime/debug"
	"strings"

	"github.com/gorilla/mux"
	"golang.org/x/sys/unix"
)

func isatty() bool {
	_, err := unix.IoctlGetWinsize(int(os.Stdout.Fd()), unix.TIOCGWINSZ)
	if err != nil {
		return false
	} else {
		return true
	}
}

func addressFromRequest(r *http.Request) string {
	addr := r.Header.Get("X-Real-IP")
	if addr != "" {
		return addr
	}

	addr = r.Header.Get("X-Forwarded-For")
	if addr != "" {
		return addr
	}

	return r.RemoteAddr
}

func NoTagRequest(w http.ResponseWriter, r *http.Request) {
	log.Print(fmt.Sprintf("Request for index page using %s from %s", r.Proto, addressFromRequest(r)))
	w.WriteHeader(http.StatusNotFound)
	w.Write([]byte("No tag specified"))
}
func TagRequest(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	match := false
	tag := params["tag"]

	source := addressFromRequest(r)

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
		log.Print(fmt.Sprintf("Received request from %s but URL file was missing or could not be opened (attempted path was "+lurl_path+")", source))
		return
	}

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		v := strings.Fields(scanner.Text())
		if v[0] == tag {
			log.Print(fmt.Sprintf("Request for %s (%s) using %s from %s", tag, v[1], r.Proto, source))
			match = true
			http.Redirect(w, r, v[1], http.StatusFound)
			break
		}
	}
	if match == false {
		log.Print(fmt.Sprintf("Request for unknown tag %s using %s from %s", tag, r.Proto, source))
		w.WriteHeader(http.StatusNotFound)
		w.Write([]byte("Tag not found"))
	}
}

func main() {
	router := mux.NewRouter()
	var log_file *os.File
	var err error

	// in headless mode we write logs to a file, else to the screen
	if !isatty() {
		log_file, err = os.OpenFile("/var/log/lurl/lurl.log", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0640)
		if err != nil {
			log.Fatal(err)
		} else {
			log.SetOutput(log_file)
		}
	}

	router.HandleFunc("/{tag}", TagRequest).Methods("GET")
	router.HandleFunc("/", NoTagRequest).Methods("GET")

	var port, port_set = os.LookupEnv("PORT")
	if !port_set {
		port = "8080"
		log.Print(fmt.Sprintf("Using default port: %s", port))
	} else {
		log.Print(fmt.Sprintf("Port override found: %s", port))
	}

	if info, ok := debug.ReadBuildInfo(); ok {
		for _, setting := range info.Settings {
			if setting.Key == "vcs.revision" {
				log.Print(fmt.Sprintf("Version: %s", setting.Value))
			}
		}
	} else {
		log.Print("Unknown version. Not built from a git repo?")
	}

	log.Print(fmt.Sprintf("Server listening at %s", port))
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", port), router))

}
