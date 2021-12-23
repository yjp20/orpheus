package main

import (
	"os/exec"
	"strings"
	"strconv"
	"time"
)

func fetchSongsFromURL(url string, playlist bool) (songs []*Song, err error) {
	listFlag := "--no-playlist"
	if playlist {
		listFlag = "--yes-playlist"
	}
	metaDataProcess := exec.Command("yt-dlp", listFlag, "--id", "-q", "-f", "ba", "-x", "--audio-format", "wav", "--no-simulate", "--print", "%(title)s\n%(id)s\n%(duration)d\n%(acodec)s", "-P", "./data", "--no-warnings", url)
	metaData, err := metaDataProcess.Output()
	if err != nil {
		return nil, err
	}

	tokens := strings.Split(string(metaData), "\n")
	numSongs := len(tokens) / 4
	songs = make([]*Song, numSongs)
	for i := 0; i < numSongs; i++ {
		title, id, duration, _ := tokens[4*i], tokens[4*i+1], tokens[4*i+2], tokens[4*i+3]
		seconds, _:= strconv.Atoi(duration)
		songs[i] = &Song{
			Name:   title,
			ID:     id,
			File:   "./data/" + id  + ".wav",
			Length: time.Duration(seconds) * time.Second,
			Format: Wav,
		}
	}

	return songs, nil
}
