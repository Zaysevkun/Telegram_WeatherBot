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
	ChatId int `json:"chat"`
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
