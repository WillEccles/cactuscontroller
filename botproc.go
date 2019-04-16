package main

import (
	"log"
	"fmt"
	"os/exec"
	"syscall"
	"github.com/bwmarrin/discordgo"
)

var BotCmd *exec.Cmd
var BotStatus int

const (
	BotRunning = 1 << 0
	BotStopped = 1 << 1
	BotStopping = 1 << 2
)

func StartBot(s *discordgo.Session) {
	if BotStatus == BotRunning {
		return
	}

	// wait for the bot to shut down
	if BotStatus == BotStopping {
		log.Println("StartBot(): Waiting for bot to shutdown...")
		for BotStatus == BotStopping { }
		log.Println("StartBot(): Bot shutdown.")
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
			s.ChannelMessageSend(DebugChannel, fmt.Sprintf("Bot just crashed, trying to restart! <@%v>", AdminID))
			loghandler(&discordgo.MessageCreate{
				&discordgo.Message{
					ChannelID: DebugChannel,
					Content: "cc log bot",
				},
			}, s)
			BotStatus = BotStopped
			StartBot(s)
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

func RestartBot(s *discordgo.Session) {
	log.Println("Restarting bot process.")
	StopBot()
	StartBot(s)
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
