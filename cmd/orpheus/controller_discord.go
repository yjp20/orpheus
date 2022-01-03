package main

import (
	"fmt"
	"log"
	"os/exec"
	"strings"
	"time"

	"github.com/yjp20/orpheus/pkg/music"
	"github.com/yjp20/orpheus/pkg/queue"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{Name: "join", Description: "Join channel"},
	{Name: "queue", Description: "Show queue",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "index",
				Description: "location in current queue to show",
				Required:    false,
			},
		},
	},
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
					{Name: "now", Value: queue.Now},
					{Name: "next", Value: queue.Next},
					{Name: "last", Value: queue.Last},
					{Name: "normal", Value: queue.Smart},
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
					{Name: "now", Value: queue.Now},
					{Name: "next", Value: queue.Next},
					{Name: "last", Value: queue.Last},
					{Name: "normal", Value: queue.Smart},
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
					{Name: "one", Value: queue.LoopSong},
					{Name: "all", Value: queue.LoopQueue},
					{Name: "off", Value: queue.NoLoop},
				},
			},
		},
	},
	{Name: "shuffle", Description: "Shuffles the queue after the current song"},
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
	g := GetGuild(m.GuildID)

	switch m.ApplicationCommandData().Name {
	case "add":
		if g.Player.Voice == nil {
			g.joinVoiceOfUser(s, m.GuildID, m.Member.User.ID)
		}
		err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		song, err := music.FetchFromURL(m.ApplicationCommandData().Options[0].StringValue(), false)
		if err != nil {
			log.Printf("failed to add song\nerror: %s\n", err)
			return
		}
		policy := queue.Smart
		if len(m.ApplicationCommandData().Options) > 1 {
			policy = queue.AddPolicy(m.ApplicationCommandData().Options[1].IntValue())
		}
		queueItem := g.Queue.Add(song, s.State.User.ID, false, policy)
		s.InteractionResponseEdit(*appID, m.Interaction, &discordgo.WebhookEdit{
			Content: formatSong("Add", g, queueItem[0]),
		})

	case "addlist":
		err := s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
		})
		songs, err := music.FetchFromURL(m.ApplicationCommandData().Options[0].StringValue(), true)
		if err != nil {
			log.Printf("failed to add song\nerror: %s\n", err)
			return
		}
		shuffle := false
		if len(m.ApplicationCommandData().Options) > 1 {
			shuffle = m.ApplicationCommandData().Options[1].BoolValue()
		}
		queueItems := g.Queue.Add(songs, s.State.User.ID, shuffle, queue.Smart)
		s.InteractionResponseEdit(*appID, m.Interaction, &discordgo.WebhookEdit{
			Content: formatSong("Add", g, queueItems[0]) + fmt.Sprintf(" and %d others", len(queueItems)-1),
		})

	case "queue":
		if len(g.Queue.List) == 0 {
			response = "Queue is empty\n"
			break
		}
		center := g.Queue.Index
		if len(m.ApplicationCommandData().Options) >= 1 {
			center = int(m.ApplicationCommandData().Options[0].IntValue())
		}
		response = PrintQueue(g, center)

	case "pause":
		g.Player.Pause()
		response = formatCurrentSong("Paused", g)

	case "resume":
		if g.Queue.Index == -1 && len(g.Queue.List) > 0 {
			g.Queue.SkipTo(0)
		}
		g.Player.Resume()
		response = formatCurrentSong("Resumed", g)

	case "fastforward":
		if g.Queue.Index == -1 {
			response = "Not playing any song to fast-forward"
		}
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		g.Player.FastForward(seconds)
		response = formatCurrentSong("Fast-forwarded", g)

	case "rewind":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		g.Player.FastForward(-seconds)
		response = formatCurrentSong("Rewound", g)

	case "seek":
		seconds := m.ApplicationCommandData().Options[0].FloatValue()
		if time.Second*time.Duration(seconds) >= g.Queue.List[g.Queue.Index].Song.Length || seconds < 0 {
			response = "Seek value out of range\n"
			break
		}
		g.Player.Seek(seconds)
		response = formatCurrentSong("Seek", g)

	case "skip":
		skip := 1
		if len(m.ApplicationCommandData().Options) >= 1 {
			skip = int(m.ApplicationCommandData().Options[0].IntValue())
		}
		index := (g.Queue.Index + skip) % len(g.Queue.List)
		g.Queue.SkipTo(index)
		response = formatCurrentSong("Skipped to", g)

	case "goto":
		index := int(m.ApplicationCommandData().Options[0].IntValue())
		if index >= len(g.Queue.List) || index < 0 {
			response = "index out of range"
			break
		}
		g.Queue.SkipTo(index)
		response = formatCurrentSong("Go to", g)

	case "remove":
		index := g.Queue.Index
		if len(m.ApplicationCommandData().Options) >= 1 {
			index = int(m.ApplicationCommandData().Options[0].IntValue())
		}
		queueItem, err := g.Queue.Remove(index)
		if err != nil {
			log.Printf("Failed to remove song\nerror: %s\n", err)
		}
		response = formatSong("Remove", g, queueItem)

	case "join":
		response = g.joinVoiceOfUser(s, m.GuildID, m.Member.User.ID)

	case "shuffle":
		g.Queue.Shuffle()
		response = fmt.Sprintf("Shuffled queue")

	case "nowplaying":
		response = formatCurrentSong("Currently Playing: ", g)

	case "loop":
		g.Queue.NextPolicy = queue.NextPolicy(m.ApplicationCommandData().Options[0].IntValue())
		switch g.Queue.NextPolicy {
		case queue.NoLoop:
			response = "Looping turned off"
		case queue.LoopSong:
			response = "Now looping current song"
		case queue.LoopQueue:
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
		song, err := music.FetchFromURL(string(metaData), false)
		if err != nil {
			log.Printf("failed to fetch song\nerror: %s\n", err)
			return
		}
		policy := queue.Smart
		if len(m.ApplicationCommandData().Options) > 1 {
			policy = queue.AddPolicy(m.ApplicationCommandData().Options[1].IntValue())
		}
		queueItem := g.Queue.Add(song, s.State.User.ID, false, policy)
		s.InteractionResponseEdit(*appID, m.Interaction, &discordgo.WebhookEdit{
			Content: formatSong("Add", g, queueItem[0]),
		})

	case "help":
		lines := make([]string, 0)
		lines = append(lines, "Command - Description")
		for _, command := range commands {
			lines = append(lines, fmt.Sprintf("%s - %s", command.Name, command.Description))
		}
		response = strings.Join(lines, "\n")

	case "move":
		from := int(m.ApplicationCommandData().Options[0].IntValue())
		to := int(m.ApplicationCommandData().Options[1].IntValue())
		item, err := g.Queue.Move(from, to)
		if err != nil {
			// TODO handle error
		}
		response = fmt.Sprintf("Move **%s** from index %d to index %d\n", item.Song.Name, from, to)

	case "clear":
		g.Queue.Clear()
		response = "Cleared queue"
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

func (g *Guild) joinVoiceOfUser(s *discordgo.Session, guildID string, userID string) string {
	guild, err := s.State.Guild(guildID)
	if err != nil {
		log.Print(err)
		return fmt.Sprintf("failed to fetch guild '%s'\nerror: %s\n", guildID, err)
	}
	channelID := ""
	for _, vs := range guild.VoiceStates {
		if vs.UserID == userID {
			channelID = vs.ChannelID
		}
	}
	if channelID == "" {
		return "User not in a voice channel"
	}
	channel, err := s.State.Channel(channelID)
	if err != nil {
		log.Print(err)
		return fmt.Sprintf("failed to join channel '%s'\nerror: %s\n", channelID, err)
	}
	g.Player.Voice, err = s.ChannelVoiceJoin(guildID, channelID, false, false)
	if err != nil {
		log.Print(err)
		return fmt.Sprintf("failed to join channel '%s'\nerror: %s\n", channelID, err)
	}
	return fmt.Sprintf("Joining channel '%s'", channel.Name)
}

func joinHandler(s *discordgo.Session, m *discordgo.GuildCreate) {
	if m.Guild.Unavailable {
		return
	}
	GetGuild(m.Guild.ID)
	err := initCommands(s, m.Guild.ID)
	if err != nil {
		log.Printf("failed to init commands in guild '%s'\n%s\n", m.Guild.ID, err)
	}
}

func formatCurrentSong(status string, g *Guild) string {
	if g.Queue.Index == -1 {
		return status
	}
	queueItem := g.Queue.List[g.Queue.Index]
	return fmt.Sprintf(
		"%s **%s** (%s/%s)",
		status,
		queueItem.Song.Name,
		formatDuration(g.Player.Time),
		formatDuration(queueItem.Song.Length),
	)
}

func formatSong(status string, g *Guild, queueItem *queue.QueueItem) string {
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

func PrintQueue(g *Guild, center int) string {
	lines := make([]string, 0)
	for index, queueItem := range g.Queue.List {
		indexString := fmt.Sprintf("%d. ", index)
		if center-10 <= index && index <= center+10 {
			if index == g.Queue.Index {
				lines = append(lines, formatCurrentSong(indexString, g))
			} else {
				lines = append(lines, formatSong(indexString, g, queueItem))
			}
		}
	}
	return strings.Join(lines, "\n")
}
