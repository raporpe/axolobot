package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
)

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
			fmt.Println("😋 Got one mention -> " + mention.Text)
			mentionExchanger <- mention
		}

		time.Sleep(9 * time.Second)
	}

}

func MentionWorker(mentionExchanger chan Tweet, twitterClient *TwitterClient) {

	for {

		// Get the mentions that have to be processed
		mention := <-mentionExchanger

		// Get the tweets in the same conversation as the mention
		tweetsToAnalyze, err := twitterClient.GetTweetsByConversationID(mention.ConversationID)
		if err != nil {
			log.Fatal("Error when getting the tweets to analyze from twitter in conversaionID -> " + mention.ConversationID)
			continue
		}
		// Analyze the tweets using the neural network
		result, err := GetSentimentFromTweets(tweetsToAnalyze)
		if err != nil {
			log.Fatal("Error when passing the tweets to the neural network: " + err.Error())
			continue
		}

		// Different responses depending on the amount of tweets that can be analyzed

		var responseText string
		l := len(tweetsToAnalyze)

		switch {
		case l == 0:
			responseText = "There are no tweets for me to analyse!\n" +
				"Twitter only lets me read tweets published in the last 7 days."
		case l < 10:
			responseText = "Thanks for calling me! The average sentiment of all the responses is " +
				strconv.Itoa(result) + "/100.\n" +
				"Note that I could only analyse " + strconv.Itoa(len(tweetsToAnalyze)) + " tweets."
		default:
			responseText = "Thanks for calling me! The average sentiment of all the responses is " +
				strconv.Itoa(result) + "/100.\nHave a nice day!"
		}

		response := Tweet{
			InReplyToID: mention.ID,
			Text:        responseText,
			UserID:      mention.UserID,
		}
		err = twitterClient.PostResponse(response)
		if err != nil {
			log.Fatal("Could not process mention with id " + mention.ID + ": " + err.Error())
			continue
		}

		err = twitterClient.SetMentionDone(mention)
		if err != nil {
			log.Fatal("Error when inserting done mentions " + err.Error())
			continue
		}

	}

}

func NewTwitterClient() *TwitterClient {

	auth_tokens := os.Getenv("AUTH_TOKENS")
	hostname := "https://api.twitter.com"

	// If in development environment, use mock api
	if auth_tokens == "" {
		log.Println("⚒️ Development mode, using mockup api.")
		return &TwitterClient{
			httpClient: &http.Client{},
			db:         NewDatabaseManager(),
			hostname:   "http://mockup-api:10090",
		}
	}

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
		hostname:   hostname,
	}

}

func GetSentimentFromTweets(tweets []Tweet) (int, error) {

	if len(tweets) == 0 {
		return 0, nil
	}

	averageSentiment := 0

	for _, tweet := range tweets {

		// Delete the '/' that can cause trouble when making queries
		tweet.Text = strings.Replace(tweet.Text, "/", "", -1)
		resp, err := http.Get("http://neural-network:8081/v1/sentiment/" + url.PathEscape(tweet.Text))
		if err != nil {
			return -1, err
		}

		j, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return -1, err
		}
		var data map[string]string
		err = json.Unmarshal(j, &data)

		if err != nil {
			log.Println("Respuesta incorrecta por parte de neural-network")
			return -1, err
		}
		value, err := strconv.Atoi(data["score"])
		if err != nil {
			return -1, err
		}

		averageSentiment += value

	}

	return averageSentiment / len(tweets), nil
}
