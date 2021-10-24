package main

import (
	"fmt"
	"log"
    "sort"
)

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
    user.LengthSum += s.Length
    server.Queue = append(server.Queue, &item)
    user.Songs = append(user.Songs, &item)
    sortServerQueue(server)
    return s
}

func main() {
    servers = make(map[string](*Server))
    s := fetchMusicFromURL("https://www.youtube.com/watch?v=_mMyPJSx8RU")
    fmt.Printf("%+v\n", s)
    addServer(string("macklin"))
    add("macklin", "https://www.youtube.com/watch?v=_mMyPJSx8RU", "mycho")
    add("macklin", "https://www.youtube.com/watch?v=dQw4w9WgXcQ", "mycho")
    add("macklin", "https://www.youtube.com/watch?v=Llr2dcd-VBo", "theory")
    add("macklin", "https://www.youtube.com/watch?v=IMXwy8KawGI", "theory")
    fmt.Printf("%+v %+v\n %+v %+v\n", *servers["macklin"].Queue[0], *servers["macklin"].Queue[1], *servers["macklin"].Queue[2], *servers["macklin"].Queue[3])
}
