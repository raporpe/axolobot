package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type Tweet struct {
	ConversationID string `json:"conversation_id"`
	ID             string
	Text           string
	InReplyToID    string `json:"in_reply_to_status_id"`
	UserID         string `json:"author_id"`
	Language       string `json:"lang"`
}

// Response structure from Twitter api v2 endpoint
type TwitterResponse struct {
	Tweets []Tweet `json:"data"`
	Meta   struct {
		OldestID    string
		NewestID    string
		ResultCount int
	}
}

// Response structure when looking for username by userID
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

// Makes an authenticated request to any of the Twitter api v1 and v2 endpoints
func (c *TwitterClient) makeRequest(method string, url string) (string, error) {

	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		return "", err
	}

	// Execute request. It is authrized by httpClient
	resp, err := c.httpClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	// Return error is the response is different from 200
	if resp.StatusCode != http.StatusOK {
		log.Println("URL -> " + url)
		log.Println("Response -> " + string(responseData))
		return "", fmt.Errorf("the response from %v was not 200: %v", url, string(responseData))
	}

	return string(responseData), nil
}

// Pulls the mention timeline from the axolobot twitter username
// Only returns new mentions.
// Note: uses api v2
func (c *TwitterClient) GetNewMentions(number int) ([]Tweet, error) {

	last_mention := c.db.GetLastRegisteredMentionID()
	log.Println("The last mention is -> " + last_mention)

	// The author id of axolobot
	axolobotUser := "1451497427098275860"

	// Set parameters in the query url that are necessary
	params := url.Values{}
	params.Add("max_results", strconv.Itoa(number))
	params.Add("since_id", last_mention)
	params.Add("tweet.fields", "conversation_id,author_id,lang")

	url := c.hostname + "/2/users/" + axolobotUser + "/mentions?" + params.Encode()

	j, err := c.makeRequest("GET", url)
	if err != nil {
		fmt.Println("Error when requesting twitter api: ", err.Error())
		return nil, err
	}

	var twitterResponse TwitterResponse

	err = json.Unmarshal([]byte(j), &twitterResponse)
	if err != nil {
		log.Println("Error when convertion mention list to tweet strucut")
		return nil, err
	}

	tweets := twitterResponse.Tweets
	var newTweets []Tweet

	// Check if the mentions shall be processed
	for _, tweet := range tweets {

		// To avoid processing twice the same mention
		mentionNotRegistered := !c.db.IsMentionRegistered(tweet)

		// To avoid processing a mention from the own bot and entering a loop
		mentionNotFromBot := tweet.UserID != axolobotUser

		// To avoid answering to conversations in which the bot is not directly mentioned
		mentionContainsBotName := strings.Contains(tweet.Text, "axolobot")

		if mentionNotRegistered && mentionNotFromBot && mentionContainsBotName {
			newTweets = append(newTweets, tweet)
			c.db.RegisterMention(tweet)
		}
	}

	return newTweets, nil

}

// Gets up to 100 Tweets given a conversation ID
func (c *TwitterClient) GetTweetsByConversationID(conversation string) ([]Tweet, error) {

	params := url.Values{}
	query := fmt.Sprintf("conversation_id:%v -has:media (lang:en OR lang:es)", conversation)
	params.Add("query", query)
	params.Add("tweet.fields", "conversation_id,author_id,lang")
	params.Add("max_results", "100")

	url := c.hostname + "/2/tweets/search/recent?" + params.Encode()

	j, err := c.makeRequest("GET", url)
	if err != nil {
		log.Println("Error when retrieving tweets by conversation id")
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

// Returns username given userID (or authorID)
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

func (c *TwitterClient) GetTweetByID(tweetID string) (Tweet, error) {

	params := url.Values{}
	params.Add("tweet.fields", "conversation_id,author_id,lang")

	url := c.hostname + "/2/tweets/" + tweetID + "?" + params.Encode()

	j, err := c.makeRequest("GET", url)

	var data map[string]Tweet
	json.Unmarshal([]byte(j), &data)
	if err != nil {
		log.Println("Error en la respuesta al extraer username por userID")
		return Tweet{}, err
	}

	return data["data"], nil
}

// Posts a Tweet in response to the given userID and
// Note: uses api v1.1
func (c *TwitterClient) PostResponse(tweet Tweet) error {

	username, err := c.GetUsernameByUserID(tweet.UserID)
	if err != nil {
		return err
	}

	params := url.Values{}

	params.Add("status", "@"+username+" "+tweet.Text)
	params.Add("in_reply_to_status_id", tweet.InReplyToID)

	url := c.hostname + "/1.1/statuses/update.json?" + params.Encode()

	_, err = c.makeRequest("POST", url)
	if err != nil {
		log.Println("Error when publishing tweet: " + err.Error())
	}

	return nil

}

// Stores in the datbase the Tweet ID of the passed Tweet
// When pulling new mentions, those with a already registered TweetID will be ignored
func (c *TwitterClient) SetMentionDone(tweet Tweet) error {
	err := c.db.setMentionDone(tweet)
	if err != nil {
		return err
	}
	return nil
}
