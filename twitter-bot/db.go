package main

import (
	"fmt"
	"log"
	"os"
	"time"

	"github.com/upper/db/v4"
	"github.com/upper/db/v4/adapter/mysql"
)

type DatabaseManager struct {
	sess db.Session
}

// Performs connection with the database
func NewDatabaseManager() *DatabaseManager {
	dbPassword := os.Getenv("DB_PASSWORD")
	settings := mysql.ConnectionURL{
		User:     "root",
		Password: dbPassword,
		Host:     "database",
		Database: "axolobot",
	}

	time.Sleep(5 * time.Second)
	session, err := mysql.Open(settings)
	for err != nil {
		log.Println("Cannot open database connection. Retrying...")
		session, err = mysql.Open(settings)
	}
	log.Println("Connected to the database")
	return &DatabaseManager{sess: session}

}

// Checks if the mention ID is marked as done in the database
func (db *DatabaseManager) IsMentionDone(tweet Tweet) bool {

	var mentionData map[string]string
	result := db.sess.Collection("mention").Find("mention_id", tweet.ID)
	err := result.One(&mentionData)
	if err != nil {
		log.Fatal("There was an error checking if mention was done in database: " + err.Error())
		return false
	}

	return mentionData["is_done"] == "1"
}

// Checks if the mention ID is registered in the database
func (db *DatabaseManager) IsMentionRegistered(tweet Tweet) bool {

	result, err := db.sess.Collection("mention").Find("mention_id", tweet.ID).Count()
	if err != nil {
		log.Fatal("There was an error checking if mention is registered in database: " + err.Error())
		return false
	}

	return result > 0
}

// Inserts the Tweet ID in the table Mention of the database
func (db *DatabaseManager) RegisterMention(tweet Tweet) error {
	toInsert := map[string]string{
		"mention_id": tweet.ID,
	}
	_, err := db.sess.Collection("mention").Insert(toInsert)
	if err != nil {
		log.Fatal("There was an error registering mention in database: " + err.Error())
	}

	return nil
}

// Inserts the Tweet ID in the table Mention of the database
func (db *DatabaseManager) setMentionDone(tweet Tweet) error {

	err := db.sess.Collection("mention").Find("mention_id", tweet.ID).Update(map[string]string{"is_done": "1"})

	if err != nil {
		log.Fatal("There was an error setting mention as done in database: " + err.Error())
	}

	return nil
}

// Gets the most recent TweetID stored in Mention table in the database
func (db *DatabaseManager) GetLastRegisteredMentionID() string {
	var lastMention map[string]string
	res := db.sess.Collection("mention").Find().OrderBy("-mention_id").Limit(1)
	err := res.One(&lastMention)
	if err != nil {
		log.Println("There was an error getting last registered mention in database: " + err.Error())
		return "1"
	}
	fmt.Println(lastMention)
	return (lastMention["mention_id"])
}
