package main

import (
	"flag"
	"log"

	"github.com/bwmarrin/discordgo"
)

var (
	token = flag.String("token", "", "token of discord bot")
	appID = flag.String("appID", "", "appID of discord bot")
	addr  = flag.String("addr", ":4000", "address of bot api")
	cors  = flag.String("cors", "http://localhost:3000", "cors")
)

func main() {
	flag.Parse()

	session, err := discordgo.New("Bot " + *token)
	if err != nil {
		log.Fatal(err)
	}
	err = session.Open()
	if err != nil {
		log.Fatal(err)
	}

	// Discord command-based controller
	session.AddHandler(commandHandler)
	session.AddHandler(joinHandler)

	// Web-based controller
	server := serverAPI(session, *addr, *cors)
	server.ListenAndServe()
}
