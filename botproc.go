package main

import (
	"log"
	"os/exec"
	"syscall"
)

var BotCmd *exec.Cmd
var BotStatus int

const (
	BotRunning = 1 << 0
	BotStopped = 1 << 1
	BotStopping = 1 << 2
)

func StartBot() {
	if BotStatus == BotRunning {
		return
	}

	BotCmd = exec.Command("./cactusbot")
	BotCmd.Dir = "../cactusbot/"

	BotCmd.Stdout = BotWriter{}
	BotCmd.Stderr = BotWriter{}

	err := BotCmd.Start()
	if err != nil {
		log.Println(err)
		return
	}

	BotStatus = BotRunning
	log.Println("Started bot process.")
	
	go func() {
		err = BotCmd.Wait()
		if err != nil {
			log.Println("Error in StartBot go func(): ", err)
		}
		if BotStatus != BotStopping {
			log.Println("Bot just died, attempting to restart it.")
			BotStatus = BotStopped
			StartBot()
		} else {
			BotStatus = BotStopped
			log.Println("Bot process stopped.")
		}
	}()
}

func StopBot() {
	if BotStatus != BotRunning {
		return
	}

	BotStatus = BotStopping
	err := BotCmd.Process.Signal(syscall.SIGTERM)
	if err != nil {
		log.Println("Error in StopBot: ", err)
	}

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

	initialized = true
}
