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
	if !CanCommand {
		return
	}
	CanCommand = false

	cmd.Handler(msg, s)

	CanCommand = true
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
	{
		Name: "log <option>",
		Description: "`<option>` can be one of:\n- `bot`: shows the last 40 bot log messages\n- `controller`: shows the last 40 lines of controller logs\n- `all`: shows the last 40 bot log messages, followed by the last 40 controller log messages\n- `<none>`: same as `all`",
		Pattern: regexp.MustCompile(`(?i)^cc\s+log(\s+(bot|controller|all))?`),
		Handler: loghandler,
	},
}
