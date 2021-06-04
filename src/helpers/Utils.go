package helpers

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"net/http"

	"gopkg.in/yaml.v2"
)

type AppConfig struct {
	WebpLibPath string `yaml:"webpLibPath"`
	Port        string `yaml:"port"`
}

func InitAppParams() AppConfig {
	var configPath string

	flag.StringVar(&configPath, "c", "", "[Required] Path to config")

	flag.Parse()

	config := AppConfig{}

	configFile, err := ioutil.ReadFile(configPath)

	CheckError(err)

	err = yaml.Unmarshal(configFile, &config)

	CheckError(err)

	return config
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
