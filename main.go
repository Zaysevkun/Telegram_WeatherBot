package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/joho/godotenv"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

func main() {
	botToken, _ := fetchApiKey()
	telegramApi := "https://api.telegram.org/bot"
	getCityApiPart1 := "http://api.openweathermap.org/data/2.5/find?q="
	getCityApiPart2 := ",RU&type=like&APPID="
	botUrl := telegramApi + botToken
	city := "Moscow"
	keyTest := "2e709d8234d5940dadfee59807e51ddd"
	cityUrl := getCityApiPart1 + city + getCityApiPart2 + keyTest
	for {
		//_, err := getUpdates(botUrl)
		updates, err := getUpdates(botUrl)
		if err != nil {
			log.Println("error in GetUpdates", err.Error())
		}
		citySearch, err := findCity(cityUrl)
		if err != nil {
			log.Println("error in GetUpdates", err.Error())
		}
		getForecast(botUrl, updates)
		fmt.Println(updates)
		//for _,update := range updates {
		//
		//}

		fmt.Println(citySearch)
	}
}

func fetchApiKey() (string, string) {
	if err := godotenv.Load("apiKey.env"); err != nil {
		log.Print("No .env file found")
	}
	botApiKey, err := os.LookupEnv("BOT_API_KEY")
	if !err {
		log.Println("bot api key not found")
	}
	weatherApiKey, err := os.LookupEnv("BOT_API_KEY")
	if !err {
		log.Println("weather api key not found")
	}
	return botApiKey, weatherApiKey
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

func findCity(Url string) (City, error) {
	resp, err := http.Get(Url)
	var response City
	if err != nil {
		return response, err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return response, err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}
	return response, nil
}

func getForecast(botUrl string, update []Update) {
	if update[len(update)-1].Message.Text == "/get_forecast" {
		sendMessage(botUrl, update[len(update)-1])
	}
}

func sendMessage(botUrl string, update Update) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = "Напиши мне город,на территории которого хочешь узнать погоду"
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	sendUrl := botUrl + "/sendMessage"
	response, err := http.Post(sendUrl, "application/json", bytes.NewBuffer(buf))

	if err != nil {
		return err
	}
	fmt.Println(response)
	return nil
}
