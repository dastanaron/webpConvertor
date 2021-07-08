package main

import (
	"bytes"
	"errors"
	"fmt"
	"image"
	"image/jpeg"
	_ "image/png"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/dastanaron/webpConvertor/convertor"
	"github.com/dastanaron/webpConvertor/helpers"
	"github.com/disintegration/imaging"
)

var appConfig helpers.AppConfig

func main() {
	appConfig = helpers.InitAppParams()
	http.HandleFunc("/", handleRequest)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%s", appConfig.Port), nil))
}

func handleRequest(response http.ResponseWriter, request *http.Request) {
	cwebp := convertor.NewCWebP()
	cwebp.SetBinPath(appConfig.WebpLibPath)

	err, resizeParameters, quality := buildConvertParameters(request)

	if err != nil {
		helpers.BuildErrorResponse(err, response, 422)
		return
	}

	cwebp.SetQuality(*quality)
	cwebp.Mode(appConfig.Mode)

	downloadedImage, err := downloadImage(request)

	if err != nil {
		helpers.BuildErrorResponse(err, response, 422)
		return
	}
	convert(response, *resizeParameters, downloadedImage, *cwebp)
}

func convert(response http.ResponseWriter, resizeParameters convertor.ResizeParameters, imageBytes *bytes.Buffer, cwebp convertor.WebP) {

	//computedResizeParameters, cropParameters, err := computeFitCropping(resizeParameters, image)

	dstImage, _, err := image.Decode(imageBytes)

	if err != nil {
		helpers.BuildErrorResponse(errors.New("Cannot decode image"), response, 422)
		return
	}

	imageBuf := bytes.Buffer{}

	var dstImageResized image.Image

	if resizeParameters.Type == "fill" {
		dstImageResized = imaging.Fill(dstImage, resizeParameters.Width, resizeParameters.Height, imaging.Center, imaging.Lanczos)
	}

	if resizeParameters.Type == "fit" {
		dstImageResized = imaging.Fit(dstImage, resizeParameters.Width, resizeParameters.Height, imaging.Lanczos)
	}

	if resizeParameters.Type == "" {
		dstImageResized = dstImage
		cwebp.SetResize(resizeParameters)
	}

	err = jpeg.Encode(&imageBuf, dstImageResized, nil)

	if err != nil {
		helpers.BuildErrorResponse(errors.New("Cannot encode image"), response, 422)
		return
	}

	cwebp.Input(&imageBuf).Output(response)

	err = cwebp.Run()

	if err != nil {
		helpers.BuildErrorResponse(err, response, 500)
	}
}

func buildConvertParameters(request *http.Request) (error, *convertor.ResizeParameters, *int) {
	params := request.URL.Query()

	qualityParam, ok := params["q"]

	var err error

	quality := 80

	if ok {
		formattedQuality, err := strconv.Atoi(qualityParam[0])

		if err == nil {
			quality = formattedQuality
		}
	}

	resizeParameters := convertor.ResizeParameters{}
	resizeParamWidth, ok := params["w"]

	if ok {
		resizeParameters.Width, err = strconv.Atoi(resizeParamWidth[0])

		if err != nil {
			return errors.New("[w] parameter is not string"), nil, nil
		}
	}

	resizeParamHeight, ok := params["h"]

	if ok {
		resizeParameters.Height, err = strconv.Atoi(resizeParamHeight[0])

		if err != nil {
			return errors.New("[h] parameter is not string"), nil, nil
		}
	}

	resizeParamType, ok := params["type"]

	if ok {
		resizeParameters.Type = resizeParamType[0]
	}

	if (resizeParameters.Type != "") && (resizeParameters.Height == 0 || resizeParameters.Width == 0) {
		return errors.New("You need to specify the height and width if the type is specified"), nil, nil
	}

	return nil, &resizeParameters, &quality
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
