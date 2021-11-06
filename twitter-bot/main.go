package main

import (
	"database/sql"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"time"

	_ "github.com/go-sql-driver/mysql"

	"github.com/dghubble/oauth1"
)

type Tweet struct {
	ConversationID string `json:"conversation_id"`
	ID             string
	Text           string
	InReplyToID    string `json:"in_reply_to_status_id"`
}

type TwitterResponse struct {
	Tweets []Tweet `json:"data"`
	Meta   struct {
		OldestID    string
		NewestID    string
		ResultCount int
	}
}

type TwitterClient struct {
	httpClient             *http.Client
	clientUserID           string
	lastRetrievedMentionID string
	db                     *sql.DB
}

func remove(slice []int, s int) []int {
	return append(slice[:s], slice[s+1:]...)
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

	last_mention := 1424882508848435204

	axolobot_user := 1451497427098275860
	url := fmt.Sprintf(
		"https://api.twitter.com/2/users/%v/mentions?max_results=%v&since_id=%v&tweet.fields=conversation_id",
		axolobot_user, number, last_mention)

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
		row := c.db.QueryRow("SELECT count(*) as count from axolobot.mention where mention_id = ?", tweet.ID)
		if err != nil {
			return nil, err
		}
		var count int
		row.Scan(&count)
		if count == 0 {
			_, err := c.db.Query("INSERT INTO mention VALUES (?)", tweet.ID)
			if err != nil {
				return nil, err
			}
			newTweets = append(newTweets, tweet)
		}
	}

	return newTweets, nil

}

func (c *TwitterClient) GetTweetsByConversationID(conversation string, number int) ([]Tweet, error) {

	url := fmt.Sprintf(
		"https://api.twitter.com/2/tweets/search/recent?query=conversation_id:%v%%20-has:media%%20lang:en&max_results=%v&tweet.fields=conversation_id",
		conversation, number)

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

func (c *TwitterClient) PostResponse(tweet Tweet) error {

	params := url.Values{}

	params.Add("status", tweet.Text)
	params.Add("in_reply_to_status_id", tweet.InReplyToID)

	url := "https://api.twitter.com/1.1/statuses/update.json?" + params.Encode()

	_, err := c.makeRequest("POST", url)
	if err != nil {
		log.Fatal("Error when publishing tweet")
		return err
	}

	return nil

}

func NewTwitterClient() *TwitterClient {

	auth_tokens := os.Getenv("AUTH_TOKENS")
	splitted := strings.Split(auth_tokens, ":")
	consumerKey := splitted[0]
	consumerSecret := splitted[1]
	token := splitted[2]
	tokenSecret := splitted[3]

	config := oauth1.NewConfig(consumerKey, consumerSecret)
	oAuthToken := oauth1.NewToken(token, tokenSecret)

	// httpClient will automatically authorize http.Request's
	client := config.Client(oauth1.NoContext, oAuthToken)

	dbPassword := os.Getenv("DB_PASSWORD")
	database, err := sql.Open("mysql", "root:"+dbPassword+"@tcp(db:3306)/axolobot")
	if err != nil {
		return nil
	}

	return &TwitterClient{
		httpClient: client,
		db:         database,
	}

}

func main() {

	twitterClient := NewTwitterClient()

	//tweets, err := twitterClient.GetTweetsByConversationID("1453481950052700161", 10)
	//if err != nil {
	//	log.Fatal(err.Error())
	//}
	//fmt.Println("Un tweet de la conversacion", tweets[0].Text)

	mentionExchanger := make(chan Tweet, 100)

	go MentionListener(mentionExchanger, twitterClient)

	// Spawn a pool of 4 workers
	for w := 0; w < 5; w++ {
		go MentionWorker(mentionExchanger, twitterClient)
	}

	select {}

}

func MentionListener(mentionExchanger chan Tweet, twitterClient *TwitterClient) {

	for {

		log.Println(" 🧏 Listening for new Tweets... ")

		mentions, err := twitterClient.GetNewMentions(10)
		if err != nil {
			log.Fatal("Error consiguiendo menciones: ", err.Error())
		}
		for _, mention := range mentions {
			fmt.Println("Una mencion al bot -> ", mention.Text)
			mentionExchanger <- mention
		}

		time.Sleep(9 * time.Second)
	}

}

func MentionWorker(mentionExchanger chan Tweet, twitterClient *TwitterClient) {

	for {

		mention := <-mentionExchanger
		log.Println(" 🌟 Answering to Tweet -> " + mention.Text)
		response := Tweet{
			InReplyToID: mention.ID,
			Text:        "I am still in development",
		}
		twitterClient.PostResponse(response)

	}

}
