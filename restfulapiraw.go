package main

import (
	"fmt"
	"io/ioutil"
	"net"
	"net/http"
)

func main() {
	// Create a new HTTP router
	http.HandleFunc("/upload", uploadHandler)

	// Start the server
	http.ListenAndServe(":8080", nil)
}

func uploadHandler(w http.ResponseWriter, r *http.Request) {

	// Replace "localhost:9100" with the address of your server
	serverAddr := "localhost:9100"

	// Read the request body to get the PDF file data
	fileBytes, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	// Send the PDF data to the print server
	// Connect to the server's TCP port
	conn, err := net.Dial("tcp", serverAddr)
	if err != nil {
		panic(err)
	}

	// Send the PDF file data to the server
	_, err = conn.Write(fileBytes)
	if err != nil {
		panic(err)
	}

	// Close the connection
	err = conn.Close()
	if err != nil {
		panic(err)
	}

	// Respond with a success message
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	fmt.Fprint(w, `{"message": "File uploaded successfully"}`)
}
