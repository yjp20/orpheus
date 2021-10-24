package main

import (
    "os/exec"
    "log"
    "io"
    "github.com/tcolgate/mp3"
    "os"
    "strings"
)

//cache := map[string]Song

func fetchMusicFromURL(url string) Song{
    cmd := exec.Command("yt-dlp", "--id", "-x", "--audio-format", "mp3", "-P", "./data", url)
    err := cmd.Run()
    if err != nil {
        log.Fatal(err)
    }

    staff, err := exec.Command("yt-dlp", "--skip-download", "--get-title", "--get-id", "--no-warnings", url).Output()
    if err != nil {
        log.Fatal(err)
    }
    stuff := strings.Split(string(staff), "\n")
    title := string(stuff[0])
    fil := "./data/"+string(stuff[1])+".mp3"

    t := 0.0
    path, err := os.Open(fil)
    if err != nil {
        log.Fatal(err)
    }
    d := mp3.NewDecoder(path)
    var f mp3.Frame
    skipped := 0

    for {
        if err := d.Decode(&f, &skipped); err != nil {
            if err == io.EOF {
                break
            }
        }
        t = t + f.Duration().Seconds()
    }
    path, _ = os.Open(fil)

    s := Song{title, url, t, path}
    return s
}
