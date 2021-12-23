package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	mp3 "github.com/hajimehoshi/go-mp3"
	"layeh.com/gopus"
)

const (
	BYTES_IN_FRAME = int64(4)
)

type Player struct {
	Playing  bool
	Time     time.Duration
	Song     *Song
	Voice    *discordgo.VoiceConnection
	Callback func(killed bool)

	events chan playerEvent
	mu     sync.Mutex
}

type playerEvent int

const (
	kill playerEvent = iota
	pause
	resume
	seek
)

func (p *Player) PlaySong(song *Song) error {
	p.killWorker()

	p.mu.Lock()
	p.Playing = true
	p.Song = song
	p.Time = 0
	p.mu.Unlock()

	err := p.startWorker()
	if err != nil {
		log.Println(err)
	}
	return err
}

func (p *Player) Resume() {
	p.mu.Lock()
	p.Playing = true
	if p.events != nil {
		p.events <- resume
	}
	p.mu.Unlock()
}

func (p *Player) Pause() {
	p.mu.Lock()
	p.Playing = false
	if p.events != nil {
		p.events <- pause
	}
	p.mu.Unlock()
}

func (p *Player) FastForward(seconds float64) {
	p.mu.Lock()
	p.Time = p.clampTime(p.Time + time.Duration(float64(time.Second)*seconds))
	if p.events != nil {
		p.events <- seek
	}
	p.mu.Unlock()
}

func (p *Player) Seek(seconds float64) {
	p.mu.Lock()
	if p.events != nil {
		p.Time = p.clampTime(time.Duration(float64(time.Second) * seconds))
		p.events <- seek
	}
	p.mu.Unlock()
}

func (p *Player) seek(decoder *mp3.Decoder) error {
	offset := BYTES_IN_FRAME * (int64(p.Time) / (int64(time.Second) / int64(decoder.SampleRate())))
	_, err := decoder.Seek(offset, io.SeekStart)
	return err
}

func (p *Player) killWorker() {
	p.mu.Lock()
	if p.events != nil {
		events := p.events
		p.events = nil
		events <- kill
	}
	p.mu.Unlock()
}

func (p *Player) startWorker() error {
	if p.Song == nil {
		return fmt.Errorf("player must have song initialized")
	}

	for !p.Song.IsDownloaded {
		<-p.Song.download
	}

	file, err := os.Open(p.Song.File)
	if err != nil {
		return err
	}
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}
	err = p.seek(decoder)
	if err != nil {
		return err
	}

	go p.audioWorker(decoder)
	return nil
}

func (p *Player) audioWorker(decoder *mp3.Decoder) {
	killed := false
	events := make(chan playerEvent)
	p.mu.Lock()
	p.events = events
	p.mu.Unlock()
	sampleRate := decoder.SampleRate()
	frameSize := sampleRate / 50
	encoder, _ := gopus.NewEncoder(sampleRate, 2, gopus.Audio)
	buffer16 := make([]int16, frameSize*2)

	for p.Voice == nil {
		<-events
	}

	if p.Playing {
		goto playing
	} else {
		goto paused
	}

playing:
	p.Voice.Speaking(true)
	for {
		select {
		case e := <-events:
			switch e {
			case kill:
				killed = true
				goto cleanup
			case resume:
				continue
			case pause:
				goto paused
			case seek:
				p.seek(decoder)
			}

		default:
			err := binary.Read(decoder, binary.LittleEndian, &buffer16)
			if err != nil {
				goto cleanup
			}
			res, err := encoder.Encode(buffer16, frameSize, frameSize*4)
			if err != nil {
				goto cleanup
			}
			p.Time += time.Duration(int64(time.Second) / int64(decoder.SampleRate()) * int64(len(buffer16)) / BYTES_IN_FRAME)
			p.Voice.OpusSend <- res
		}
	}

paused:
	p.Voice.Speaking(false)
	for {
		select {
		case e := <-events:
			switch e {
			case kill:
				killed = true
				goto cleanup
			case resume:
				goto playing
			case pause:
				continue
			case seek:
				p.seek(decoder)
			}
		}
	}

cleanup:
	p.Voice.Speaking(false)
	p.mu.Lock()
	p.events = nil
	p.mu.Unlock()
	close(events)
	if p.Callback != nil {
		go p.Callback(killed)
	}
}

func (p *Player) clampTime(t time.Duration) time.Duration {
	if t < 0 {
		return time.Duration(0)
	}
	if t > p.Song.Length {
		return p.Song.Length
	}
	return t
}
