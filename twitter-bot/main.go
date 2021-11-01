package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
)

type Data struct {
	ConversationID string
	ID             string
	Text           string
}

type TwitterResponse struct {
	Data []Data
	Meta struct {
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

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	return string(responseData), nil
}

func getNewMentions(number int) ([]Data, error) {

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

	fmt.Println(j)

	var tr TwitterResponse

	err = json.Unmarshal([]byte(j), &tr)
	if err != nil {
		log.Println("Error convirtiendo json a struct")
		return nil, err
	}

	log.Println(tr)

	return tr.Data, nil

}

func main() {
	log.Println("The twitter bot is programmed in Go")
	mentions, err := getNewMentions(10)
	if err != nil {
		log.Fatal("Error consiguiendo menciones: ", err.Error())
	}

	fmt.Println(mentions)
}
