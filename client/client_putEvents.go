package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"os"

	_ "github.com/go-sql-driver/mysql"
)

var eventIndex = 1

func main() {
	apiUrl := "http://localhost:3232/event"

	data := url.Values{}
	u, _ := url.ParseRequestURI(apiUrl)

	client := &http.Client{}
	events := getEvents() // retrive events from JSON file
	//now send post request per event in the following for loop

	for _, p := range events {
		//retrive necessary parameters to create the url
		data.Set("api_key", p.API_KEY)
		data.Set("unix_timestamp", p.UNIX_TIMESTAMP)
		data.Set("user_id", p.USER_ID)
		u.RawQuery = data.Encode()
		urlStr := fmt.Sprintf("%v", u)
		//now the post request with the created url can be sent to server
		r, _ := http.NewRequest("POST", urlStr, nil)
		resp, err := client.Do(r)
		//check if server is reachable
		if err != nil {
			panic(err)
		}
		//use response variable to check for messages from server side to handle error conditions
		//first handle error conditions, server sends 404 Not Found response in case of errors
		if resp.Status == "404 Not Found" {
			handleErrorConditions(resp)
		} else {
			fmt.Println(eventIndex, "- Event is successfully saved in database.")
		}
		eventIndex = eventIndex + 1
	}
	//fmt.Println(toJson(events)) -- events' definitions can be printed
}

func handleErrorConditions(resp *http.Response) {
	//server puts error conditions in headers, so check the response header
	if resp.Header.Get("Answer") == "input_apiKeyMissing" {
		fmt.Println(eventIndex, "- api key parameter is missing. Please specify the api key in your url.")
	}
	else if resp.Header.Get("Answer") == "input_userIdMissing" {
		fmt.Println(eventIndex, "- user id parameter is missing. Please specify the user id in your url.")
	}
	else if resp.Header.Get("Answer") == "input_timeMissing" {
		fmt.Println(eventIndex, "- unix timestamp parameter is missing. Please specify unix timestamp in your url.")
	}
	else if resp.Header.Get("Answer") == "input_userIdError" {
		fmt.Println(eventIndex, "- User Id value must be integer")
	}
	else if resp.Header.Get("Answer") == "apikeyNotFound" {
		fmt.Println(eventIndex, "- The api key: (", resp.Header.Get("apikey"), ") does not exist in database")
	}
	else if resp.Header.Get("Answer") == "input_TimeError" {
		fmt.Println(eventIndex, "- Unix timestamp value must be integer")
	}
	else if resp.Header.Get("Answer") == "input_apiKeyTypeError" {
		fmt.Println(eventIndex, "- The api key must be integer")
	}
}

// event definition
type Event struct {
	API_KEY        string `json:"api_key"`
	USER_ID        string `json:"user_id"`
	UNIX_TIMESTAMP string `json:"unix_timestamp"`
}

func getEvents() []Event { //retrives all the events from JSON file
	raw, err := ioutil.ReadFile("./events.json")
	if err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	}

	var c []Event
	json.Unmarshal(raw, &c)
	return c
}
