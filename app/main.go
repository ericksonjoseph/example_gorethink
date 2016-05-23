package main

import (
	"fmt"
	"log"

	r "github.com/dancannon/gorethink"
)

func main() {

	// This will hold our DB connection
	var db *r.Session

	// Connect to DB
	db, err := r.Connect(r.ConnectOpts{
		Address:  "192.168.99.100:28015",
		Database: "test",
	})
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Place connection in our context object
	context := AppContext{
		DB: db,
	}

	// Use routines & channels to listen to changefeed
	ChangeFeeds("transactions", context)

	// Keep master thread alive (press ENTER to end program)
	fmt.Scanln()
}

/**
 * Creates a channel and spawns two routines.
 * One to listen to the changefeed and push the changes to a channel
 * and one to accept data comming from the channel
 */
func ChangeFeeds(table string, context AppContext) {

	channel := make(chan interface{})

	// Open changefeed
	res, err := r.Table(table).Changes().Run(context.DB)
	if err != nil {
		log.Fatalln(err.Error())
	}

	// Listen to changefeed and push data into channel
	go func() {
		fmt.Println("listening to rethinkDB")
		var value interface{}
		for res.Next(&value) {
			// Wont accep values if there is no where to put them (buffer)
			channel <- value
		}
	}()

	// Accept data from channel
	go func() {
		fmt.Println("listening to channel")
		for {
			msg_from_channel := <-channel
			fmt.Println(msg_from_channel)
		}
	}()
}
