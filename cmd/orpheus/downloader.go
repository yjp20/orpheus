package main

import (
	"io"
	"log"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/tcolgate/mp3"
)

func fetchMusicFromURL(url string) Song {
	downloadProcess := exec.Command("yt-dlp", "--id", "-x", "--audio-format", "mp3", "-P", "./data", url)
	err := downloadProcess.Run()
	if err != nil {
		log.Fatal(err)
	}

	titleProcess := exec.Command("yt-dlp", "--skip-download", "--get-title", "--get-id", "--no-warnings", url)
	titleString, err := titleProcess.Output()
	if err != nil {
		log.Fatal(err)
	}
	tokens := strings.Split(string(titleString), "\n")
	title := string(tokens[0])
	path := "./data/" + string(tokens[1]) + ".mp3"

	file, err := os.Open(path)
	if err != nil {
		log.Fatal(err)
	}
	length := getMP3Length(file)

	s := Song{title, url, length, path}
	return s
}

func getMP3Length(file *os.File) time.Duration {
	d := mp3.NewDecoder(file)
	var f mp3.Frame
	length := time.Duration(0)
	skipped := 0

	for {
		if err := d.Decode(&f, &skipped); err != nil {
			if err == io.EOF {
				break
			}
		}
		length = length + f.Duration()
	}

	return length
}
