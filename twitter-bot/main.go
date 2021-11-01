package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Tweet struct {
	ConversationID string `json:"conversation_id"`
	ID             string
	Text           string
}

type TwitterResponse struct {
	Tweets []Tweet `json:"data"`
	Meta   struct {
		OldestID    string
		NewestID    string
		ResultCount int
	}
}

func getJson(url string) (string, error) {
	bearer_token := os.Getenv("BEARER_TOKEN")

	client := http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return "", err
	}

	req.Header.Add("Authorization", "Bearer "+bearer_token)
	resp, err := client.Do(req)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println(resp.Request.URL)
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func getNewMentions(number int) ([]Tweet, error) {

	last_mention := 1424882508848435204

	axolobot_user := 1451497427098275860
	url := fmt.Sprintf(
		"https://api.twitter.com/2/users/%v/mentions?max_results=%v&since_id=%v&tweet.fields=conversation_id",
		axolobot_user, number, last_mention)

	j, err := getJson(url)
	if err != nil {
		fmt.Println("Error when requesting twitter api: ", err.Error())
		return nil, err
	}

	var tr TwitterResponse

	err = json.Unmarshal([]byte(j), &tr)
	if err != nil {
		log.Println("Error convirtiendo json a struct obteniendo menciones")
		return nil, err
	}

	return tr.Tweets, nil

}

func getTweetsByConversationID(conversation string, number int) ([]Tweet, error) {

	url := fmt.Sprintf(
		"https://api.twitter.com/2/tweets/search/recent?query=conversation_id:%v%%20-has:media%%20lang:en&max_results=%v&tweet.fields=conversation_id",
		conversation, number)

	j, err := getJson(url)
	if err != nil {
		log.Fatal("Error when retrieving tweets by conversation id")
		return nil, err
	}

	var tr TwitterResponse

	err = json.Unmarshal([]byte(j), &tr)
	if err != nil {
		log.Println("Error convirtiendo json a struct obteniendo conversaciones")
		return nil, err
	}

	return tr.Tweets, nil
}

func main() {
	log.Println("The twitter bot is programmed in Go")
	mentions, err := getNewMentions(10)
	if err != nil {
		log.Fatal("Error consiguiendo menciones: ", err.Error())
	}

	fmt.Println("Una mencion al bot", mentions[0].Text)

	tweets, err := getTweetsByConversationID("1453481950052700161", 10)
	if err != nil {
		log.Fatal(err.Error())
	}
	fmt.Println("Un tweet de la conversacion", tweets[5].Text)
}
