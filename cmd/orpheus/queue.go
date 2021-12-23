package main

import (
	"fmt"
	"math/rand"
	"sort"
	"time"
)

type addPolicy int

const (
	Next addPolicy = iota
	Now
	Last
	Smart
)

func (server *Server) Add(songs []*Song, userId string, shuffle bool, policy addPolicy) []*QueueItem {
	if shuffle {
		rand.Seed(time.Now().UnixNano())
		rand.Shuffle(len(songs), func(i, j int) { songs[i], songs[j] = songs[j], songs[i] })
	}
	queueItems := make([]*QueueItem, 0)
	for _, song := range songs {
		queueItem := &QueueItem{
			Song:     song,
			QueuedBy: userId,
			Dynamic:  true,
		}
		switch policy {
		case Next, Now:
			if len(server.Queue) == 0 || server.Index == len(server.Queue)-1 {
				queueItem.Rank = server.currentRank() + int64(song.Length)
			} else {
				queueItem.Rank = (server.Queue[server.Index].Rank + server.Queue[server.Index+1].Rank) / 2
			}
		case Last:
			if len(server.Queue) == 0 {
				queueItem.Rank = int64(song.Length)
			} else {
				queueItem.Rank = server.Queue[len(server.Queue)-1].Rank + int64(song.Length)
			}
		case Smart:
			queueItem.Rank = max(server.getQueueSum(userId), server.currentRank()) + int64(song.Length)
		}
		server.Queue = append(server.Queue, queueItem)
		queueItems = append(queueItems, queueItem)
		server.sortQueue()
		server.triggerUpdate()
	}
	if policy == Now {
		server.nextSong(false)
	}
	return queueItems
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
		target.Rank = server.Queue[to].Rank + int64(target.Song.Length)
	} else if to == 0 {
		target.Rank = server.Queue[to].Rank - int64(target.Song.Length)
	} else if to > from {
		target.Rank = (server.Queue[to].Rank + server.Queue[to+1].Rank) / 2
	} else if to < from {
		target.Rank = (server.Queue[to].Rank + server.Queue[to-1].Rank) / 2
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
		item.Rank = 0
		item.Dynamic = false
	}
	rand.Seed(time.Now().UnixNano())
	rand.Shuffle(len(server.Queue)-server.Index-1, func(i, j int) {
		server.Queue[i+server.Index+1], server.Queue[j+server.Index+1] =
			server.Queue[j+server.Index+1], server.Queue[i+server.Index+1]
	})
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
			server.Player.PlaySong(server.Queue[server.Index].Song)
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
		return server.Queue[i].Rank < server.Queue[j].Rank
	})
	for index, queueItem := range server.Queue {
		if target == queueItem {
			server.Index = index
		}
	}
}

func (server *Server) currentRank() int64 {
	if len(server.Queue) == 0 {
		return 0
	}
	return server.Queue[server.Index].Rank
}

func max(x, y int64) int64 {
	if x < y {
		return y
	}
	return x
}
