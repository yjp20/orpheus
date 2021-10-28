package main

import (
	"fmt"
	"log"
	"strconv"
	"time"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{Name: "join", Description: "Join channel"},
	{Name: "queue", Description: "Show queue"},
	{Name: "pause", Description: "Pause playing"},
	{Name: "resume", Description: "Resume playing"},
	{Name: "add", Description: "Adds a song to the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "link to music",
				Required:    true,
			},
		},
	},
	{Name: "fastforward", Description: "Goes forward a given amount of seconds",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "seconds",
				Description: "number of seconds to skip forward",
				Required:    true,
			},
		},
	},
	{Name: "rewind", Description: "Goes back a given amount of seconds",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "seconds",
				Description: "number of seconds to skip backward",
				Required:    true,
			},
		},
	},
	{Name: "seek", Description: "Goes to a given amount of seconds",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "seconds",
				Description: "number of seconds to seek to",
				Required:    true,
			},
		},
	},
	{Name: "skip", Description: "Skips a song in the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "index",
				Description: "index to skip to",
				Required:    false,
			},
		},
	},
	{Name: "remove", Description: "Removes a song in the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "index",
				Description: "index of song to remove",
				Required:    false,
			},
		},
	},
}

func initCommands(s *discordgo.Session, guildId string) error {
	for _, command := range commands {
		_, err := s.ApplicationCommandCreate(s.State.User.ID, guildId, command)
		if err != nil {
			return err
		}
	}
	return nil
}

func commandHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	//perms, _ := s.State.MessagePermissions(m.Message)
	//perms>>discordgo.PermissionAdministrator&1 == 1

	var response string
	server, _ := servers[m.GuildID]

	switch m.ApplicationCommandData().Name {
	case "add":
		err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})

		queueItem, err := add(m.GuildID, m.ApplicationCommandData().Options[0].StringValue(), s.State.User.ID, s)
		if err != nil {
			log.Printf("failed to add song\nerror: %s\n", err)
			return
		}

		s.InteractionResponseEdit(*appID, m.Interaction, &discordgo.WebhookEdit{
			Content: formatSong(s, "Added ", server, queueItem),
		})

	case "queue":
		response = ""
		for i, v := range server.Queue {
			response += strconv.Itoa(i) + ".  " + v.Song.Name + "\n"
		}

	case "pause":
		server.Player.Pause()
		response = formatCurrentSong(s, "Pause", server)

	case "resume":
		server.Player.Resume()
		response = formatCurrentSong(s, "Resume", server)

	case "fastforward":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		server.Player.FastForward(seconds)
		response = formatCurrentSong(s, "Fast-forward", server)

	case "rewind":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		server.Player.FastForward(-seconds)
		response = formatCurrentSong(s, "Rewind", server)

	case "seek":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		server.Player.Seek(seconds)
		response = formatCurrentSong(s, "Seek", server)

	case "join":
		guild, err := s.State.Guild(m.GuildID)
		if err != nil {
			log.Printf("failed to fetch guild '%s'\nerror: %s\n", m.GuildID, err)
			return
		}
		channelID := ""
		for _, vs := range guild.VoiceStates {
			if vs.UserID == m.Member.User.ID {
				channelID = vs.ChannelID
			}
		}
		channel, err := s.State.Channel(channelID)
		if err != nil {
			log.Printf("failed to fetch channel '%s'\nerror: %s\n", channelID, err)
			return
		}
		server.Player.Voice, err = s.ChannelVoiceJoin(m.GuildID, channelID, false, false)
		if err != nil {
			log.Printf("failed to join channel '%s'\nerror: %s\n", channelID, err)
			return
		}
		response = fmt.Sprintf("Joining channel '%s'", channel.Name)
	}

	if response != "" {
		s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})
	}
}

func joinHandler(s *discordgo.Session, m *discordgo.GuildCreate) {
	if m.Guild.Unavailable {
		return
	}
	addServer(m.Guild.ID)
	err := initCommands(s, m.Guild.ID)
	if err != nil {
		log.Printf("failed to init commands in guild '%s'\n%s\n", m.Guild.ID, err)
	}
}

func formatCurrentSong(session *discordgo.Session, status string, server *Server) string {
	queueItem := server.Queue[server.Index]
	return fmt.Sprintf(
		"%s **%s** (%s/%s)",
		status,
		queueItem.Song.Name,
		formatDuration(server.Player.Time),
		formatDuration(queueItem.Song.Length),
	)
}

func formatSong(session *discordgo.Session, status string, server *Server, queueItem *QueueItem) string {
	return fmt.Sprintf(
		"%s **%s** (%s)",
		status,
		queueItem.Song.Name,
		formatDuration(queueItem.Song.Length),
	)
}

func formatDuration(duration time.Duration) string {
	rawSeconds := int(duration.Seconds())
	seconds := rawSeconds % 60
	minutes := rawSeconds / 60 % 60
	hours := rawSeconds / 60 / 60
	if hours > 0 {
		return fmt.Sprintf("%d:%02d:%02d", hours, minutes, seconds)
	}
	return fmt.Sprintf("%d:%02d", minutes, seconds)
}
