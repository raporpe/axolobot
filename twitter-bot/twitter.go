package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
)

type Tweet struct {
	ConversationID string `json:"conversation_id"`
	ID             string
	Text           string
	InReplyToID    string `json:"in_reply_to_status_id"`
	UserID         string `json:"author_id"`
}

type TwitterResponse struct {
	Tweets []Tweet `json:"data"`
	Meta   struct {
		OldestID    string
		NewestID    string
		ResultCount int
	}
}

type UserLookupResponse struct {
	Data struct {
		Username string `json:"username"`
	}
}

type TwitterClient struct {
	httpClient *http.Client
	db         *DatabaseManager
	hostname   string
}

func (c *TwitterClient) makeRequest(method string, url string) (string, error) {

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}

	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode != http.StatusOK {
		fmt.Println("URL -> " + url)
		fmt.Println("Response -> " + string(responseData))
		return "", errors.New("error de autentiación")
	}

	return string(responseData), nil
}

func (c *TwitterClient) GetNewMentions(number int) ([]Tweet, error) {

	last_mention := c.db.GetLastMentionID()
	log.Println("The last mention is -> " + last_mention)

	axolobotUser := "1451497427098275860"

	params := url.Values{}
	params.Add("max_results", strconv.Itoa(number))
	params.Add("since_id", last_mention)
	params.Add("tweet.fields", "conversation_id")
	params.Add("expansions", "author_id")

	url := c.hostname + "/2/users/" + axolobotUser + "/mentions?" + params.Encode()

	j, err := c.makeRequest("GET", url)
	if err != nil {
		fmt.Println("Error when requesting twitter api: ", err.Error())
		return nil, err
	}

	var twitterResponse TwitterResponse

	err = json.Unmarshal([]byte(j), &twitterResponse)
	if err != nil {
		log.Println("Error convirtiendo json a struct obteniendo menciones")
		return nil, err
	}

	tweets := twitterResponse.Tweets
	newTweets := []Tweet{}

	// Discard the ones in the database and insert the new ones
	for _, tweet := range tweets {
		if !c.db.IsMentionDone(tweet) && tweet.UserID != axolobotUser {
			newTweets = append(newTweets, tweet)
		}
	}

	return newTweets, nil

}

func (c *TwitterClient) GetTweetsByConversationID(conversation string) ([]Tweet, error) {

	params := url.Values{}
	query := fmt.Sprintf("conversation_id:%v -has:media lang:en", conversation)
	params.Add("query", query)
	params.Add("tweet.fields", "conversation_id")
	params.Add("max_results", "100")
	params.Add("expansions", "author_id")

	url := c.hostname + "/2/tweets/search/recent?" + params.Encode()

	j, err := c.makeRequest("GET", url)
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

func (c *TwitterClient) GetUsernameByUserID(userID string) (string, error) {

	j, err := c.makeRequest("GET", c.hostname+"/2/users/"+userID)
	if err != nil {
		return "", err
	}

	var userLookupResponse UserLookupResponse
	json.Unmarshal([]byte(j), &userLookupResponse)
	if err != nil {
		log.Println("Error en la respuesta al extraer username por userID")
		return "", err
	}

	return userLookupResponse.Data.Username, nil

}

func (c *TwitterClient) PostResponse(tweet Tweet) error {

	username, err := c.GetUsernameByUserID(tweet.UserID)
	if err != nil {
		return err
	}

	params := url.Values{}

	params.Add("status", "@"+username+" "+tweet.Text)
	params.Add("in_reply_to_status_id", tweet.InReplyToID)

	url := c.hostname + "/1.1/statuses/update.json?" + params.Encode()
	fmt.Println(url)

	_, err = c.makeRequest("POST", url)
	if err != nil {
		log.Fatal("Error when publishing tweet: " + err.Error())
	}

	return nil

}

func (c *TwitterClient) SetMentionDone(tweet Tweet) error {
	err := c.db.InsertMention(tweet)
	if err != nil {
		return err
	}
	return nil
}
