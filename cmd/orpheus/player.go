package main

import (
	"encoding/binary"
	"fmt"
	"io"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	mp3 "github.com/hajimehoshi/go-mp3"
	"layeh.com/gopus"
)

var (
	BYTES_IN_FRAME = int64(4)
)

type Player struct {
	Playing  bool
	Time     time.Duration
	Song     *Song
	Voice    *discordgo.VoiceConnection
	Callback *func(killed bool)
	process  *playerProcess
}

type playerProcess struct {
	kill   chan int
	exit   chan int
	pause  chan int
	resume chan int
	seek   chan int
}

func (p *Player) PlaySong(song *Song) error {
	p.killWorker()
	p.Playing = true
	p.Song = song
	p.Time = 0
	return p.startWorker()
}

func (p *Player) Resume() {
	if p.process != nil {
		p.Playing = true
		p.process.resume <- 0
	}
}

func (p *Player) Pause() {
	if p.process != nil {
		p.Playing = false
		p.process.pause <- 0
	}
}

func (p *Player) FastForward(seconds float64) {
	p.Time = p.clampTime(p.Time + time.Duration(float64(time.Second)*seconds))
	if p.process != nil {
		p.process.seek <- 0
	}
}

func (p *Player) Seek(seconds float64) {
	p.Time = p.clampTime(time.Duration(float64(time.Second) * seconds))
	if p.process != nil {
		p.process.seek <- 0
	}
}

func (p *Player) seek(decoder *mp3.Decoder) error {
	offset := BYTES_IN_FRAME * (int64(p.Time) / (int64(time.Second) / int64(decoder.SampleRate())))
	_, err := decoder.Seek(offset, io.SeekStart)
	return err
}

func (p *Player) killWorker() {
	if p.process != nil {
		process := p.process
		p.process = nil
		process.kill <- 0
		<-process.exit
	}
}

func (p *Player) startWorker() error {
	if p.Song == nil {
		return fmt.Errorf("player must have song initialized")
	}
	if p.Voice == nil {
		return fmt.Errorf("player must have voice channel initialized")
	}

	file, err := os.Open(p.Song.File)
	if err != nil {
		return err
	}
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}
	p.seek(decoder)
	go p.audioWorker(decoder)
	return nil
}

func (p *Player) audioWorker(decoder *mp3.Decoder) {
	process := &playerProcess{
		kill:   make(chan int),
		exit:   make(chan int),
		pause:  make(chan int),
		resume: make(chan int),
		seek:   make(chan int),
	}
	p.process = process
	killed := false
	sampleRate := decoder.SampleRate()
	frameSize := sampleRate / 50
	encoder, _ := gopus.NewEncoder(sampleRate, 2, gopus.Audio)
	buffer16 := make([]int16, frameSize*2)

	if p.Playing {
		goto playing
	} else {
		goto paused
	}

playing:
	p.Voice.Speaking(true)
	for {
		select {
		case <-process.kill:
			killed = true
			goto cleanup
		case <-process.resume:
			continue
		case <-process.pause:
			goto paused
		case <-process.seek:
			p.seek(decoder)

		default:
			binary.Read(decoder, binary.LittleEndian, &buffer16)
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
		case <-process.kill:
			killed = true
			goto cleanup
		case <-process.resume:
			goto playing
		case <-process.seek:
			p.seek(decoder)
		case <-process.pause:
			continue
		}
	}

cleanup:
	p.Voice.Speaking(false)
	if p.Callback != nil {
		(*p.Callback)(killed)
	}
	process.exit <- 0
	close(process.kill)
	close(process.exit)
	close(process.pause)
	close(process.resume)
	close(process.seek)
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
