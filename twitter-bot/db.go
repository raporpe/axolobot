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
		log.Fatal("Cannot open database connection. Retrying...")
		session, err = mysql.Open(settings)
	}
	log.Println("Connected to the database")
	return &DatabaseManager{sess: session}

}

// Checks for existence of Tweet ID in the table Mention of the database
func (db *DatabaseManager) IsMentionDone(tweet Tweet) bool {

	result, err := db.sess.Collection("mention").Find("mention_id", tweet.ID).Count()
	if err != nil {
		log.Fatal("There was an error checking if mention was in database: " + err.Error())
		return false
	}

	return result > 0
}

// Inserts the Tweet ID in the table Mention of the database
func (db *DatabaseManager) InsertMention(tweet Tweet) error {
	toInsert := map[string]string{
		"mention_id": tweet.ID,
	}
	_, err := db.sess.Collection("mention").Insert(toInsert)
	if err != nil {
		log.Fatal("There was an error inserting mention in database: " + err.Error())
	}

	return nil
}

// Gets the most recent TweetID stored in Mention table in the database
func (db *DatabaseManager) GetLastMentionID() string {
	var lastMention map[string]string
	res := db.sess.Collection("mention").Find().OrderBy("-mention_id").Limit(1)
	err := res.One(&lastMention)
	if err != nil {
		log.Fatal("There was an error getting mentions in database: " + err.Error())
		return "1"
	}
	fmt.Println(lastMention)
	return (lastMention["mention_id"])
}
