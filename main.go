package main

import (
	"fmt"
	"log"
	"os"

	"samhofi.us/x/keybase"
	"samhofi.us/x/keybase/types/chat1"
	"samhofi.us/x/keybase/types/stellar1"
)

// this global controls debug printing
var debug bool

// Bot holds the necessary information for the bot to work.
type bot struct {
	k        *keybase.Keybase
	handlers keybase.Handlers
	opts     keybase.RunOptions
	payments map[stellar1.PaymentID]botReply
	config   botConfig
}

// botConfig hold env and cli flags and options
// fields must be exported for package env (reflect) to work
type botConfig struct {
	Debug              bool   `env:"BOT_DEBUG" envDefault:"false"`
	LogConvIDStr       string `env:"BOT_LOG_CONVID" envDefault:""`
	FeedbackConvIDStr  string `env:"BOT_FEEDBACK_CONVID" envDefault:""`
	FeedbackTeamAdvert string `env:"BOT_FEEDBACK_TEAM_ADVERT" envDefault:""`
}

// hold reply information when needed
type botReply struct {
	convID chat1.ConvIDStr
	msgID  chat1.MessageID
}

// Debug provides printing only when --debug flag is set or BOT_DEBUG env var is set
func (b *bot) debug(s string, a ...interface{}) {
	if b.config.Debug {
		log.Printf(s, a...)
		if b.config.LogConvIDStr != "" {
			b.logToChat(s, a...)
		}
	}
}

// logToChat will send this message to the keybase chat configured in b.logConv
func (b *bot) logToChat(s string, a ...interface{}) {
	// if the ConvIdStr isn't blank try to log
	if b.config.LogConvIDStr != "" {
		// if you can't send the message, log the error to stdout
		if _, err := b.k.SendMessageByConvID(chat1.ConvIDStr(b.config.LogConvIDStr), s, a...); err != nil {
			log.Printf("Unable to log to keybase chat: %s", err)
		}
	} else {
		// otherwise (and you shouldn't be here but....) log it to stdout
		log.Println("Unable to log to keybase chat, logging ConvIDStr is not set")
	}
}

// newBot returns a new empty bot
func newBot() *bot {
	var b bot
	b.k = keybase.NewKeybase()
	b.handlers = keybase.Handlers{}
	b.opts = keybase.RunOptions{}
	b.payments = make(map[stellar1.PaymentID]botReply)
	return &b
}

// this handles setting up command advertisements and aliases
func (b *bot) registerCommands() {
	opts := keybase.AdvertiseCommandsOptions{
		Alias: "Jitsi Meet",
		Advertisements: []chat1.AdvertiseCommandAPIParam{
			{
				Typ: "public",
				Commands: []chat1.UserBotCommandInput{
					{
						Name:        "jitsi",
						Description: "Starts a meet.jit.si meeting",
						Usage:       "",
					},
					{
						Name:                fmt.Sprintf("%s feedback", b.k.Username),
						Description:         "Tell us how we're doing!",
						Usage:               "",
						ExtendedDescription: getFeedbackExtendedDescription(b.config),
					},
				},
			},
		},
	}
	b.k.AdvertiseCommands(opts)
}

// run performs a proxy main function
func (b *bot) run(args []string) error {
	// parse the arguments
	err := b.parseArgs(args)
	if err != nil {
		return err
	}

	b.registerHandlers()
	// clear the commands and advertise the new commands
	b.k.ClearCommands()
	b.registerCommands()

	log.Println("Starting...")
	b.k.Run(b.handlers, &b.opts)
	return nil
}

// main is a thin skeleton, proxied to Bot.Run()
func main() {
	b := newBot()
	if err := b.run(os.Args); err != nil {
		log.Printf("%s", err)
		os.Exit(1)
	}
}
