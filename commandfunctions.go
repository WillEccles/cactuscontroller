package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"strings"
	"regexp"
	"fmt"
	"time"
	"io/ioutil"
	"os"
	"encoding/json"
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

	StartBot(s)

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
		StartBot(s)

		_, err := s.ChannelMessageSend(msg.ChannelID, "Restarted bot.")
		if err != nil {
			log.Printf("Error in restarthandler:\n%v\n", err)
		}
	} else {
		StartBot(s)
		
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

func confhandler(msg *discordgo.MessageCreate, s *discordgo.Session) {
	re := regexp.MustCompile(`(?i)^cc\s+conf\s+`)
	clean := re.ReplaceAllString(msg.Content, "")
	if strings.ToLower(strings.Fields(clean)[0]) == "get" {
		file, err := os.Open("../cactusbot/config.json")
		if err != nil {
			log.Printf("Error opening config.json:\n%v\n", err)
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error opening config.json:\n```\n%v\n```", err))
			return
		}
		defer file.Close()
		decoder := json.NewDecoder(file)
		conf := make(map[string]interface{})
		err = decoder.Decode(&conf)
		if err != nil {
			log.Printf("Error decoding json:\n%v\n", err)
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error decoding config.json:\n```\n%v\n```", err))
			return
		}
		b, err := json.MarshalIndent(conf, "", "    ")
		if err != nil {
			log.Printf("Error mashalling JSON:\n%v\n", err)
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error marshalling json:\n```\n%v\n```", err))
			return
		}
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Contents of `config.json`:\n```json\n%v\n```", string(b)))
	} else {
		re2 := regexp.MustCompile("(?i)set[\\s\\n\\t]+(```(json)?\\n?)?")
		clean2 := re2.ReplaceAllString(clean, "")
		re3 := regexp.MustCompile("(?i)\\n?```")
		clean3 := re3.ReplaceAllString(clean2, "")

		err := ioutil.WriteFile("../cactusbot/config.json", []byte(clean3), 0644)
		if err != nil {
			log.Printf("Error writing file:\n%v\n", err)
			s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("Error writing file:\n%v\n", err))
			return
		}
		s.ChannelMessageSend(msg.ChannelID, "Wrote `config.json`. Remember to do `cc restart` to update the bot's config!")
	}
}
