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
	"strconv"
)

var ExpectingCity int = 0

//var CityName = ""
var botToken = ""
var weatherApiKey = ""

func main() {
	botToken, weatherApiKey = fetchApiKey()
	telegramApi := "https://api.telegram.org/bot"
	//getCityApiPart1 := "http://api.openweathermap.org/data/2.5/find?q="
	//getCityApiPart2 := ",RU&type=like&APPID="
	botUrl := telegramApi + botToken
	offset := 0
	//keyTest := "2e709d8234d5940dadfee59807e51ddd"
	//cityUrl := getCityApiPart1 + CityName + getCityApiPart2 + weatherApiKey
	for {
		//_, err := getUpdates(botUrl)
		updates, err := getUpdates(botUrl, offset)
		if err != nil {
			log.Println("error in GetUpdates", err.Error())
		}
		//citySearch, err := findCity(cityUrl)
		//if err != nil {
		//	log.Println("error in find in GetUpdates", err.Error())
		//}
		fmt.Println(updates)
		for _, update := range updates {
			if ExpectingCity == 1 {
				err := getCity(update)
				if err != nil {
					log.Println(err)
				}
			}
			err := getForecast(botUrl, update)
			if err != nil {
				log.Println("error in getting command signal")
			}
			offset = update.UpdateId + 1
		}

		//fmt.Println(citySearch)
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
	weatherApiKey, err := os.LookupEnv("WEATHER_API_KEY")
	if !err {
		log.Println("weather api key not found")
	}
	return botApiKey, weatherApiKey
}

func getUpdates(Url string, offset int) ([]Update, error) {
	resp, err := http.Get(Url + "/getUpdates" + "?offset=" + strconv.Itoa(offset))
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

func getForecast(botUrl string, update Update) error {
	if update.Message.Text == "/get_forecast" {
		ExpectingCity = 1
		msg := "Напиши мне населенный пункт,на территории которого хочешь узнать погоду"
		err := sendMessage(botUrl, update, msg)
		if err != nil {
			return err
		}

	}
	return nil
}

func sendMessage(botUrl string, update Update, message string) error {
	var botMessage BotMessage
	botMessage.ChatId = update.Message.Chat.ChatId
	botMessage.Text = message
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	sendUrl := botUrl + "/sendMessage"
	response, err := http.Post(sendUrl, "application/json", bytes.NewBuffer(buf))
	if err != nil {
		return err
	}
	fmt.Printf("res.StatusCode: %d\n", response.StatusCode)
	if response.StatusCode != 200 {
		log.Fatal(response.Status)
	}
	return nil
}

func getCity(update Update) error {
	cityName := update.Message.Text
	if cityName == "Отмена" || cityName == "отмена" {
		return nil
	}
	Url := assembleCityUrl(cityName)
	var city City
	city, err := findCity(Url)
	if err != nil {
		return err
	}
	if city.Count == 0 {
		err := sendMessage(botToken, update, "населенный пункт введен неправильно, попробуйте еще раз")
		if err != nil {
			return err
		}
	} else {

	}
	return nil
}

func assembleCityUrl(cityName string) string {
	getCityApiPart1 := "http://api.openweathermap.org/data/2.5/find?q="
	getCityApiPart2 := ",RU&type=like&APPID="
	cityUrl := getCityApiPart1 + cityName + getCityApiPart2 + weatherApiKey
	return cityUrl
}
