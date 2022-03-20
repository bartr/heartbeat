package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"
)

// default values
var uri = "/heartbeat/"
var port = 8080
var buffer = "0123456789ABCDEF"
var minSize = 1
var maxSize = 1024 * 1024
var bufferSize = 1
var logResults = false

// version gets set in the build / dockerfile
var Version = "dev"

// main app
func main() {

	parseCommandLine()

	setupHandlers()

	displayConfig()

	// run the web server
	if err := http.ListenAndServe(":"+strconv.Itoa(port), nil); err != nil {
		log.Fatal(err)
	}
}

// setup handlers
func setupHandlers() {
	// build the result buffer
	buffer = strings.Repeat(buffer, bufferSize*1024/16)
	bufferSize = len(buffer)

	// handle /heartbeat/
	http.Handle(uri, http.HandlerFunc(heartbeatHandler))

	// handle /healthz
	http.Handle("/healthz", http.HandlerFunc(healthzHandler))

	// handle /readyz
	http.Handle("/readyz", http.HandlerFunc(readyzHandler))

	// handle /version
	http.Handle("/version", http.HandlerFunc(versionHandler))

	// handle / and /index.htm[l]
	http.Handle("/", http.HandlerFunc(rootHandler))
}

// display config
func displayConfig() {
	log.Println("Version:    ", Version)
	log.Println("URI:        ", uri)
	log.Println("Port:       ", port)
	log.Println("Buffer Size ", bufferSize/1024, "KB")
	log.Println("MinResult:  ", minSize)
	log.Println("MaxResult:  ", maxSize)
	log.Println("Log Results ", logResults)
}

// parseCommandLine
func parseCommandLine() {
	// parse flags
	u := flag.String("u", uri, "URI to listen on")
	p := flag.Int("p", port, "port to listen on")
	min := flag.Int("min", minSize, "minimum response size")
	max := flag.Int("max", maxSize, "maximum response size")
	l := flag.Bool("log", logResults, "log incoming requests")
	v := flag.Bool("v", false, "display version")
	b := flag.Int("b", bufferSize, "buffer size (KB) (1-1024)")

	flag.Parse()

	// add  trailing /
	if !strings.HasSuffix(*u, "/") {
		*u += "/"
	}

	// check URI
	if !strings.HasPrefix(*u, "/") || len(*u) < 3 {
		flag.Usage()
		log.Fatal("invalid URI")
	}

	// check port
	if *p <= 0 || *p >= 64*1024 {
		flag.Usage()
		log.Fatal("invalid port")
	}

	// check buffer size
	if *b < 1 || *b > 1024 {
		flag.Usage()
		log.Fatal("buffer size must be between 1 and 1024")
	}

	// check min
	if *min < 1 || *min >= *max {
		flag.Usage()
		log.Fatal("min must be > 0 and <= max")
	}

	// check max
	if *max > 1024*1024*1024 {
		flag.Usage()
		log.Fatal("max must be <= 1GB")
	}

	// display version and exit
	if *v {
		fmt.Println(Version)
		os.Exit(0)
	}

	// set variables to args (or defaults)
	uri = *u
	port = *p
	minSize = *min
	maxSize = *max
	logResults = *l
}

// very basic logging
func logToConsole(code int, path string, duration time.Duration) {
	if logResults {
		log.Println(code, "\t", duration, "\t", path)
	}
}

// handle /  /index.*  and /default.*
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	if r.URL.Path == "/" || strings.HasPrefix(r.URL.Path, "/index.") || strings.HasPrefix(r.URL.Path, "/default.") {

		http.Redirect(w, r, uri+"17", http.StatusMovedPermanently)

		logToConsole(http.StatusMovedPermanently, r.URL.Path, time.Since(start))
	} else {
		logToConsole(http.StatusNotFound, r.URL.Path, time.Since(start))
		w.WriteHeader(http.StatusNotFound)
	}
}

// handle /healthz
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	w.Header().Add("Cache-Control", "no-cache")
	fmt.Fprintf(w, "Pass\n")

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}

// handle /readyz
func readyzHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	w.Header().Add("Cache-Control", "no-cache")
	fmt.Fprintf(w, "Ready\n")

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}

// handle /version
func versionHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	w.Header().Add("Cache-Control", "no-cache")
	fmt.Fprintf(w, Version+"\n")

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}

// handle /heartbeat (or -u)
func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	// get size from URI and convert to int
	size, err := strconv.Atoi(strings.ToLower(r.URL.Path)[len(uri):])

	// size is invalid
	if err != nil || size < minSize || size > maxSize {
		w.WriteHeader(http.StatusBadRequest)

		logToConsole(http.StatusBadRequest, r.URL.Path, time.Since(start))

		return
	}

	// set no-cache
	w.Header().Add("Cache-Control", "no-cache")

	// send the data in bufferSize chunks
	if size >= bufferSize {
		for i := 0; i < size/bufferSize; i++ {
			fmt.Fprintf(w, buffer)
		}
	}

	// send the remaining data
	fmt.Fprintf(w, buffer[0:size%bufferSize])

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}
