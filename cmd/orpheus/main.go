package main

import (
	"flag"
)

var (
	token = flag.String("token", "", "token of discord bot")
	addr  = flag.String("addr", ":4000", "address of bot api")
	cors  = flag.String("cors", "http://localhost:3000", "cors")
)

func main() {
	flag.Parse()

	session := Login(*token)
	session.Open()
	session.AddHandler(commandHandler)

	serveAPI(session, *addr, *cors)
}
