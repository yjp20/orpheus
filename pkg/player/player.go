package player

import (
	"fmt"
	"io"
	"log"
	"os"
	"sync"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/go-audio/audio"
	"github.com/go-audio/wav"
	"layeh.com/gopus"

	"github.com/yjp20/orpheus/pkg/music"
)

const (
	CHANNELS    = 2
	SAMPLE_RATE = 48000
	FRAME_SIZE  = SAMPLE_RATE / 50
)

type Player struct {
	Playing       bool
	Time          time.Duration
	Song          *music.Song
	Voice         *discordgo.VoiceConnection
	FinishHandler func()

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

func (p *Player) Stop() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Playing = false
	p.Time = 0
	p.killWorker()
}

func (p *Player) PlaySong(song *music.Song) {
	p.mu.Lock()
	defer p.mu.Unlock()

	p.killWorker()
	p.Song = song
	p.Time = 0
	p.Playing = true

	if p.Song == nil {
		p.Playing = false
		return
	}

	err := p.startWorker()
	if err != nil {
		log.Println(err)
	}
}

func (p *Player) Resume() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Playing = true
	if p.events != nil {
		p.events <- resume
	}
}

func (p *Player) Pause() {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Playing = false
	if p.events != nil {
		p.events <- pause
	}
}

func (p *Player) FastForward(seconds float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	p.Time = p.clampTime(p.Time + time.Duration(float64(time.Second)*seconds))
	if p.events != nil {
		p.events <- seek
	}
}

func (p *Player) Seek(seconds float64) {
	p.mu.Lock()
	defer p.mu.Unlock()
	if p.events != nil {
		p.Time = p.clampTime(time.Duration(float64(time.Second) * seconds))
		p.events <- seek
	}
}

func (p *Player) seek(decoder io.Seeker) error {
	offset := 2 * CHANNELS * int64(p.Time) * SAMPLE_RATE / int64(time.Second)
	_, err := decoder.Seek(offset, io.SeekStart)
	return err
}

func (p *Player) killWorker() {
	if p.events != nil {
		p.events <- kill
	}
}

func (p *Player) startWorker() error {
	if p.Song == nil {
		return fmt.Errorf("player must have song initialized")
	}

	file, err := os.Open(p.Song.File)
	if err != nil {
		return err
	}
	decoder := wav.NewDecoder(file)
	err = p.seek(decoder)
	if err != nil {
		return err
	}

	p.events = make(chan playerEvent)
	go p.audioWorker(decoder)
	return nil
}

func (p *Player) audioWorker(decoder *wav.Decoder) {
	killed := false

	encoder, _ := gopus.NewEncoder(SAMPLE_RATE, CHANNELS, gopus.Audio)
	buffer16 := make([]int16, FRAME_SIZE*CHANNELS)
	buffer := audio.IntBuffer{Data: make([]int, FRAME_SIZE*CHANNELS)}

	for p.Voice == nil {
		select {
		case <-p.events:
			continue
		default:
			continue
		}
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
		case e := <-p.events:
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
			_, err := decoder.PCMBuffer(&buffer)
			if err != nil {
				goto cleanup
			}
			for i, v := range buffer.Data {
				buffer16[i] = int16(v)
			}
			res, err := encoder.Encode(buffer16, FRAME_SIZE, FRAME_SIZE*4)
			if err != nil {
				goto cleanup
			}
			p.Time += time.Second * FRAME_SIZE / SAMPLE_RATE
			p.Voice.OpusSend <- res
			if p.Time >= p.Song.Length {
				goto cleanup
			}
		}
	}

paused:
	p.Voice.Speaking(false)
	for {
		switch <-p.events {
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

cleanup:
	p.Voice.Speaking(false)
	close(p.events)
	p.events = nil
	if !killed && p.FinishHandler != nil {
		go p.FinishHandler()
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
