package main

import (
	"os"
	"os/exec"
	"strings"
	"time"

	mp3 "github.com/hajimehoshi/go-mp3"
)

func fetchMusicFromURL(url string) (*Song, error) {
	metaDataProcess := exec.Command("yt-dlp", "--skip-download", "--get-title", "--get-id", "--no-warnings", url)
	metaData, err := metaDataProcess.Output()
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(string(metaData), "\n")
	title, id := tokens[0], tokens[1]
	path := "./data/" + id + ".mp3"

	downloadProcess := exec.Command("yt-dlp", "--id", "-x", "--audio-format", "mp3", "-P", "./data", url)
	err = downloadProcess.Run()
	if err != nil {
		return nil, err
	}

	sampleRate, length, err := getMP3MetaData(path)
	if err != nil {
		return nil, err
	}

	if sampleRate != 48000 {
		tempPath := "/tmp/" + id + ".mp3"
		resampleProcess := exec.Command("ffmpeg", "-i", path, "-ar", "48000", tempPath)
		out, err := resampleProcess.CombinedOutput()
		if err != nil {
			println(string(out))
			return nil, err
		}
		err = os.Rename(tempPath, path)
		if err != nil {
			return nil, err
		}
	}

	return &Song{title, url, length, path}, nil
}

func getMP3MetaData(path string) (sampleRate int, length time.Duration, err error) {
	file, err := os.Open(path)
	if err != nil {
		return
	}
	d, err := mp3.NewDecoder(file)
	if err != nil {
		return
	}
	sampleRate = d.SampleRate()
	length = time.Duration(d.Length() * int64(time.Second) / int64(sampleRate) / 4)
	return
}
