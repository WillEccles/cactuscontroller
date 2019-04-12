package main

import (
	"log"
	"fmt"
	"os/exec"
	"errors"
	"strings"
	"github.com/bwmarrin/discordgo"
)

// this one requires inputs to allow logging in the discord channel
func UpdateBot(msg *discordgo.MessageCreate, s *discordgo.Session) {
	log.Println("Attempting to update bot...")
	
	log.Println("Pulling repository...")
	bytes, err := PullRepo()
	if err != nil {
		log.Println("Error pulling repository! Will restart bot.")
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ git pull\n%v\n```\n**Failed.** Restarting bot.", string(bytes)))
		StartBot()
		return
	}
	
	log.Println("Successfully pulled from repository.")
	s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ git clone\n%v\n```\n", string(bytes)))

	log.Println("Building bot...")
	bytes, err = BuildBot()
	if err != nil {
		log.Println("Error building bot! Will restart bot.")
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ go build\n%v\n```\n**Failed.** Restarting bot.", string(bytes)))
		StartBot()
		return
	}

	log.Println("Successfully built bot.")
	s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ go build\n%v\n```\n", string(bytes)))

	log.Println("Upgrade successful, restarting bot...")
	s.ChannelMessageSend(msg.ChannelID, "Upgrade successful, starting bot.")

	StartBot()
}

func PullRepo() (l string, e error) {
	gitcomm := exec.Command("git", "pull")
	gitcomm.Dir = "../cactusbot/"

	var combinedlog []string
	
	gitcomm.Stdout = UtilWriter{
		Log: &combinedlog,
	}
	gitcomm.Stderr = UtilWriter{
		Log: &combinedlog,
	}

	fmt.Println("$ git pull")
	err := gitcomm.Run()
	l = strings.Join(combinedlog, "")
	if err != nil {
		e = errors.New(fmt.Sprintf("Error pulling repository: %v", err))
		return
	}
	
	e = nil
	return
}

func BuildBot() (l string, e error) {
	gocomm := exec.Command("go", "build")
	gocomm.Dir = "../cactusbot/"

	var combinedlog []string
	
	gocomm.Stdout = UtilWriter{
		Log: &combinedlog,
	}
	gocomm.Stderr = UtilWriter{
		Log: &combinedlog,
	}

	fmt.Println("$ go build")
	err := gocomm.Run()
	l = strings.Join(combinedlog, "")
	if err != nil {
		e = errors.New(fmt.Sprintf("Error building bot: %v", err))
		return
	}

	e = nil
	return
}
