package main

import (
	"log"
	"fmt"
	"os/exec"
	"errors"
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

/* These functions aren't used right now, but are here just in case they are needed in the future.
func Backup() (l string, e error) {
	backupcomm := exec.Command("cp", "-v", "cactusbot", "cactusbot.old")
	backupcomm.Dir = "../cactusbot/"
	backupcomm.Stdout = UtilWriter{}
	backupcomm.Stderr = UtilWriter{}
	
	fmt.Println("$ cp -v cactusbot cactusbot.old")
	bytes, err := backupcomm.CombinedOutput()
	l = string(bytes)
	if err != nil {
		e = errors.New("Error backing up!")
		return
	}

	e = nil
	return
}

func Revert() (l string, e error) {
	revertcomm := exec.Command("cp", "-v", "cactusbot.old", "cactusbot")
	revertcomm.Dir = "../cactusbot/"
	revertcomm.Stdout = UtilWriter{}
	revertcomm.Stderr = UtilWriter{}
	
	fmt.Println("$ cp -v cactusbot.old cactusbot")
	bytes, err := revertcomm.CombinedOutput()
	l = string(bytes)
	if err != nil {
		e = errors.New("Error restoring from backup!")
		return
	}

	e = nil
	return
}
*/

func PullRepo() (l string, e error) {
	gitcomm := exec.Command("git", "pull")
	gitcomm.Dir = "../cactusbot/"
	gitcomm.Stdout = UtilWriter{}
	gitcomm.Stderr = UtilWriter{}

	fmt.Println("$ git pull")
	bytes, err := gitcomm.CombinedOutput()
	l = string(bytes)
	if err != nil {
		e = errors.New("Error pulling repository!")
		return
	}
	
	e = nil
	return
}

func BuildBot() (l string, e error) {
	gocomm := exec.Command("go", "build")
	gocomm.Dir = "../cactusbot/"
	gocomm.Stdout = UtilWriter{}
	gocomm.Stderr = UtilWriter{}

	fmt.Println("$ go build")
	bytes, err := gocomm.CombinedOutput()
	l = string(bytes)
	if err != nil {
		e = errors.New("Error building bot!")
		return
	}

	e = nil
	return
}
