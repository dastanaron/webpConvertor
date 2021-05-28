package main

import (
	"bytes"
	"encoding/json"
	"errors"
	"io"
	"log"
	"net/http"

	"github.com/nickalie/go-webpbin"
)

func main() {
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func handleRequest(response http.ResponseWriter, request *http.Request) {
	downloadedImage, err := downloadImage(request)
	if err != nil {
		BuildErrorResponse(err, response, 422)
		return
	}
	convert(response, request, downloadedImage)
}

func convert(response http.ResponseWriter, request *http.Request, image *bytes.Buffer) {
	err := webpbin.NewCWebP().
		Quality(80).
		Input(image).
		Output(response).
		Run()
	CheckError(err)
}

func downloadImage(request *http.Request) (*bytes.Buffer, error) {
	params := request.URL.Query()

	srcParam, ok := params["src"]

	if !ok || len(srcParam) < 1 {
		return nil, errors.New("Not found source of image")
	}

	imageUrl := srcParam[0]

	response, err := http.Get(imageUrl)

	if err != nil {
		return nil, err
	}

	defer response.Body.Close()

	if response.StatusCode != 200 {
		return nil, errors.New("Received non 200 response code")
	}

	var downloadedFile bytes.Buffer

	_, err = io.Copy(&downloadedFile, response.Body)

	if err != nil {
		return nil, err
	}

	return &downloadedFile, nil
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func BuildErrorResponse(err error, response http.ResponseWriter, statusCode int) {
	errorStructure := ErrorResponse{"error", err.Error()}

	json, err := json.Marshal(errorStructure)
	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}

	response.Header().Set("Content-Type", "application/json")
	response.WriteHeader(statusCode)
	_, err = response.Write(json)

	if err != nil {
		http.Error(response, err.Error(), http.StatusInternalServerError)
		return
	}
}

type ErrorResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
