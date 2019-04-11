package main

import (
	"log"
	"os/exec"
	"fmt"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var BotCmd *exec.Cmd
var BotStatus int
var CanCommand bool // whether or not admin commands can be handled at the moment
var ProcLog []string

const (
	BotRunning = 1 << 0
	BotStopped = 1 << 1
	BotStopping = 1 << 2
)

const MaxOutputBuffer = 15

type BotWriter struct {}

func (bw BotWriter) Write(p []byte) (n int, err error) {
	if len(ProcLog) != MaxOutputBuffer {
		ProcLog = append(ProcLog, string(p))
	} else {
		ProcLog = append(ProcLog[1:], string(p))
	}
	fmt.Print(string(p))
	n = len(p)
	return
}

func StartBot() {
	if BotStatus == BotRunning {
		return
	}

	BotCmd = exec.Command("./cactusbot", "-t", BotToken)
	BotCmd.Dir = "../cactusbot/"

	ProcLog = nil

	BotCmd.Stdout = BotWriter{}
	BotCmd.Stderr = BotWriter{}

	err := BotCmd.Start()
	if err != nil {
		log.Fatal(err)
	}

	BotStatus = BotRunning
}

func StopBot() {
	if BotStatus != BotRunning {
		return
	}

	BotStatus = BotStopping
	err := BotCmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Fatal("Error in StopBot: ", err)
	}

	err = BotCmd.Wait()
	if err != nil {
		log.Fatal("Error in StopBot: ", err)
	}

	BotStatus = BotStopped
}

func RestartBot() {
	StopBot()
	StartBot()
}

func EndProc() {
	StopBot()
}

// this one requires inputs to allow logging in the discord channel
func UpdateBot(msg *discordgo.MessageCreate, s *discordgo.Session) {
	gitcomm := exec.Command("git", "pull")
	gitcomm.Dir = "../cactusbot/"

	bytes, err := gitcomm.CombinedOutput()
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ git clone\n%v\n```\n**Failed.** Restarting bot.", string(bytes)))
		StartBot()
		return
	}
	
	s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ git clone\n%v\n```\n", string(bytes)))

	gocomm := exec.Command("go", "build")
	gocomm.Dir = "../cactusbot/"

	bytes, err = gocomm.CombinedOutput()
	if err != nil {
		s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ go build\n%v\n```\n**Failed!** You'll need to fix it and try again.", string(bytes)))
		return
	}

	s.ChannelMessageSend(msg.ChannelID, fmt.Sprintf("```\n$ go build\n%v\n```\n", string(bytes)))

	s.ChannelMessageSend(msg.ChannelID, "Succeeded in pulling and building, starting bot.")

	StartBot()

	s.ChannelMessageSend(msg.ChannelID, "Bot is started.")
}

var initialized bool
func InitProc() {
	if initialized {
		return
	}

	BotStatus = BotStopped
	CanCommand = true

	initialized = true
}
