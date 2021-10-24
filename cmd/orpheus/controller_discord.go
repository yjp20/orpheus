package main

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func messageHandler(s *discordgo.Session, m *discordgo.MessageCreate) {
	if m.Author.ID == s.State.User.ID {
		return
	}
	perms, _ := s.State.MessagePermissions(m.Message)
	isAdmin := perms >> discordgo.PermissionAdministrator & 1 == 1
	fmt.Println(isAdmin)
}
