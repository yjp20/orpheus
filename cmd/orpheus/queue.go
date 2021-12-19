package main

import (
	"fmt"
	"sort"
	"math/rand"
)

func (server *Server) Add(url string, userId string, shuf bool) ([]*QueueItem, error) {
	songs, err := fetchSongFromURL(url)
	if err != nil {
		return nil, err
	}
	if shuf {
		for i := range songs {
			j := rand.Intn(i+1)
			songs[i], songs[j] = songs[j], songs[i]
		}
	}
	queuedList := make([]*QueueItem, 0)
	for _, song := range songs {
		queueItem := &QueueItem{
			Song:     song,
			QueuedBy: userId,
			Index:    max(server.getQueueSum(userId), server.currentIndex()) + int64(song.Length),
			Dynamic: true,
		}
		server.Queue = append(server.Queue, queueItem)
		queuedList = append(queuedList, queueItem)
		server.sortQueue()
		server.triggerUpdate()
	}
	return queuedList, err
}

func (server *Server) SkipTo(index int) (*QueueItem, error) {
	if index >= len(server.Queue) || index < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	server.Index = index
	server.triggerUpdate()
	return server.Queue[index], nil
}

func (server *Server) Move(from, to int) (*QueueItem, error) {
	if from >= len(server.Queue) || to >= len(server.Queue) || from < 0 || to < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	target := server.Queue[from]
	if from == to {
		return target, nil
	}
	if to+1 == len(server.Queue) {
		target.Index = server.Queue[to].Index + int64(target.Song.Length)
	} else if to == 0 {
		target.Index = server.Queue[to].Index - int64(target.Song.Length)
	} else if to > from {
		target.Index = (server.Queue[to].Index + server.Queue[to+1].Index) / 2
	} else if to < from {
		target.Index = (server.Queue[to].Index + server.Queue[to-1].Index) / 2
	}
	target.Dynamic = false
	server.sortQueue()
	server.triggerUpdate()
	return target, nil
}

func (server *Server) Remove(index int) (*QueueItem, error) {
	if index >= len(server.Queue) || index < 0 {
		return nil, fmt.Errorf("index out of bounds")
	}
	queueItem := server.Queue[index]
	server.Queue = append(server.Queue[:index], server.Queue[index+1:]...)
	if index < server.Index {
		server.Index -= 1
	}
	server.triggerUpdate()
	return queueItem, nil
}

func (server *Server) Shuffle() {
	if len(server.Queue) == 0 {
		return
	}
	for _, item := range server.Queue {
		item.Index = 0
		item.Dynamic = false
	}
	for i := range server.Queue {
		j := rand.Intn(i+1)
		server.Queue[i], server.Queue[j] = server.Queue[j], server.Queue[i]
	}
}

func (server *Server) triggerUpdate() {
	if server.Index < len(server.Queue) {
		queueItem := server.Queue[server.Index]
		if queueItem.Song != server.Player.Song {
			server.Player.PlaySong(queueItem.Song)
		}
	}
}

func (server *Server) nextSong(killed bool) {
	if !killed {
		switch server.NextPolicy {
		case LoopQueue:
			server.Index = (server.Index + 1) % len(server.Queue)
			server.triggerUpdate()
		case LoopSong:
			server.Player.Seek(0)
		}
	}
}

func (server *Server) getQueueSum(userId string) int64 {
	sum := int64(0)
	for _, item := range server.Queue {
		if item.Dynamic && item.QueuedBy == userId {
			sum += int64(item.Song.Length)
		}
	}
	return sum
}

func (server *Server) sortQueue() {
	target := server.Queue[server.Index]
	sort.Slice(server.Queue, func(i, j int) bool {
		return server.Queue[i].Index < server.Queue[j].Index
	})
	for index, queueItem := range server.Queue {
		if target == queueItem {
			server.Index = index
		}
	}
}

func (server *Server) currentIndex() int64 {
	if len(server.Queue) == 0 {
		return 0
	}
	return server.Queue[server.Index].Index
}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
