package main

import (
	"log"
	"sort"
	"time"
)

type Song struct {
	Name   string
	URL    string
	Length time.Duration
	File   string
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

var servers map[string](*Server)

func sortServerQueue(server *Server) {
	sort.Slice(server.Queue, func(i, j int) bool {
		return server.Queue[i].Index < server.Queue[j].Index
	})
}

func addServer(id string) {
	_, ok := servers[id]
	if !ok {
		servers[id] = &Server{id, make([](*QueueItem), 0, 5), 0, make(map[string](*User))}
	}
}

func getServers(access []string) []string {
	ans := make([]string, 0, 10)
	for _, v := range access {
		if _, ok := servers[v]; ok {
			ans = append(ans, v)
		}
	}
	return ans
}

func add(serverId string, url string, userId string) Song {
	server, ok := servers[serverId]
	if !ok {
		log.Fatal()
	}
	s := fetchMusicFromURL(url)
	user, ok := server.Users[userId]
	if !ok {
		user = &User{make([](*QueueItem), 0, 10), userId, 0.0}
		server.Users[userId] = user
	}
	item := QueueItem{s, userId, user.LengthSum}
	user.LengthSum += s.Length.Seconds()
	server.Queue = append(server.Queue, &item)
	user.Songs = append(user.Songs, &item)
	sortServerQueue(server)
	return s
}
