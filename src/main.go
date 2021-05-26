package main

import (
	"log"
	"net/http"

	"github.com/nickalie/go-webpbin"
)

func main() {
	http.HandleFunc("/", convert)
	log.Fatal(http.ListenAndServe(":8080", nil))
}

func convert(response http.ResponseWriter, request *http.Request) {
	//response.Header().Set("content-type", "image/webp")
	err := webpbin.NewCWebP().
		Quality(80).
		InputFile("../data/test.jpg").
		Output(response).
		Run()
	CheckError(err)
}

func CheckError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
