package main

import (
	"flag"
)

var (
	token = flag.String("token", "", "token of discord bot")
)

func main() {
	flag.Parse()

	session := Login(*token)
	session.Open()
	session.AddHandler(commandHandler)
}
