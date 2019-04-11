package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"fmt"
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
	_, err := s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n%v\n```", strings.Join(ProcLog, "")))
	if err != nil {
		log.Printf("Error in loghandler:\n%v\n", err)
	}
}
