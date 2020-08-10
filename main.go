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
var currentCityId = 0

func main() {
	botToken, weatherApiKey = fetchApiKey()
	telegramApi := "https://api.telegram.org/bot"
	//getCityApiPart1 := "http://api.openweathermap.org/data/2.5/find?q="
	//getCityApiPart2 := ",RU&type=like&APPID="
	botUrl := telegramApi + botToken
	offset := 0
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
			if ExpectingCity == 11 {
				err := handleYesNo(botUrl, update)
				if err != nil {
					log.Println(err)
				}
			}
			if ExpectingCity == 2 {
				err := get5Days(botUrl, update)
				if err != nil {
					log.Println(err)
				}
			}
			if ExpectingCity == 1 {
				err := getCity(botUrl, update)
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
		msg := "Город России,в котором ты хочешь узнать погоду?(Отмена для выхода)"
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

func getCity(botUrl string, update Update) error {
	cityName := update.Message.Text
	if cityName == "Отмена" || cityName == "отмена" {
		ExpectingCity = 0
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
		ExpectingCity = 11
		currentCityId = city.List[0].Id
		weather, err := getCurrentWeather(currentCityId)
		if err != nil {
			return err
		}
		wtr := weather.Weather[0].Description
		tmp := strconv.FormatFloat(weather.Main.Temp, 'g', -1, 64)
		flslk := strconv.FormatFloat(weather.Main.FeelsLike, 'g', -1, 64)
		prs := strconv.Itoa(weather.Main.Pressure)
		hum := strconv.Itoa(weather.Main.Humidity)
		spd := strconv.FormatFloat(weather.Wind.Speed, 'g', -1, 64)
		deg := computeDirection(weather.Wind.Deg)
		nm := weather.Name
		msg := "Город: " + nm + "\nПогода: " + wtr + "\nТемпература: " + tmp + "С'\nЧувствуется как: " + flslk + "С'\nДавление: " + prs + "мм. рт. ст.\nВлажность: " + hum + "%\nСкорость ветра: " + spd + "м/с\nнаправление ветра: " + deg
		err = sendMessage(botUrl, update, msg)
		if err != nil {
			return err
		}
		err = sendImage(botUrl, update, weather.Weather[0].Icon)
		if err != nil {
			return err
		}
		err = sendMessage(botUrl, update, "Хотите увидеть прогноз на 5 дней?(да/нет)")
		if err != nil {
			return err
		}
		//ExpectingCity = 2
	}
	return nil
}

func assembleCityUrl(cityName string) string {
	getCityApiPart1 := "http://api.openweathermap.org/data/2.5/find?q="
	getCityApiPart2 := ",RU&type=like&APPID="
	cityUrl := getCityApiPart1 + cityName + getCityApiPart2 + weatherApiKey
	return cityUrl
}

func assembleCityFiveUrl(cityName string) string {
	getCityApiPart1 := "http://api.openweathermap.org/data/2.5/forecast?id="
	getCityApiPart2 := "&units=metric&lang=ru&appid="
	cityUrl := getCityApiPart1 + cityName + getCityApiPart2 + weatherApiKey
	return cityUrl
}

func getCurrentWeather(id int) (WeatherResponse, error) {
	apiUrl := "http://api.openweathermap.org/data/2.5/weather?id=" + strconv.Itoa(id) + "&units=metric" + "&lang=ru" + "&appid=" + weatherApiKey
	resp, err := http.Get(apiUrl)
	var response WeatherResponse
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

func sendImage(Url string, update Update, imgCode string) error {
	var botMessage ImageMessage
	botMessage.Id = update.Message.Chat.ChatId
	botMessage.Photo = "http://openweathermap.org/img/wn/" + imgCode + ".png"
	buf, err := json.Marshal(botMessage)
	if err != nil {
		return err
	}
	sendUrl := Url + "/sendPhoto"
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

func get5Days(botUrl string, update Update) error {
	cityId := strconv.Itoa(currentCityId)
	apiUrl := assembleCityFiveUrl(cityId)
	ExpectingCity = 0
	resp, err := http.Get(apiUrl)
	var response WeatherList
	if err != nil {
		return err
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return err
	}
	nm := response.City.Name
	var tmp float64
	var flslk float64
	for i := 0; i < 5; i++ {
		for j := i; j < i+8; j++ {
			tmp += response.List[j].Main.Temp
			flslk += response.List[j].Main.FeelsLike
		}
		tmp := strconv.FormatFloat(tmp/8, 'g', 2, 64)
		flslk := strconv.FormatFloat(flslk/8, 'g', 2, 64)
		wtr := response.List[i*8].Weather[0].Description
		prs := strconv.Itoa(response.List[i*8].Main.Pressure)
		hum := strconv.Itoa(response.List[i*8].Main.Humidity)
		spd := strconv.FormatFloat(response.List[i*8].Wind.Speed, 'g', 2, 64)
		deg := computeDirection(response.List[i*8].Wind.Deg)
		date := response.List[i*8].Date[0:10]
		msg := date + "\nГород: " + nm + "\nПогода: " + wtr + "\nТемпература: " + tmp + "С'\nЧувствуется как: " + flslk + "С'\nДавление: " + prs + "мм. рт. ст.\nВлажность: " + hum + "%\nСкорость ветра: " + spd + "м/с\nнаправление ветра: " + deg
		err = sendMessage(botUrl, update, msg)
		if err != nil {
			return err
		}
	}
	return nil
}

func handleYesNo(botUrl string, update Update) error {
	msg := update.Message.Text
	switch {
	case msg == "да" || msg == "Да":
		ExpectingCity = 2
	case msg == "нет" || msg == "Нет":
		ExpectingCity = 11
	default:
		err := sendMessage(botUrl, update, "Некорректный ввод.Пожалуйста, используйте для ответа 'да' или 'нет'")
		if err != nil {
			return err
		}
	}
	return nil
}

func computeDirection(deg int) string {
	switch {
	case deg > 337:
		return "С"
	case deg > 292:
		return "СЗ"
	case deg > 247:
		return "З"
	case deg > 202:
		return "ЮЗ"
	case deg > 157:
		return "Ю"
	case deg > 122:
		return "ЮВ"
	case deg > 67:
		return "В"
	case deg > 22:
		return "СВ"
	}
	return "С"
}
