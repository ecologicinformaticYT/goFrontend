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

// development password (definition)
var devpassword string = "godev65" //default : godev65

// functions
func check(e error) {
	/*
		Function to panic the server in case of error
		Parameters :
			(error) e : the error (nil if there's no error)
	*/

	if e != nil {
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
	check(err)
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

	file, err := os.OpenFile(file_, os.O_CREATE|os.O_APPEND|mode, 0600)
	check(err)
	if _, err := file.WriteString(text); err != nil {
		panic(err)
	}
}

func writef(text_ string, file__ string) {
	/*
		Function to write/create a file into the 'www' directory
		Parameters :
			(string) text : text content of the file
			(string) file_ : path to the file
			(int) mode : os mode to write into the file (flag).
	*/

	write(text_, file__, os.O_WRONLY)
}

func checkDevPass(pass string) int {
	/*
		Function to check if a string is the devlopment password
		Parameters :
			(string) pass : the string
	*/
	if pass == devpassword {
		return 200
	} else {
		return 403
	}

}

func njslog(content string) {
	/*
		Function to log a result into the console
		Parameters :
			(string) content : the string to write into the console
	*/
	fmt.Println(content)
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

	//check if the user is looking for the 'index.html' file
	if p == "" || p == "/" {
		p = "index.html"
		page = "www/index.html"
	}

	var split = strings.Split(p, ".")
	var ct string = ""

	//define content-type
	if len(split) == 2 && r.Method == "GET" {
		ct = "text/" + strings.Split(p, ".")[1]
	} else if r.Method == "GET" {
		ct = "text/plain"
	}
	if ct == "text/js" && r.Method == "GET" {
		ct = "text/javascript"
	} else if ct == "text/txt" && r.Method == "GET" {
		ct = "text/plain"
	} else if (ct == "text/jpeg" || ct == "text/jpg" || ct == "text/svg" || ct == "text/img" || ct == "text/png" || ct == "text/ico" || ct == "text/gif" || ct == "text/pdf") && r.Method == "GET" {
		ct = "image"
	}

	//check if the file exists and isn't empty
	if checkFileNotEmpty(page) && r.Method == "GET" {
		//send the page to the client (HTTP 200) | HTML content
		w.Header().Set("Content-Type", ct)
		w.WriteHeader(http.StatusOK)
		fmt.Fprint(w, read(page))

	} else if ct == "text/html" && r.Method == "GET" {
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

	//--Define the ports number--

	var httpp string = ":65080"
	var httpsp string = ":65443"

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
		fmt.Println("Starting HTTPS sandbox on port " + httpsp)
		tlsConfig := &tls.Config{}
		server := &http.Server{
			Addr:      httpsp,
			Handler:   mux,
			TLSConfig: tlsConfig,
		}
		log.Fatal(server.ListenAndServeTLS(certFile, keyFile))
	} else {
		// Else, start the http server
		fmt.Println("Starting HTTP sandbox on port " + httpp)
		log.Fatal(http.ListenAndServe(httpp, mux))
	}
}
