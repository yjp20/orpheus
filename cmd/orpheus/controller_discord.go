package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{Name: "join", Description: "Join channel"},
	{Name: "queue", Description: "Show queue"},
	{Name: "pause", Description: "Pause playing"},
	{Name: "resume", Description: "Resume playing"},
	{Name: "add", Description: "Adds a single song to the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "link to music",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "when",
				Description: "when to play queued song",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "now",
						Value: Now,
					},
					{
						Name:  "next",
						Value: Next,
					},
					{
						Name:  "last",
						Value: Last,
					},
					{
						Name:  "normal",
						Value: Smart,
					},
				},
			},
		},
	},
	{Name: "addlist", Description: "Adds a playlist to the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "url",
				Description: "link to music",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionBoolean,
				Name:        "shuffle",
				Description: "shuffle playlist before queueing",
				Required:    false,
			},
		},
	},
	{Name: "search", Description: "Searches Youtube for a song and queues the first result",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "song",
				Description: "name of song",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "when",
				Description: "when to play queued song",
				Required:    false,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "now",
						Value: Now,
					},
					{
						Name:  "next",
						Value: Next,
					},
					{
						Name:  "last",
						Value: Last,
					},
					{
						Name:  "normal",
						Value: Smart,
					},
				},
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
	{Name: "skip", Description: "Skips a song (or multiple songs) in the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "index",
				Description: "number of indices to skip",
				Required:    false,
			},
		},
	},
	{Name: "goto", Description: "Plays the song at a given index",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "index",
				Description: "index of song to play",
				Required:    true,
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
	{Name: "move", Description: "Moves a song in the queue to a desired position",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "from",
				Description: "index of song to move",
				Required:    true,
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "to",
				Description: "new index of song",
				Required:    true,
			},
		},
	},
	{Name: "loop", Description: "Change the looping policy of the queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "mode",
				Description: "looping mode of the queue",
				Required:    true,
				Choices: []*discordgo.ApplicationCommandOptionChoice{
					{
						Name:  "one",
						Value: LoopSong,
					},
					{
						Name:  "all",
						Value: LoopQueue,
					},
					{
						Name:  "off",
						Value: NoLoop,
					},
				},
			},
		},
	},
	{Name: "shuffle", Description: "Shuffles the queue"},
	{Name: "nowplaying", Description: "Shows the currently playing song"},
	{Name: "clear", Description: "Removes all songs from the queue and stops playing"},
	{Name: "help", Description: "Prints all available commands"},
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
	server := getServer(m.GuildID)

	switch m.ApplicationCommandData().Name {
	case "add":
		err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		song, err := fetchSongsFromURL(m.ApplicationCommandData().Options[0].StringValue(), false)
		if err != nil {
			log.Printf("failed to add song\nerror: %s\n", err)
			return
		}
		policy := Smart
		if len(m.ApplicationCommandData().Options) > 1 {
			policy = addPolicy(m.ApplicationCommandData().Options[1].IntValue())
		}
		queueItem := server.Add(song, s.State.User.ID, false, policy)
		s.InteractionResponseEdit(*appID, m.Interaction, &discordgo.WebhookEdit{
			Content: formatSong("Add", server, queueItem[0]),
		})

	case "addlist":
		err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		songs, err := fetchSongsFromURL(m.ApplicationCommandData().Options[0].StringValue(), true)
		if err != nil {
			log.Printf("failed to add song\nerror: %s\n", err)
			return
		}
		shuffle := false
		if len(m.ApplicationCommandData().Options) > 1 {
			shuffle = m.ApplicationCommandData().Options[1].BoolValue()
		}
		queueItems := server.Add(songs, s.State.User.ID, shuffle, Smart)
		s.InteractionResponseEdit(*appID, m.Interaction, &discordgo.WebhookEdit{
			Content: formatSong("Add", server, queueItems[0]) + fmt.Sprintf(" and %d others", len(queueItems)-1),
		})

	case "queue":
		if len(server.Queue) == 0 {
			response = "Queue is empty\n"
			break
		}
		response = PrintQueue(server)

	case "pause":
		server.Player.Pause()
		response = formatCurrentSong("Pause", server)

	case "resume":
		server.Player.Resume()
		response = formatCurrentSong("Resumed", server)

	case "fastforward":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		server.Player.FastForward(seconds)
		response = formatCurrentSong("Fast-forward", server)

	case "rewind":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		server.Player.FastForward(-seconds)
		response = formatCurrentSong("Rewind", server)

	case "seek":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		server.Player.Seek(seconds)
		response = formatCurrentSong("Seek", server)

	case "skip":
		skip := 1
		if len(m.ApplicationCommandData().Options) >= 1 {
			skip = int(m.ApplicationCommandData().Options[0].IntValue())
		}
		index := (server.Index + skip) % len(server.Queue)
		server.SkipTo(index)
		response = formatCurrentSong("Skip to", server)

	case "goto":
		index := int(m.ApplicationCommandData().Options[0].IntValue())
		if index >= len(server.Queue) || index < 0 {
			response = "index out of range"
			break
		}
		server.SkipTo(index)
		response = formatCurrentSong("Skip to", server)

	case "remove":
		index := server.Index
		if len(m.ApplicationCommandData().Options) >= 1 {
			index = int(m.ApplicationCommandData().Options[0].IntValue())
		}
		queueItem, err := server.Remove(index)
		if err != nil {
			log.Printf("Failed to remove song\nerror: %s\n", err)
		}
		response = formatSong("Remove", server, queueItem)

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

	case "shuffle":
		server.Shuffle()
		response = fmt.Sprintf("Shuffling the queue...")

	case "nowplaying":
		response = formatCurrentSong("Currently Playing: ", server)

	case "loop":
		server.NextPolicy = NextPolicy(m.ApplicationCommandData().Options[0].IntValue())
		switch server.NextPolicy {
		case NoLoop:
			response = "Looping turned off"
		case LoopSong:
			response = "Now looping current song"
		case LoopQueue:
			response = "Now looping queue"
		}

	case "search":
		err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		metaDataProcess := exec.Command("yt-dlp", "--default-search", "auto", "--print", "%(id)s", m.ApplicationCommandData().Options[0].StringValue())
		metaData, err := metaDataProcess.Output()
		if err != nil || string(metaData) == "" {
			log.Printf("failed to load search result matching \"%s\"\n", m.ApplicationCommandData().Options[0].StringValue())
			return
		}
		song, err := fetchSongsFromURL(string(metaData), false)
		if err != nil {
			log.Printf("failed to fetch song\nerror: %s\n", err)
			return
		}
		policy := Smart
		if len(m.ApplicationCommandData().Options) > 1 {
			policy = addPolicy(m.ApplicationCommandData().Options[1].IntValue())
		}
		queueItem := server.Add(song, s.State.User.ID, false, policy)
		s.InteractionResponseEdit(*appID, m.Interaction, &discordgo.WebhookEdit{
			Content: formatSong("Add", server, queueItem[0]),
		})

	case "help":
		lines := make([]string, 0)
		lines = append(lines, "Command - Description")
		for _, command := range commands {
			lines = append(lines, fmt.Sprintf("%s - %s", command.Name, command.Description))
		}
		response = strings.Join(lines, "\n")

	case "move":
		server.Move(int(m.ApplicationCommandData().Options[0].IntValue()), int(m.ApplicationCommandData().Options[1].IntValue()))
		response = fmt.Sprintf("Move **%s** from index %d to index %d\n", server.Queue[m.ApplicationCommandData().Options[1].IntValue()].Song.Name,
			m.ApplicationCommandData().Options[0].IntValue(), m.ApplicationCommandData().Options[1].IntValue())

	case "clear":
		server.Clear()
		response = "Queue has been cleared\n"
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
	getServer(m.Guild.ID)
	err := initCommands(s, m.Guild.ID)
	if err != nil {
		log.Printf("failed to init commands in guild '%s'\n%s\n", m.Guild.ID, err)
	}
}

func formatCurrentSong(status string, server *Server) string {
	queueItem := server.Queue[server.Index]
	return fmt.Sprintf(
		"%s **%s** (%s/%s)",
		status,
		queueItem.Song.Name,
		formatDuration(server.Player.Time),
		formatDuration(queueItem.Song.Length),
	)
}

func formatSong(status string, server *Server, queueItem *QueueItem) string {
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

func PrintQueue(server *Server) string {
	lines := make([]string, 0)
	for index, queueItem := range server.Queue {
		indexString := fmt.Sprintf("%d. ", index)
		if index == server.Index {
			lines = append(lines, formatCurrentSong(indexString, server))
		} else {
			lines = append(lines, formatSong(indexString, server, queueItem))
		}
	}
	return strings.Join(lines, "\n")
}
