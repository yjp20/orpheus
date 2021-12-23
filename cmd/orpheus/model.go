package main

import (
	"time"
)

type Server struct {
	ID         string       `json:"id"`
	Queue      []*QueueItem `json:"queue"`
	Index      int          `json:"index"`
	Player     Player
	NextPolicy NextPolicy
}

type QueueItem struct {
	Song     *Song  `json:"song"`
	QueuedBy string `json:"queued_by"`
	Dynamic  bool   `json:"dynamic"`
	Rank     int64  `json:"rank"`
}

type Song struct {
	ID     string        `json:"id"`
	Name   string        `json:"name"`
	URL    string        `json:"url"`
	Length time.Duration `json:"length"`
	File   string        `json:"file"`
	Format Format        `json:"format"`
}

type NextPolicy int

const (
	LoopSong NextPolicy = iota
	LoopQueue
	NoLoop
)

type Format int

const (
	Opus Format = iota
	Mp3
	M4a
	Wav
	Aac
	Vorbis
	Flac
)

var servers = make(map[string]*Server)

func getServer(id string) *Server {
	_, ok := servers[id]
	if !ok {
		servers[id] = &Server{
			ID:         id,
			Queue:      make([]*QueueItem, 0),
			Index:      0,
			NextPolicy: LoopQueue,
		}
		servers[id].Player.Callback = servers[id].nextSong
	}
	return servers[id]
}

func unionServers(serverIDs []string) []string {
	ans := make([]string, 0)
	for _, id := range serverIDs {
		if _, ok := servers[id]; ok {
			ans = append(ans, id)
		}
	}
	return ans
}
