package main

import (
	"flag"
	"fmt"
	"syscall"
	"os"
	"os/signal"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

const (
	/*Perms = 8 // admin
	ClientID = "565946505026863145"
	InvURL = "https://discordapp.com/oauth2/authorize?&client_id=%v&scope=bot&permissions=%v"*/
	ConsoleChannel = "565947628718915587"
	DebugChannel = "245649734302302208"
	AdminID = "111943010396229632"
	BotID = "237605108173635584"
	AdminGuild = "111944249544564736"
)

func init() {
	flag.StringVar(&token, "t", "", "Controller Token")
	flag.StringVar(&BotToken, "b", "", "Bot Token")
	flag.Parse()
}

var token string
var BotToken string
var HelpEmbed discordgo.MessageEmbed
var SigChan chan os.Signal
var firstReady bool

func main() {
	firstReady = true
	InitLogger()

	if token == "" {
		log.Println("No controller token provided. Please run: cactuscontroller -t <controllertoken> -b <bottoken>")
		return
	}

	if BotToken == "" {
		log.Println("No child token supplied. Please run: cactuscontroller -t <controllertoken> -b <bottoken>")
		return
	}

	InitProc()

	// prepare a help embed to reduce CPU load later on
	HelpEmbed.Title = "**Here's what I can do!**"
	HelpEmbed.Description = "You should begin each command with `cc`.\nFor example: `cc upgrade`."

	for _, cmd := range(Commands) {
		newfield := discordgo.MessageEmbedField{
			Name: "**`" + cmd.Name + "`**",
			Value: cmd.Description,
			Inline: false,
		}
		if len(cmd.Aliases) != 0 {
			if len(cmd.Aliases) == 1 {
				newfield.Value += "\n**Alias:** "
			} else {
				newfield.Value += "\n**Aliases:** "
			}
			for i, a := range(cmd.Aliases) {
				if i > 0 {
					newfield.Value += ", "
				}
				newfield.Value += "`" + a + "`"
			}
		}
		HelpEmbed.Fields = append(HelpEmbed.Fields, &newfield)
	}

	dg, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Println("Error creating Discord session: ", err)
		return
	}

	dg.AddHandler(ready)
	dg.AddHandler(messageCreate)
	dg.AddHandler(connect)
	dg.AddHandler(resume)
	dg.AddHandler(disconnect)

	err = dg.Open()
	if err != nil {
		log.Println("Error opening Discord session: ", err)
		return
	}
	defer log.Println("Goodbye.")
	defer dg.Close() // close the session after Control-C
	defer EndProc()

	SigChan = make(chan os.Signal)
	signal.Notify(SigChan, syscall.SIGINT, syscall.SIGTERM, os.Interrupt, os.Kill)
	<-SigChan
}

func ready(s *discordgo.Session, event *discordgo.Ready) {
	log.Println("Controller ready.")

	i := 0
	usd := discordgo.UpdateStatusData{
		IdleSince: &i,
		AFK: false,
		Status: "online",
		Game: &discordgo.Game {
			Name: "over cactusbot",
			Type: discordgo.GameTypeWatching,
		},
	}

	err := s.UpdateStatusComplex(usd)
	if err != nil {
		log.Printf("Error in ready:\n%v\n", err)
	}

	if firstReady {
		firstReady = false

		// start the bot since we know we should now
		StartBot()
	}
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// only listen to me
	if m.Author.ID != AdminID || m.ChannelID != ConsoleChannel {
		return
	}

	for _, cmd := range(Commands) {
		if cmd.Pattern.MatchString(m.Content) {
			cmd.Handle(m, s)
			break
		}
	}
}

func connect(s *discordgo.Session, event *discordgo.Connect) {
	log.Println("Controller connected.")
}

func disconnect(s *discordgo.Session, event *discordgo.Disconnect) {
	log.Println("Controller disconnected!")
}

func resume(s *discordgo.Session, event *discordgo.Resumed) {
	log.Println("Controller resumed, attempting to send debug message.")
	_, err := s.ChannelMessageSend(DebugChannel, fmt.Sprintf("Just recovered from error(s)! <@%v>\n```\n%v\n```", AdminID, strings.Join(event.Trace, "\n")))
	if err != nil {
		log.Printf("Error in resume (this is awkward):\n%v\n", err)
	}
}
