package main

import (
	"bufio"
	"encoding/binary"
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
	mp3 "github.com/hajimehoshi/go-mp3"
	"layeh.com/gopus"
)

func Login(token string) *discordgo.Session {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

type PlayInstance struct {
	Playing bool
	Time    time.Duration
	Target  *discordgo.VoiceConnection
	Song    *Song
	kill    chan int
}

func NewPlayer(bot *discordgo.Session, guildId string, channelId string) (*PlayInstance, error) {
	voice, err := bot.ChannelVoiceJoin(guildId, channelId, false, false)
	if err != nil {
		return nil, err
	}
	instance := PlayInstance{
		Playing: false,
		Time:    0,
		Target:  voice,
	}
	return &instance, nil
}

func (instance *PlayInstance) changeSong(song *Song) error {
	if instance.kill != nil {
		instance.kill <- 0
	}
	file, err := os.Open(song.File)
	if err != nil {
		return err
	}
	decoder, err := mp3.NewDecoder(file)
	if err != nil {
		return err
	}
	instance.kill = make(chan int)
	instance.Song = song
	go instance.streamAudio(decoder, instance.kill)
	return nil
}

func (instance *PlayInstance) streamAudio(decoder *mp3.Decoder, kill chan int) {
	instance.Target.Speaking(true)
	encoder, _ := gopus.NewEncoder(48000, 2, gopus.Audio)
	buffered := bufio.NewReader(decoder)
	buffer16 := make([]int16, 960*2)
	for {
		select {
		case <-kill:
			goto done
		default:
			binary.Read(buffered, binary.LittleEndian, &buffer16)
			res, err := encoder.Encode(buffer16, 960, 960*4)
			if err != nil {
				goto done
			}
			instance.Target.OpusSend <- res
		}
	}
done:
	instance.Target.Speaking(false)
}
