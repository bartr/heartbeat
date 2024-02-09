package main

import (
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"sort"
	"strconv"
	"strings"
	"time"
)

// default values
var title = "Heartbeat"
var uri = "/heartbeat"
var root = "/"
var port = 8080
var buffer = "0123456789ABCDEF"
var minSize = 1
var maxSize = 1024 * 1024
var bufferSize = 1
var logResults = false

// version gets set in the build / dockerfile
var Version = "0.4.0"

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

	// handle /heartbeat
	root = uri
	if root != "/" {
		root = strings.TrimRight(uri, "/")

		http.Handle(root, http.HandlerFunc(heartbeatHandler))
	}

	// handle /heartbeat
	http.Handle(uri, http.HandlerFunc(heartbeatHandler))
}

// display config
func displayConfig() {
	log.Println("Version:    ", Version)
	log.Println("URI:        ", uri)
	log.Println("Root:       ", root)
	log.Println("Port:       ", port)
	log.Println("Buffer Size ", bufferSize/1024, "KB")
	log.Println("MinResult:  ", minSize)
	log.Println("MaxResult:  ", maxSize)
	log.Println("Log Results ", logResults)
}

// parseCommandLine
func parseCommandLine() {
	// get env vars
	if s := os.Getenv("URI"); s != "" {
		uri = s
	}

	if s := os.Getenv("TITLE"); s != "" {
		title = s
	}

	// parse flags
	u := flag.String("u", uri, "URI to listen on")
	p := flag.Int("p", port, "port to listen on")
	min := flag.Int("min", minSize, "minimum response size")
	max := flag.Int("max", maxSize, "maximum response size")
	l := flag.Bool("log", logResults, "log incoming requests")
	v := flag.Bool("v", false, "display version")

	flag.Parse()

	// add  trailing /
	if !strings.HasSuffix(*u, "/") {
		*u += "/"
	}

	// add  leading /
	if !strings.HasPrefix(*u, "/") {
		*u = "/" + *u
	}

	// check port
	if *p <= 0 || *p >= 64*1024 {
		flag.Usage()
		log.Fatal("invalid port")
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

// handle root
func rootHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	w.Header().Add("Cache-Control", "no-cache")

	envMap := make(map[string]string)

	// Retrieve environment variables
	envVars := os.Environ()

	// Iterate through environment variables
	for _, envVar := range envVars {
		if strings.HasPrefix(envVar, "e_") || true {
			parts := strings.SplitN(envVar, "=", 2)
			key := strings.Replace(strings.Replace(parts[0], "e_", "", 1), "_", " ", -1)
			value := parts[1]
			envMap[key] = value
		}
	}

	// Sort the keys
	var keys []string
	for key := range envMap {
		keys = append(keys, key)
	}
	sort.Strings(keys)

	// Format map as HTML table
	table := ""

	for _, key := range keys {
		value := envMap[key]
		table += fmt.Sprintf("                                    <tr><td>%s</td><td>%s</td></tr>\n", key, value)
	}

	html := getTemplate()

	html = strings.Replace(html, "{{url}}", uri, -1)
	html = strings.Replace(html, "{{title}}", title, -1)
	html = strings.Replace(html, "{{table}}", strings.TrimRight(table, "\n"), -1)

	fmt.Fprintln(w, html)

	logToConsole(http.StatusOK, root, time.Since(start))
}

// handle /healthz
func healthzHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	w.Header().Add("Cache-Control", "no-cache")
	fmt.Fprintf(w, "pass")

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}

// handle /readyz
func readyzHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	w.Header().Add("Cache-Control", "no-cache")
	fmt.Fprintf(w, "ready")

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}

// handle /version
func versionHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	w.Header().Add("Cache-Control", "no-cache")
	fmt.Fprintf(w, Version)

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}

// handle /heartbeat (or -u)
func heartbeatHandler(w http.ResponseWriter, r *http.Request) {
	// start the request timer
	start := time.Now()

	path := strings.ToLower(r.URL.Path)

	if path == root {
		rootHandler(w, r)
		return
	}

	path = path[len(uri):]

	// handle /heartbeat/healthz
	if path == "healthz" {
		healthzHandler(w, r)
		return
	}

	if path == "readyz" {
		readyzHandler(w, r)
		return
	}

	if path == "version" {
		versionHandler(w, r)
		return
	}

	// get size from URI and convert to int
	size, err := strconv.Atoi(path)

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
			fmt.Fprintln(w, buffer)
		}
	}

	// send the remaining data
	fmt.Fprintln(w, buffer[0:size%bufferSize])

	logToConsole(http.StatusOK, r.URL.Path, time.Since(start))
}

func getTemplate() string {
	return `<!DOCTYPE html>
<html lang="en">
	<head>
		<meta charset="UTF-8">
		<meta name="viewport" content="width=device-width, initial-scale=1.0">
		<title>{{title}}</title>
		<style>
			body {
				font-family: Arial, sans-serif;
				margin: 0;
				padding: 0;
				background-color: #f4f4f4;
			}

			.container {
				max-width: 900px;
				margin: 20px auto;
				padding: 20px;
				background-color: #fff;
				border-radius: 5px;
				box-shadow: 0 0 10px rgba(0, 0, 0, 0.1);
			}

			h1 {
				text-align: center;
			}

			table {
				width: 100%%;
				border-collapse: collapse;
			}

			th, td {
				padding: 10px;
				border-bottom: 1px solid #ddd;
				text-align: left;
			}

			th {
				background-color: #f2f2f2;
			}

			tr:hover {
				background-color: #f5f5f5;
			}
		</style>
	</head>
	<body>
		<div class="container">
			<h1>{{title}}</h1>

			<h2>Environment Variables</h2>

			<table>
				<thead>
				<tr>
					<th>Key</th>
					<th>Value</th>
				</tr>
</thead>
				<tbody>
{{table}}
				</tbody>
			</table>
		</div>
	</body>
</html>
`
}
