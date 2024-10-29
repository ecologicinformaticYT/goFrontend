package main

//imports
import (
	"crypto/tls"
	"fmt"
	"log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

// functions
func check(e error, panic_ bool) {
	/*
		Function to log the errors into the console / panic the server in case of error
		Parameters :
			(error) e : the error (nil if there's no error)
	*/

	if e != nil && !panic_ {
		fmt.Println(e)
	} else if e != nil {
		panic(e)
	}
}

func checkFileNotEmpty(path string) bool {
	/*
		Function to check if a file exists and isn't empty
		Parameters :
			(string) path : path to the file
		Returns :
			(bool) : True if the file exists and isn't empty, False otherwise
	*/

	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		return false
	}
	return true
}

func read(filePath string) string {
	/*
		Function to read a file
		Parameters :
			(string) filePath : path to the file
		Returns :
			(string) : the content of the file
	*/

	data, err := os.ReadFile(filePath)
	check(err, false)
	return string(data)
}

func write(text string, file_ string, mode int) {
	/*
		Function to write/create a file
		Parameters :
			(string) text : text content of the file
			(string) file_ : path to the file
			(int) mode : os mode to write into the file (flag).
	*/

	file, err := os.OpenFile(file_, os.O_CREATE|os.O_WRONLY|mode, 0600)
	check(err, false)
	if _, err := file.WriteString(text); err != nil {
		panic(err)
	}
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	/*
		Function to handle the received http (GET) requests
		Parameters :
			(http.ResponseWriter) w : the response writer to write the response
			(http.Request) *r : the request that has been received
	*/

	//set headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//define page location on the server
	var p string = r.URL.Path
	var page string = "www/" + p
	var bpage string = "wwwbackup/" + p

	//check if the user is looking for the 'index.html' file
	if p == "" || p == "/" {
		p = "index.html"
		page = "www/index.html"
	}

	var split = strings.Split(p, ".")
	var ct string = ""

	//define content-type
	if len(split) == 2 {
		ct = "text/" + strings.Split(p, ".")[1]
	} else {
		ct = "text/plain"
	}
	if ct == "text/js" {
		ct = "text/javascript"
	} else if ct == "text/txt" {
		ct = "text/plain"
	} else if ct == "text/jpeg" || ct == "text/jpg" || ct == "text/svg" || ct == "text/img" || ct == "text/png" || ct == "text/ico" || ct == "text/gif" || ct == "text/pdf" {
		ct = "image"
	}

	//check if the file exists and isn't empty
	if checkFileNotEmpty(page) {
		//send the page to the client (HTTP 200) | HTML content
		w.Header().Set("Content-Type", ct)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, read(page))

	} else if !checkFileNotEmpty("./www/index.html") && checkFileNotEmpty(bpage) {
		//send the page to the client (HTTP 200) | HTML content | backup mode
		w.Header().Set("Content-Type", ct)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, read(bpage))

	} else if ct == "text/html" {
		//send the 404 error page to the client (HTTP 404) | HTML content
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, read("404Page/404NotFound.html"))
	} else {
		//send a 404 error text to the client (HTTP 404) | other than HTML content
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Not Found")
	}

}

// main code
func main() {
	/*
		Main function. Starts the http or https server on the HTTP or HTTPS port
	*/

	//--Read the ports number--

	var httpp string = ":80"
	var httpsp string = ":443"

	//--Start the server--
	// Folder where TLS certificates are stored
	certDir := "https"

	// Paths to cert.pem and key.pem (definition)
	certFile := filepath.Join(certDir, "cert.pem")
	keyFile := filepath.Join(certDir, "key.pem")

	// Check if cert.pem and key.pem exist and are not empty
	certExists := checkFileNotEmpty(certFile)
	keyExists := checkFileNotEmpty(keyFile)

	// Create a router for http requests
	mux := http.NewServeMux()
	mux.HandleFunc("/", handlerFunc)

	// If TLS certificates exist and are not empty, configure the HTTPS server
	if certExists && keyExists {
		fmt.Println("Starting HTTPS server on port " + httpsp)
		tlsConfig := &tls.Config{}
		server := &http.Server{
			Addr:      httpsp,
			Handler:   mux,
			TLSConfig: tlsConfig,
		}
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		// Else, start the http server
		fmt.Println("Starting HTTP server on port " + httpp)
		log.Fatal(http.ListenAndServe(httpp, mux))
	}
}
