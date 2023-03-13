package main

import (
    "bytes"
    "fmt"
    "io/ioutil"
    "mime/multipart"
    "net/http"
    "os"
)

package main

import (
	"bytes"
	"fmt"
	"io"
	"io/ioutil"
	"mime/multipart"
	"net/http"
	"os"
)

func main() {
	// Replace "http://localhost:631/printers/printer-name" with the URL of your CUPS printer.
	printerURL := "http://localhost:631/printers/printer-name"

	// Create an HTTP handler function that handles POST requests to "/print".
	http.HandleFunc("/print", func(w http.ResponseWriter, r *http.Request) {
		// Read the content of the request body.
		body, err := ioutil.ReadAll(r.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusBadRequest)
			return
		}

		// Create a new multipart/form-data request with the PDF file content.
		bodyBuf := &bytes.Buffer{}
		bodyWriter := multipart.NewWriter(bodyBuf)

		fileHeader := make(textproto.MIMEHeader)
		fileHeader.Set("Content-Disposition", fmt.Sprintf(`form-data; name="%s"; filename="%s"`, "file", "file.pdf"))
		fileHeader.Set("Content-Type", "application/pdf")

		fileWriter, err := bodyWriter.CreatePart(fileHeader)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		_, err = io.Copy(fileWriter, bytes.NewReader(body))
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		bodyWriter.Close()

		// Create a new HTTP request with the multipart/form-data body.
		req, err := http.NewRequest("POST", printerURL, bodyBuf)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		req.Header.Set("Content-Type", bodyWriter.FormDataContentType())

		// Submit the HTTP request to the CUPS printer.
		client := &http.Client{}
		resp, err := client.Do(req)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}
		defer resp.Body.Close()

		// Read the response from the CUPS printer.
		respBody, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			http.Error(w, err.Error(), http.StatusInternalServerError)
			return
		}

		// Write the response from the CUPS printer to the HTTP response.
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(resp.StatusCode)
		w.Write(respBody)
	})

	// Start the HTTP server on port 8080.
	fmt.Println("Listening on port 8080...")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		panic(err)
	}
}

