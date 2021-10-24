package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

var commands = []*discordgo.ApplicationCommand{
	{Name: "/play", Description: ""},
}

func initCommands(s *discordgo.Session, guildId string) {
	for _, command := range commands {
		s.ApplicationCommandCreate(s.State.User.ID, guildId, command)
	}
}

func commandHandler(s *discordgo.Session, m *discordgo.InteractionCreate) {
	if addServer(m.GuildID) {
		initCommands(s, m.GuildID)
	}
	perms, _ := s.State.MessagePermissions(m.Message)
	isAdmin := perms>>discordgo.PermissionAdministrator&1 == 1
	fmt.Println(isAdmin)
}
