package main

import (
    "os/exec"
    "strconv"
    "log"
    "strings"
)

//cache := map[string]Song

func fetchMusicFromURL(url string) Song{
    cmd := exec.Command("youtube-dl", "--id", "-x", url)
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }
    title, err := exec.Command("youtube-dl", "--skip-download", "--get-title", "--no-warnings", url).Output()
    if err != nil {
        log.Fatal(err)
    }
    duration, err := exec.Command("youtube-dl", "--skip-download", "--get-duration", "--no-warnings", url).Output()
    length := strings.Split(string(duration[:]), string(":"))
    minute, _ := strconv.Atoi(length[0])
    second, _ := strconv.Atoi(length[1])
    id, err := exec.Command("youtube-dl", "--skip-download", "--get-id", "--no-warnings", url).Output()
    f := string(id[:len(id)-1])+".m4a"
    s := Song{string(title[:]), url, minute*60+second, f}
    return s
}
