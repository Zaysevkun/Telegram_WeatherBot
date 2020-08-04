package main

import (
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	Token := fetchApiKey()
	Api := "https://api.telegram.org/bot"
	Url := Api + Token
	for {
		updates, err := getUpdates(Url)
		if err != nil {
			log.Println("error in GetUpdates", err.Error())
		}
		fmt.Println(updates)
	}
}

func fetchApiKey() string {
	if err := godotenv.Load("apiKey.env"); err != nil {
		log.Print("No .env file found")
	}
	apiKey, err := os.LookupEnv("BOT_API_KEY")

	if !err {
		log.Println("api key not found")
	}
	return apiKey
}

func getUpdates(Url string) ([]Update, error) {
	resp, err := http.Get(Url + "/getUpdates")
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	var response RestResponse
	err = json.Unmarshal(body, &response)
	if err != nil {
		return nil, err
	}
	return response.Result, nil
}

func respond() {

}
