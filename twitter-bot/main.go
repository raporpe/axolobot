package main

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/dghubble/oauth1"
)

const spanishLang = "es"
const englishLang = "en"
const undefinedLang = "und"

func main() {

	twitterClient := NewTwitterClient()

	mentionExchanger := make(chan Tweet, 100)

	go MentionListener(mentionExchanger, twitterClient)

	// Spawn a pool of 4 workers
	for w := 0; w < 5; w++ {
		go MentionWorker(mentionExchanger, twitterClient)
	}

	// Sleep for ever
	select {}

}

// Worker funcion that is constantly pulling mentions from Twitter
// The found mentions are passed to the mentionExchanger where several instances of MentionWorker are waiting
func MentionListener(mentionExchanger chan Tweet, twitterClient *TwitterClient) {

	for {

		log.Println(" 🧏 Listening for new Tweets... ")

		mentions, err := twitterClient.GetNewMentions(10)
		if err != nil {
			log.Println("Error getting mentions: ", err.Error())
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
		lang := mention.Language

		// If the mention contains no text or the mention is in a language different from english or spanish,
		// set the language to english
		if lang == undefinedLang ||
			(lang != spanishLang && lang != englishLang) {
			tweet, err := twitterClient.GetTweetByID(mention.ConversationID)

			if err != nil {
				lang = englishLang
			} else if tweet.Language == englishLang || tweet.Language == spanishLang {
				lang = tweet.Language
			} else {
				lang = englishLang
			}

		}

		// Get the tweets in the same conversation as the mention
		tweetsToAnalyze, err := twitterClient.GetTweetsByConversationID(mention.ConversationID)
		if err != nil {
			log.Println("Error when getting the tweets to analyze from twitter in conversaionID -> " + mention.ConversationID)
			continue
		}

		// Analyze the tweets using the neural network api
		results, err := GetSentimentFromTweets(tweetsToAnalyze)
		if err != nil {
			log.Println("Error when passing the tweets to the neural network: " + err.Error())
			continue
		}

		// There are different replies depending on the amount of tweets that can be analyzed

		var negativeTweets int
		var positiveTweets int

		for _, result := range results {
			if result > 50 {
				positiveTweets++
			} else {
				negativeTweets++
			}
		}

		welcomeMessages := map[string][]string{
			englishLang: {
				"Hi there! 😊",
				"So nice to see you! 😉",
				"Hello! 💁",
				"Hi! 💜 ",
				"Greetings! 🧐",
			},
			spanishLang: {
				"¡Holaaa! 😊",
				"¡Me alegro de verte! 😉",
				"¡Hola! 💁",
				"¡Aquí estoy! 💜 ",
				"¡Listo para analizar! 🧐",
			},
		}

		byeMessages := map[string][]string{
			englishLang: {
				"Bye! 👋",
				"See you soon! 🙃",
				"Bye bye! 😺",
				"¡Adiós! 🤠",
			},
			spanishLang: {
				"¡Hasta luego! 🤙",
				"¡Tá luego! 🤠",
				"¡Nos vemos! 🙃",
				"¡Adiós! 😺",
			},
		}

		negativeReaction := []string{
			"🙀", "😰", "😢", "😿", "😮", "🥴", "😱", "😪",
			"😥", "😨", "😭", "😢", "😲", "😧", "☹️", "🙁",
			"😦", "😵", "🥵", "🥶", "😳", "🤒", "😈",
			"☠️", "💀",
		}

		positiveReaction := []string{
			"🤙", "😄", "👍", "😁", "😺", "😃", "☀️", "😉",
			"😍", "🥰", "😊", "😎", "🤩", "🙂", "😇", "🥳",
			"😸",
		}

		// Generate random indexes for all the messages
		welcomeIndex := rand.Intn(len(welcomeMessages[englishLang]))
		byeIndex := rand.Intn(len(byeMessages[englishLang]))
		negativeIndex := rand.Intn(len(negativeReaction))
		positiveIndex := rand.Intn(len(positiveReaction))

		var responseText string
		l := len(results)

		// The response will be in the same language as the mention or in english as default

		responseNoTweets := map[string]string{

			englishLang: "There are no replies for me to analyse! " + negativeReaction[negativeIndex] + "\n" +
				"I can only see replies posted in the last 7 days!\n" +
				"Anyway, thank you for calling me " + negativeReaction[negativeIndex],

			spanishLang: "¡No he podido analizar ninguna respuesta! " + negativeReaction[negativeIndex] + "\n" +
				"Solo puedo ver respuestas publicadas en los últimos 7 días.\n" +
				"Muchas gracias por llamarme de todas formas " + positiveReaction[positiveIndex],
		}

		responseFewTweetsPositive := map[string]string{
			englishLang: "There are %v positive replies out of %v.\n",
			spanishLang: "Hay %v tweets positivos de un total de %v.\n",
		}

		responseFewTweetsNegative := map[string]string{
			englishLang: "There are %v negative replies out of %v.\n",
			spanishLang: "Hay %v tweets negativos de un total de %v.\n",
		}

		responseFewTweetsExtra := map[string]string{
			englishLang: "I could only analyse %v replies.\n" +
				"Notice that I can only see replies posted in the last 7 days!",
			spanishLang: "Solo pude analizar %v respuestas.\n" +
				"Nota: solo puedo analizar respuestas publicadas en los últimos 7 días.",
		}

		responseGeneralBalanced := map[string]string{
			englishLang: "The replies are perfectly balanced, as all things should be! ✨\n",
			spanishLang: "Las respuestas están perfectamente equilibradas, como todo debería ser ✨\n",
		}

		responseGeneralSlightNegative := map[string]string{
			englishLang: "The replies are slightly negative! 😶 \n",
			spanishLang: "Las respuestas son ligeramente negativas! 😶 \n",
		}

		responseGeneralSlightPositive := map[string]string{
			englishLang: "The replies are slightly positive! 😶 \n",
			spanishLang: "¡Las respuestas son ligeramente positivas! 😶 \n",
		}

		responseGeneralVeryNegative := map[string]string{
			englishLang: "%v%% of the replies are negative! %v \n",
			spanishLang: "¡El %v%% de las respuestas son negativas! %v \n",
		}

		responseGeneralVeryPositive := map[string]string{
			englishLang: "%v%% of the replies are positive! %v \n",
			spanishLang: "¡El %v%% de las respuestas son positivas! %v \n",
		}

		switch {

		// When there is not a single tweet that can be analyzed
		case l == 0:
			responseText = responseNoTweets[lang]

		// When there are only a few tweets (less than 10)
		case l < 10:
			// Most of them are negative
			if negativeTweets >= positiveTweets {
				responseText = fmt.Sprintf(responseFewTweetsNegative[lang], negativeTweets, l)

			} else {
				// Most of them are positive
				responseText = fmt.Sprintf(responseFewTweetsPositive[lang], positiveTweets, l)
			}
			responseText += fmt.Sprintf(responseFewTweetsExtra[lang], l)

		default:
			// General response: enough tweets to analyze
			// Welcome text
			responseText += welcomeMessages[lang][welcomeIndex] + "\n"
			negativePercentage := math.Round(float64(negativeTweets) * 100 / float64(l))
			positivePercentage := 100 - negativePercentage

			// If the percentages are very close
			if negativePercentage == 50 || positivePercentage == 50 {
				responseText += fmt.Sprintf(responseGeneralBalanced[lang])

				// Tweets are slightly negative
			} else if negativePercentage > 50 && negativePercentage <= 55 {
				responseText += fmt.Sprintf(responseGeneralSlightNegative[lang])

				// Tweets are slightly positive
			} else if positivePercentage > 50 && positivePercentage <= 55 {
				responseText += fmt.Sprintf(responseGeneralSlightPositive[lang])

				// Most tweets negative
			} else if negativeTweets >= positiveTweets {
				responseText += fmt.Sprintf(responseGeneralVeryNegative[lang], negativePercentage, negativeReaction[negativeIndex])

			} else if negativeTweets < positiveTweets {
				// Most tweets positive
				responseText += fmt.Sprintf(responseGeneralVeryPositive[lang], positivePercentage, positiveReaction[positiveIndex])
			}
			// Add the farewell message
			responseText += byeMessages[lang][byeIndex]
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
			log.Println("Could not process mention with id " + mention.ID + ": " + err.Error())
			continue
		}

		// Store the Tweet ID of the mention in the database to avoid doing it twice
		err = twitterClient.SetMentionDone(mention)
		if err != nil {
			log.Println("Error when inserting done mentions " + err.Error())
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

	var sentiments []int

	for _, tweet := range tweets {

		// Delete the '/' that can cause trouble when making queries
		tweet.Text = strings.Replace(tweet.Text, "/", "", -1)

		url := "http://neural-network:8081/v1/sentiment/" + tweet.Language

		client := &http.Client{}
		req, _ := http.NewRequest("GET", url, nil)
		req.Header.Set("sentiment", base64.StdEncoding.EncodeToString([]byte(tweet.Text)))

		resp, err := client.Do(req)
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
