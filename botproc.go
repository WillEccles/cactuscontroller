package main

import (
	"log"
	"os/exec"
	"fmt"
	"syscall"
)

var BotCmd *exec.Cmd
var BotStatus int
var CanCommand bool // whether or not admin commands can be handled at the moment
var proclog []string

const (
	BotRunning = 1 << 0
	BotStopped = 1 << 1
	BotStopping = 1 << 2
)

const MaxOutputBuffer = 15

type BotWriter struct {}

func (bw BotWriter) Write(p []byte) (n int, err error) {
	if len(proclog) != MaxOutputBuffer {
		proclog = append(proclog, string(p))
	} else {
		proclog = append(proclog[1:], string(p))
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

	proclog = nil

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

var initialized bool
func InitProc() {
	if initialized {
		return
	}

	BotStatus = BotStopped
	CanCommand = true

	initialized = true
}
