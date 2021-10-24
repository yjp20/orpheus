package main

import (
	"os"
)

type Song struct {
	Name   string
	URL    string
	Length float64
	File   *os.File
}

type QueueItem struct {
	Song     Song
	QueuedBy string
	Index    float64
}

type User struct {
	Songs     [](*QueueItem)
	Id        string
	LengthSum float64
}

type Server struct {
	Id    string
	Queue [](*QueueItem)
	Index int
	Users map[string](*User)
}
