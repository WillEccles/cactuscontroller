package main

import (
	"log"
	"os/exec"
	"syscall"
)

var BotCmd *exec.Cmd
var BotStatus int
var CanCommand bool // whether or not admin commands can be handled at the moment

const (
	BotRunning = 1 << 0
	BotStopped = 1 << 1
	BotStopping = 1 << 2
)

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
	log.Println("Started bot process.")
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
	log.Println("Bot process stopped.")
}

func RestartBot() {
	log.Println("Restarting bot process.")
	StopBot()
	StartBot()
}

func EndProc() {
	StopBot()
	log.Println("Cleaned up bot process.")
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
