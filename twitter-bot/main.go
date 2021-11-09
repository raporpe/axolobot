package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
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
}

// Still to implement
func GetClassification(text string) (int, error) {
	resp, err := http.Get("http://neural-network:8081/v1/sentiment/" + url.PathEscape(text))
	if err != nil {
		log.Fatal("Cannot contact with the neural network api")
		return -1, err
	}

	responseData, _ := ioutil.ReadAll(resp.Body)
	var data map[string]int
	err = json.Unmarshal(responseData, &data)
	if err != nil {
		log.Fatal("There was an error reading the response from the neural network api")
		return -1, err
	}

	return data["something"], nil
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

	axolobotUser := 1451497427098275860

	params := url.Values{}
	params.Add("max_results", strconv.Itoa(number))
	params.Add("since_id", last_mention)
	params.Add("tweet.fields", "conversation_id")
	params.Add("expansions", "author_id")

	url := "https://api.twitter.com/2/users/" + strconv.Itoa(axolobotUser) + "/mentions?" + params.Encode()

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

		if !c.db.IsMentionDone(tweet) {
			err := c.db.InsertMention(tweet)
			if err != nil {
				return nil, err
			}
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

	url := "https://api.twitter.com/2/tweets/search/recent?" + params.Encode()

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

	j, err := c.makeRequest("GET", "https://api.twitter.com/2/users/"+userID)
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

	url := "https://api.twitter.com/1.1/statuses/update.json?" + params.Encode()
	fmt.Println(url)

	_, err = c.makeRequest("POST", url)
	if err != nil {
		log.Fatal("Error when publishing tweet: " + err.Error())
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

	return &TwitterClient{
		httpClient: client,
		db:         NewDatabaseManager(),
	}

}

func main() {

	twitterClient := NewTwitterClient()

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

		// Get the tweets in that conversation
		tweetsToAnalyze, err := twitterClient.GetTweetsByConversationID(mention.ConversationID)
		fmt.Println(tweetsToAnalyze)
		if err != nil {
			log.Fatal("Error when getting the tweets to analyze from twitter in conversaionID -> " + mention.ConversationID)
			mentionExchanger <- mention
			continue
		}
		// Analyze the tweets using the neural network
		result, err := AnalyzeTweets(tweetsToAnalyze)
		if err != nil {
			log.Fatal("Error when passing the tweets to the neural network: " + err.Error())
			mentionExchanger <- mention
			continue
		}

		responseText := "Hi there! The result is " + strconv.Itoa(result)

		response := Tweet{
			InReplyToID: mention.ID,
			Text:        responseText,
			UserID:      mention.UserID,
		}
		err = twitterClient.PostResponse(response)
		if err != nil {
			log.Fatal("Could not process mention with id " + mention.ID + ": " + err.Error())
			mentionExchanger <- mention
			continue
		}

	}

}

func AnalyzeTweets(tweets []Tweet) (int, error) {
	return 0, nil
}
