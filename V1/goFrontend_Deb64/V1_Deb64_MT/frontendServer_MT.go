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
func check(e error) {
	if e != nil {
		go panic(e)
	}
}

func checkFileNotEmpty(path string) bool {
	fileInfo, err := os.Stat(path)
	if os.IsNotExist(err) || fileInfo.Size() == 0 {
		return false
	}
	return true
}

func read(filePath string) string {
	data, err := os.ReadFile(filePath)
	go check(err)
	return string(data)
}

/*func resMaker(req string) string {
	var res string = ""

}*/

func defineCT(p string) string {
	var split = strings.Split(p, ".")
	var ct string = ""

	if len(split) == 2 {
		ct = "text/" + strings.Split(p, ".")[1]
	} else {
		ct = "text/plain"
	}
	if ct == "text/js" {
		ct = "text/javascript"
	} else if ct == "text/txt" {
		ct = "text/plain"
	} else if ct == "text/jpeg" || ct == "text/jpg" || ct == "text/svg" || ct == "text/img" || ct == "text/png" || ct == "text/ico" {
		ct = "image"
	}

	return ct
}

func handlerFunc(w http.ResponseWriter, r *http.Request) {
	//basic headers
	w.Header().Set("Access-Control-Allow-Origin", "*")
	w.Header().Set("Access-Control-Allow-Methods", "GET")
	w.Header().Set("Access-Control-Allow-Headers", "Content-Type")

	//serve page
	var p string = r.URL.Path
	var page string = "site/user" + p

	var ct string
	go defineCT(p)

	if checkFileNotEmpty(page) {
		w.Header().Set("Content-Type", ct)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, read(page))
	} else if ct == "text/html" {
		w.Header().Set("Content-Type", "text/html")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, read("site/404Page/404NotFound.html"))
	} else {
		w.Header().Set("Content-Type", "text/plain")
		w.WriteHeader(http.StatusNotFound)
		fmt.Fprint(w, "404 Not Found")
	}
}

// main code
func main() {
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
		fmt.Println("Starting HTTPS server on port 1443...")
		tlsConfig := &tls.Config{}
		server := &http.Server{
			Addr:      ":1443",
			Handler:   mux,
			TLSConfig: tlsConfig,
		}
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		// Else, start the http server
		fmt.Println("Starting HTTP server on port 80...")
		log.Fatal(http.ListenAndServe(":80", mux))
	}
}
