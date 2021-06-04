package main

import (
	"bytes"
	"errors"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dastanaron/webpConvertor/convertor"
	"github.com/dastanaron/webpConvertor/helpers"
)

var appConfig helpers.AppConfig

func main() {
	appConfig = helpers.InitAppParams()
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appConfig.Port), nil))
}

func handleRequest(response http.ResponseWriter, request *http.Request) {
	downloadedImage, err := downloadImage(request)
	if err != nil {
		helpers.BuildErrorResponse(err, response, 422)
		return
	}
	convert(response, request, downloadedImage)
}

func convert(response http.ResponseWriter, request *http.Request, image *bytes.Buffer) {
	params := request.URL.Query()

	cwebp := convertor.NewCWebP()
	cwebp.SetBinPath(appConfig.WebpLibPath).Input(image).Output(response)

	qualityParam, ok := params["q"]

	var err error

	if ok {
		quality, err := strconv.Atoi(qualityParam[0])

		if err == nil {
			cwebp.SetQuality(quality)
		}
	}

	resizeParamWidth, ok := params["w"]

	var width, heigth int

	if ok {
		width, err = strconv.Atoi(resizeParamWidth[0])

		if err != nil {
			helpers.BuildErrorResponse(errors.New("[w] parameter is not string"), response, 422)
			return
		}
	}

	resizeParamHeight, ok := params["h"]

	if ok {
		heigth, err = strconv.Atoi(resizeParamHeight[0])

		if err != nil {
			helpers.BuildErrorResponse(errors.New("[h] parameter is not string"), response, 422)
			return
		}
	}

	if width > 0 && heigth > 0 {
		cwebp.SetResize(width, heigth)
	}

	err = cwebp.Run()

	if err != nil {
		helpers.BuildErrorResponse(err, response, 500)
	}
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
