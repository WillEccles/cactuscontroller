package main

import (
	"log"
	"fmt"
)

const MaxLogBuffer = 40 // max lines of controller/bot logs to store (each)
const CharacterLimit = 2000 // max message length on discord

type BotWriter struct {}
type UtilWriter struct {
	Log *[]string
}
type ControllerLogger struct {}

var ProcLog []string
var ControllerLog []string

func (bw BotWriter) Write(p []byte) (n int, err error) {
	if len(ProcLog) != MaxLogBuffer {
		ProcLog = append(ProcLog, string(p))
	} else {
		ProcLog = append(ProcLog[1:], string(p))
	}
	fmt.Print(string(p))
	n = len(p)
	return
}

func (uw UtilWriter) Write(p []byte) (n int, err error) {
	if len(ControllerLog) != MaxLogBuffer {
		ControllerLog = append(ControllerLog, "  " + string(p))
	} else {
		ControllerLog = append(ControllerLog[1:], "  " + string(p))
	}
	fmt.Printf("  %s", string(p))
	*(uw.Log) = append(*(uw.Log), string(p))
	n = len(p)
	return
}

func (cl ControllerLogger) Write(p []byte) (n int, err error) {
	if len(ControllerLog) != MaxLogBuffer {
		ControllerLog = append(ControllerLog, string(p))
	} else {
		ControllerLog = append(ControllerLog[1:], string(p))
	}
	fmt.Print(string(p))
	n = len(p)
	return
}

func InitLogger() {
	log.SetPrefix("[Controller] ")
	log.SetOutput(ControllerLogger{})
	log.Println("Initialized logger.")
}

// gets the last 40 controller logs as a slice of discord-friendly formatted strings
// each string may only be up to 2000 lines
func GetControllerLogs() (out []string) {
	return getLogs(ControllerLog)
}

// same as GetControllerLogs but for the bot
func GetBotLogs() (out []string) {
	return getLogs(ProcLog)
}

// backend for the above two functions
func getLogs(in []string) (out []string) {
	charcount := 30 // start with 30 just to be sure
	tempstr := ""
	for _, s := range(in) {
		if charcount + len(s) > CharacterLimit {
			out = append(out, fmt.Sprintf("```\n%s```", tempstr))
			charcount = 30
			tempstr = ""
		}
		charcount += len(s)
		tempstr += s
	}
	// one more to catch the last little bit
	out = append(out, fmt.Sprintf("```\n%s```", tempstr))
	return
}
