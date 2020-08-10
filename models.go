package main

type Update struct {
	UpdateId int     `json:"update_id"`
	Message  Message `json:"message"`
}

type Message struct {
	Chat Chat   `json:"chat"`
	Text string `json:"text"`
}

type Chat struct {
	ChatId int `json:"id"`
}

type RestResponse struct {
	Result []Update `json:"result"`
}

type City struct {
	Count int    `json:"count"`
	List  []List `json:"list"`
}

type List struct {
	Id int `json:"id"`
}

type BotMessage struct {
	ChatId int    `json:"chat_id"`
	Text   string `json:"text"`
}

type Weather struct {
	Description string `json:"description"`
	Icon        string `json:"icon"`
}

type Main struct {
	Temp      float64 `json:"temp"`
	FeelsLike float64 `json:"feels_like"`
	Pressure  int     `json:"pressure"`
	Humidity  int     `json:"humidity"`
}

type Wind struct {
	Speed float64 `json:"speed"`
	Deg   int     `json:"deg"`
}

type WeatherResponse struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
	Name    string    `json:"name"`
}

type WeatherResponse5Day struct {
	Weather []Weather `json:"weather"`
	Main    Main      `json:"main"`
	Wind    Wind      `json:"wind"`
	Date    string    `json:"dt_txt"`
}

type ImageMessage struct {
	Id    int    `json:"chat_id"`
	Photo string `json:"photo"`
}

type CityInfo struct {
	Name string `json:"name"`
}

type WeatherList struct {
	List []WeatherResponse5Day `json:"list"`
	City CityInfo              `json:"city"`
}
