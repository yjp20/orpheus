package main

import (
	"flag"
)

var (
	token = flag.String("token", "", "token of discord bot")
)

func main() {
	session := Login(*token)
	session.AddHandler(messageHandler)
}
