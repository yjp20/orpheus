package main

import (
	"log"
	"sort"
	"time"
    "math"
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

func addServer(id string) bool {
	_, ok := servers[id]
	if !ok {
		servers[id] = &Server{id, make([](*QueueItem), 0, 5), 0, make(map[string](*User))}
		return true
	}
	return false
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

// TODO: different types of add (smart-algo, add-end, add-next)
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
	item := QueueItem{s, userId, math.Max(user.LengthSum, server.Queue[server.Index].Index)+1}
	user.LengthSum += s.Length.Seconds()
	server.Queue = append(server.Queue, &item)
	user.Songs = append(user.Songs, &item)
	sortServerQueue(server)
	return s
}

func skipTo(serverId string, index int) Song {
    server, ok := servers[serverId]
    if !ok {
        log.Fatal()
    }
    if index >= len(server.Queue) || index < 0 {
        log.Fatal()
    }
    server.Index = index
    s := server.Queue[index].Song
    return s
}
/*
func move(serverId string, from_index int, to_index int) Song {
    length = len(server.Queue)
    server, ok := servers[serverId]
    if !ok {
        log.Fatal()
    }
    if from_index >= length || to_index >= length || from_index < 0 || to_index < 0 {
        log.Fatal()
    }
    if to_index > from_index {
        to_index -= 1
    }
    if from_index == to_index {
        return server.Queue[from_index].Song
    }
    q := server.Queue[from_index]
    if to_index == 0 {
        q.Index = server.Queue[0].Index-1
    }
    else if to_index == length-1 {
        q.Index = server.Queue[length-1].Index+1
    }
    else {
        q.Index = (server.Queue[to_index].Index + server.Queue[to_index-1].Index)/2
    }
    temp := make([](*QueueItem), 0, 100)
    temp = append(temp, server.Queue[:from_index]...)
    temp = append(temp, server.Queue[from_index+1:]...)
    server.Queue = make([](*QueueItem), 0, 100)
    server.Queue = append(server.Queue, temp[:to_index]...)
    server.Queue = append(server.Queue, q)
    server.Queue = append(server.Queue, temp[to_index:]...)
    if 
}
*/
func remove(serverId string, index int) *QueueItem {
    server, ok := servers[serverId]
    if !ok {
        log.Fatal()
    }
    q := server.Queue[index]
    temp := make([](*QueueItem), 0, 100)
    temp = append(temp, server.Queue[:index]...)
    server.Queue = append(temp, server.Queue[index+1:]...)
    if index < server.Index {
        server.Index -= 1
    }
    user, ok := server.Users[q.QueuedBy]
    if !ok {
        log.Fatal()
    }
    user.LengthSum -= q.Song.Length.Seconds()
    for i, v := range user.Songs {
        if v == q {
            temp = make([](*QueueItem), 0, 100)
            temp = append(temp, server.Queue[:i]...)
            user.Songs = append(temp, server.Queue[i+1:]...)
        }
    }
    return q
}

