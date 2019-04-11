package main

import (
	"regexp"
	"github.com/bwmarrin/discordgo"
)

type MsgHandler func(*discordgo.MessageCreate, *discordgo.Session)

type Command struct {
	Pattern *regexp.Regexp
	Name string
	Description string
	Aliases []string
	Handler MsgHandler
}

func (cmd *Command) Handle(msg *discordgo.MessageCreate, s *discordgo.Session) {
	cmd.Handler(msg, s)
}

var Commands = []Command {
	{
		Name: "shutdown",
		Aliases: []string{
			"stop",
		},
		Description: "Shuts down the bot.",
		Pattern: regexp.MustCompile(`(?i)^cc\s+(shutdown|stop)`),
		Handler: shutdownhandler,
	},
	{
		Name: "update",
		Description: "Updates and rebuilds the bot.",
		Aliases: []string{
			"upgrade",
		},
		Pattern: regexp.MustCompile(`(?i)^cc\s+up(date|grade)`),
		Handler: updatehandler,
	},
	{
		Name: "start",
		Description: "Starts up the bot.",
		Pattern: regexp.MustCompile(`(?i)^cc\s+start`),
		Handler: starthandler,
	},
	{
		Name: "restart",
		Description: "Restarts the bot.",
		Pattern: regexp.MustCompile(`(?i)^cc\s+restart`),
		Handler: restarthandler,
	},
}
