package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
)

func shutdownhandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	// TODO
	_, err := s.ChannelMessageSend(msg.ChannelID, "stop")
	if err != nil {
		log.Printf("Error in oodlehandler:\n%v\n", err)
	}
}

func updatehandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	// TODO
	_, err := s.ChannelMessageSend(msg.ChannelID, "update")
	if err != nil {
		log.Printf("Error in oodlehandler:\n%v\n", err)
	}
}

func starthandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	// TODO
	_, err := s.ChannelMessageSend(msg.ChannelID, "start")
	if err != nil {
		log.Printf("Error in oodlehandler:\n%v\n", err)
	}
}

func restarthandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	// TODO
	_, err := s.ChannelMessageSend(msg.ChannelID, "restart")
	if err != nil {
		log.Printf("Error in oodlehandler:\n%v\n", err)
	}
}
