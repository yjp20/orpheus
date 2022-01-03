package main

import (
	"github.com/yjp20/orpheus/pkg/player"
	"github.com/yjp20/orpheus/pkg/queue"
)

type Guild struct {
	ID     string       `json:"id"`
	Queue  queue.Queue `json:"queue"`
	Player player.Player
}

var guilds = make(map[string]*Guild)

func GetGuild(id string) *Guild {
	_, ok := guilds[id]
	if !ok {
		guilds[id] = &Guild{ID: id, Queue: queue.NewQueue()}
		guilds[id].Player.FinishHandler = guilds[id].Queue.NextSong
		guilds[id].Queue.UpdateHandler = guilds[id].Player.PlaySong
	}
	return guilds[id]
}

func QueryGuilds(serverIDs []string) []string {
	ans := make([]string, 0)
	for _, id := range serverIDs {
		if _, ok := guilds[id]; ok {
			ans = append(ans, id)
		}
	}
	return ans
}
