package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"regexp"
	"fmt"
	"time"
)

func shutdownhandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	m := ""
	if BotStatus != BotRunning {
		m = "Bot is already stopped."
	} else {
		m = "Stopping bot..."
	}
	_, err := s.ChannelMessageSend(msg.ChannelID, m)
	if err != nil {
		log.Printf("Error in shutdownhandler:\n%v\n", err)
	}

	if BotStatus != BotRunning {
		return
	}

	StopBot()

	_, err = s.ChannelMessageSend(msg.ChannelID, "Stopped bot.")
	if err != nil {
		log.Printf("Error in stophandler:\n%v\n", err)
	}
}

func updatehandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	if BotStatus == BotRunning {
		shutdownhandler(msg, s)
	}

	UpdateBot(msg, s)
}

func starthandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	m := ""
	if BotStatus == BotRunning {
		m = "Bot is already running."
	} else {
		m = "Starting bot..."
	}
	
	_, err := s.ChannelMessageSend(msg.ChannelID, m)
	if err != nil {
		log.Printf("Error in starthandler:\n%v\n", err)
	}

	if BotStatus == BotRunning {
		return
	}

	StartBot()

	_, err = s.ChannelMessageSend(msg.ChannelID, "Started bot.")
	if err != nil {
		log.Printf("Error in starthandler:\n%v\n", err)
	}
}

func restarthandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	m := ""
	if BotStatus == BotRunning {
		m = "Restarting bot..."
	} else {
		m = "Bot is already off, starting bot..."
	}
	_, err := s.ChannelMessageSend(msg.ChannelID, m)
	if err != nil {
		log.Printf("Error in restarthandler:\n%v\n", err)
	}

	if BotStatus == BotRunning {
		StopBot()
		StartBot()

		_, err := s.ChannelMessageSend(msg.ChannelID, "Restarted bot.")
		if err != nil {
			log.Printf("Error in restarthandler:\n%v\n", err)
		}
	} else {
		StartBot()
		
		_, err := s.ChannelMessageSend(msg.ChannelID, "Started bot.")
		if err != nil {
			log.Printf("Error in restarthandler:\n%v\n", err)
		}
	}

}

func loghandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	re := regexp.MustCompile(`(?i)^cc\s+log\s*`)
	clean := re.ReplaceAllString(msg.Content, "")
	clean = strings.TrimSpace(strings.ToLower(clean))
	
	dobot := false
	docontroller := false

	if clean == "" || clean == "all" {
		dobot = true
		docontroller = true
	} else if clean == "bot" {
		dobot = true
	} else if clean == "controller" {
		docontroller = true
	}

	if docontroller {
		cl := GetControllerLogs()

		for i, str := range(cl) {
			partno := ""
			if len(cl) > 1 {
				partno = fmt.Sprintf(" (part %v/%v)", i+1, len(cl))
			}
			_, err := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Controller Logs%s:\n%s", partno, str))
			if err != nil {
				log.Printf("Error sending message in loghandler:\n%v\n", err)
			}
			time.Sleep(500 * time.Millisecond) // just to make sure we don't get ratelimited or something
		}
	}

	if dobot {
		bl := GetBotLogs()
		
		for i, str := range(bl) {
			partno := ""
			if len(bl) > 1 {
				partno = fmt.Sprintf(" (part %v/%v)", i+1, len(bl))
			}
			_, err := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Bot Logs%s:\n%s", partno, str))
			if err != nil {
				log.Printf("Error sending message in loghandler:\n%v\n", err)
			}
			time.Sleep(500 * time.Millisecond) // just to make sure we don't get ratelimited or something
		}
	}
}

func helphandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	embedcolor := s.State.UserColor(s.State.User.ID, msg.ChannelID)
	embed := HelpEmbed
	embed.Color = embedcolor

	_, err := s.ChannelMessageSendEmbed(msg.ChannelID, &embed)
	if err != nil {
		log.Printf("Error in helphandler:\n%v\n", err)
	}
}
