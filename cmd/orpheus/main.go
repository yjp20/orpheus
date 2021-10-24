package main

import (
	"flag"
    "time"
)

var (
	token = flag.String("token", "", "token of discord bot")
)

func main() {
	flag.Parse()

	session := Login(*token)
	session.Open()
	session.AddHandler(commandHandler)
    initCommands(session, "833278784848658462")
    time.Sleep(time.Minute)
}
