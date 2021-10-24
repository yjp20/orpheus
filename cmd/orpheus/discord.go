package main

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/dgvoice"
	"github.com/bwmarrin/discordgo"
	"github.com/faiface/beep"
	"github.com/faiface/beep/mp3"
)

func Login(token string) *discordgo.Session {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal(err)
	}
	return session
}

type PlayInstance struct {
	Playing    bool
	Time       time.Duration
	Target     *discordgo.VoiceConnection
	Song       *Song
	SongFormat beep.Format
	SongSource beep.StreamSeekCloser
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
	go instance.startStream()
	return &instance, nil
}

func (instance *PlayInstance) startStream() {
	instance.Target.Speaking(true)
	pcm := make(chan []int16)
	dgvoice.SendPCM(instance.Target, pcm)
	for {
		data := make([][2]float64, 256)
		_, ok := instance.SongSource.Stream(data)
		if !ok {
			break
		}
		pcmData := make([]int16, 256)
		for i := 0; i < 256; i += 1 {
			pcmData[i] = int16(data[i][0])
		}
		pcm <- pcmData
	}
	close(pcm)
	instance.Target.Speaking(false)
	instance.Target.Disconnect()
}

func (instance *PlayInstance) changeSong(song *Song) error {
	file, err := os.Open(song.File)
	if err != nil {
		return err
	}
	defer file.Close()

	streamer, format, err := mp3.Decode(file)
	if err != nil {
		return err
	}
	instance.Song = song
	instance.SongFormat = format
	instance.SongSource = streamer
	return nil
}
