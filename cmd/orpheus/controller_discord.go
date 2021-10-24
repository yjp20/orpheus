package main

import (
    "strconv"
    "github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
    {Name: "add", Description: "Adds a song to the queue",
    Options: []*discordgo.ApplicationCommandOption{
        {
            Type: discordgo.ApplicationCommandOptionString,
            Name: "url",
            Description: "link to music",
            Required: true,
        },
    },
    },
    {Name: "queue", Description: "Show queue",},
    /*
    {Name: "skip", Description: "skips a song in the queue",
    Options: []*discordgo.ApplicationCommandOption{
        {
            Type: discordgo.ApplicationCommandOptionInteger,
            Name: "index",
            Description: "index to skip to",
            Required: false,
        },
    },
    },
    {Name: "remove", Description: "removes a song in the queue",
    Options: []*discordgo.ApplicationCommandOption{
        {
            Type: discordgo.ApplicationCommandOptionInteger,
            Name: "index",
            Description: "index of song to remove",
            Required: true,
        },
    },
    */
}

func initCommands(s *discordgo.Session, guildId string) {
    for _, command := range commands {
        s.ApplicationCommandCreate(s.State.User.ID, guildId, command)
    }
}

func commandHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
    println("HI")
    if addServer(m.GuildID) {
        initCommands(s, m.GuildID)
    }

    //perms, _ := s.State.MessagePermissions(m.Message)
    //perms>>discordgo.PermissionAdministrator&1 == 1

    switch m.ApplicationCommandData().Name {
    case "add":
        add(m.GuildID, m.ApplicationCommandData().Options[0].StringValue(), s.State.User.ID)
    case "queue":
        server, _ := servers[m.GuildID]
        content := ""
        for i, v := range server.Queue {
            content += strconv.Itoa(i) + ".  " + v.Song.Name + "\n"
        }
        s.InteractionRespond(m.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData {
                Content: content,
            },
        })
    }
}

