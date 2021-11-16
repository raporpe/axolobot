package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

// Worker funcion that is constantly pulling mentions from Twitter
// The found mentions are passed to the mentionExchanger where several instances of MentionWorker are waiting
func MentionListener(mentionExchanger chan Tweet, twitterClient *TwitterClient) {

	for {

		log.Println(" 🧏 Listening for new Tweets... ")

		mentions, err := twitterClient.GetNewMentions(10)
		if err != nil {
			log.Fatal("Error getting mentions: ", err.Error())
		}

		for _, mention := range mentions {
			log.Println("😋 Got one mention -> " + mention.Text)
			mentionExchanger <- mention
		}

		time.Sleep(9 * time.Second)
	}

}

// Gets the mentions passed by MentionListener trough channel and performs the following steps:
// 1. Get all the tweets in the same conversation as the mention Tweet
// 2. Get the sentiment for all the tweets of step 1
// 3. Post a response with the results of step 2
// 4. Set the mention as done in the database to avoid doing it twice.
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

		// Analyze the tweets using the neural network api
		results, err := GetSentimentFromTweets(tweetsToAnalyze)
		if err != nil {
			log.Fatal("Error when passing the tweets to the neural network: " + err.Error())
			continue
		}

		// There are different responses depending on the amount of tweets that can be analyzed

		var negativeTweets int
		var positiveTweets int

		for _, result := range results {
			if result > 50 {
				positiveTweets++
			} else {
				negativeTweets++
			}
		}

		welcomeMessages := []string{
			"Hi there! 😊",
			"So nice to see you! 😉",
			"Hello!! 🖖",
			"Hi! 💜 ",
			"Greetings! 🧐",
		}

		byeMessages := []string{
			"Bye! 👋",
			"Au revoir! 🥖",
			"Adios! 🤠",
			"See you soon! 🙃",
			"Bye bye! 😺",
		}

		negativeReaction := []string{
			"🙀",
			"😰",
			"😢",
			"😿",
			"😮",
			"🥴",
		}

		positiveReaction := []string{
			"🤙",
			"😄",
			"👍",
			"😁",
			"😺",
			"😃",
		}

		welcomeIndex := rand.Intn(len(welcomeMessages))
		byeIndex := rand.Intn(len(byeMessages))
		negativeIndex := rand.Intn(len(negativeReaction))
		positiveIndex := rand.Intn(len(positiveReaction))

		var responseText string
		l := len(results)

		switch {
		case l == 0:
			responseText = "There are no tweets for me to analyse!\n" +
				"I can only see Tweets published in the last 7 days and written in English!"
		case l < 10:
			if negativeTweets >= positiveTweets {
				responseText = fmt.Sprintf("There are %v negative Tweets out of %v.\n", negativeTweets, l)
			} else {
				responseText = fmt.Sprintf("There are %v positive Tweets out of %v.\n", positiveTweets, l)
			}
			responseText +=
				"I could only analyse " + strconv.Itoa(l) + " tweets. \n" +
					"Notice that I can only see Tweets published in the last 7 days and written in English!"
		default:
			responseText += welcomeMessages[welcomeIndex] + "\n"
			if negativeTweets >= positiveTweets {
				responseText += fmt.Sprintf("%v%% of the tweets are negative! %v \n", int((negativeTweets/l)*100), negativeReaction[negativeIndex])
			} else {
				responseText += fmt.Sprintf("%v%% of the tweets are positive! %v \n", int((positiveTweets/l)*100), positiveReaction[positiveIndex])
			}
			responseText += byeMessages[byeIndex]
		}

		// Make a Tweet struct with the response
		response := Tweet{
			InReplyToID: mention.ID,
			Text:        responseText,
			UserID:      mention.UserID,
		}

		log.Println("Response -> " + responseText)

		err = twitterClient.PostResponse(response)
		if err != nil {
			log.Fatal("Could not process mention with id " + mention.ID + ": " + err.Error())
			continue
		}

		// Store the Tweet ID of the mention in the database to avoid doing it twice
		err = twitterClient.SetMentionDone(mention)
		if err != nil {
			log.Fatal("Error when inserting done mentions " + err.Error())
			continue
		}

	}

}

// Helper function that make an instance of the Twitter Client
// The authentication is managed from this function
func NewTwitterClient() *TwitterClient {

	auth_tokens := os.Getenv("AUTH_TOKENS")
	hostname := "https://api.twitter.com"

	// If in development environment, the auth token is not set
	// Then, use the mockup api that is defined in docker-compose
	if auth_tokens == "" {
		log.Println("⚒️ Development mode, using mockup api.")
		return &TwitterClient{
			httpClient: &http.Client{},
			db:         NewDatabaseManager(),
			hostname:   "http://mockup-api:10090",
		}
	}

	// The authentication tokens are in a single environment variable
	// The 4 secrets are separated with three ":"
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

func GetSentimentFromTweets(tweets []Tweet) ([]int, error) {

	// If no tweets are passed, return zeros
	if len(tweets) == 0 {
		return nil, errors.New("No tweets given to analyze")
	}

	var sentiments []int

	for _, tweet := range tweets {

		// Delete the '/' that can cause trouble when making queries
		tweet.Text = strings.Replace(tweet.Text, "/", "", -1)
		resp, err := http.Get("http://neural-network:8081/v1/sentiment/" + url.PathEscape(tweet.Text))
		if err != nil {
			return nil, err
		}

		j, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, err
		}
		var data map[string]string
		err = json.Unmarshal(j, &data)

		if err != nil {
			log.Println("Incorrect response from neural-network")
			return nil, err
		}
		value, err := strconv.Atoi(data["score"])
		if err != nil {
			return nil, err
		}

		sentiments = append(sentiments, value)

	}

	return sentiments, nil
}
