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

func fetchSongFromURL(url string) (songs []*Song, err error) {
	metaDataProcess := exec.Command("yt-dlp", "--yes-playlist", "--print", "%(title)s\n%(id)s\n%(duration)d", "--no-warnings", url)
	metaData, err := metaDataProcess.Output()
	if err != nil {
		return nil, err
	}
	tokens := strings.Split(string(metaData), "\n")
	numsongs := len(tokens)/3
	songs = make([]*Song, 0)
	for i := 0; i < numsongs; i++ {
		song := &Song{
			Name:         tokens[i*3],
			ID:           tokens[i*3+1],
			File:         "./data/" + tokens[i*3+1] + ".mp3",
			IsDownloaded: true,
			download:     make(chan int),
		}

		if _, err := os.Stat(song.File); errors.Is(err, os.ErrNotExist) {
			song.IsDownloaded = false
			go downloadSong(song.ID, song)
		}

		if song.IsDownloaded {
			_, song.Length, err = getMP3MetaData(song.File)
			if err != nil {
				return nil, err
			}
		} else {
			seconds, err := strconv.Atoi(tokens[i*3+2])
			song.Length = time.Second * time.Duration(seconds)
			if err != nil {
				return nil, err
			}
		}
		songs = append(songs, song)
	}

	return songs, nil
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
