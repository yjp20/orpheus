package main

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	mp3 "github.com/hajimehoshi/go-mp3"
)

func fetchSongFromURL(url string) (song *Song, err error) {
	metaDataProcess := exec.Command("yt-dlp", "--print", "%(title)s\n%(id)s\n%(duration)d", "--no-warnings", url)
	metaData, err := metaDataProcess.Output()
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(string(metaData), "\n")
	song = &Song{
		Name:         tokens[0],
		ID:           tokens[1],
		File:         "./data/" + tokens[1] + ".mp3",
		IsDownloaded: true,
		download:     make(chan int),
	}

	if _, err := os.Stat(song.File); errors.Is(err, os.ErrNotExist) {
		song.IsDownloaded = false
		go downloadSong(url, song)
	}

	if song.IsDownloaded {
		_, song.Length, err = getMP3MetaData(song.File)
		if err != nil {
			return nil, err
		}
	} else {
		seconds, err := strconv.Atoi(tokens[2])
		song.Length = time.Second * time.Duration(seconds)
		if err != nil {
			return nil, err
		}
	}

	return song, nil
}

func downloadSong(url string, song *Song) error {
	downloadProcess := exec.Command("yt-dlp", "--id", "-x", "--audio-format", "mp3", "-P", "./data", url)
	err := downloadProcess.Run()
	if err != nil {
		return fmt.Errorf("Failed to download '%s'\nerror: %s\n", url, err.Error())
	}
	var sampleRate int
	sampleRate, song.Length, err = getMP3MetaData(song.File)
	if sampleRate != 48000 {
		tempPath := "/tmp/" + song.ID + ".mp3"
		resampleProcess := exec.Command("ffmpeg", "-i", song.File, "-ar", "48000", tempPath)
		err = resampleProcess.Run()
		if err != nil {
			return err
		}
		err = os.Rename(tempPath, song.File)
		if err != nil {
			return err
		}
	}

	song.IsDownloaded = true
	song.download <- 0
	close(song.download)
	return nil
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
