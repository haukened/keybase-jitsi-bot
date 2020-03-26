package main

import (
	"flag"

	"github.com/caarlos0/env"
)

// parseArgs parses command line and environment args and sets globals
func (b *bot) parseArgs(args []string) error {
	// parse the env variables into the bot config
	if err := env.Parse(&b.config); err != nil {
		return err
	}

	// then parse CLI args as overrides
	flags := flag.NewFlagSet(args[0], flag.ExitOnError)
	cliConfig := botConfig{}
	flags.BoolVar(&cliConfig.Debug, "debug", false, "enables command debugging to stdout")
	flags.StringVar(&cliConfig.LogConvIDStr, "log-convid", "", "sets the keybase chat1.ConvIDStr to log debugging to keybase chat.")
	flags.StringVar(&cliConfig.FeedbackConvIDStr, "feedback-convid", "", "sets the keybase chat1.ConvIDStr to send feedback to.")
	flags.StringVar(&cliConfig.FeedbackTeamAdvert, "feedback-team-advert", "", "sets the keybase team/channel to advertise feedback. @team#channel")
	flags.StringVar(&cliConfig.KVStoreTeam, "kvstore-team", "", "sets the keybase team where kvstore values are stored")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// then override the environment vars if there were cli args
	if flags.NFlag() > 0 {
		if cliConfig.Debug == true {
			b.config.Debug = true
		}
		if cliConfig.LogConvIDStr != "" {
			b.config.LogConvIDStr = cliConfig.LogConvIDStr
		}
		if cliConfig.FeedbackConvIDStr != "" {
			b.config.FeedbackConvIDStr = cliConfig.FeedbackConvIDStr
		}
		if cliConfig.FeedbackTeamAdvert != "" {
			b.config.FeedbackTeamAdvert = cliConfig.FeedbackTeamAdvert
		}
		if cliConfig.KVStoreTeam != "" {
			b.config.KVStoreTeam = cliConfig.KVStoreTeam
		}
	}

	// then print the running options
	b.debug("Debug Enabled")
	if b.config.LogConvIDStr != "" {
		b.debug("Logging to conversation %s", b.config.LogConvIDStr)
	}
	if b.config.FeedbackConvIDStr != "" {
		b.debug("Feedback enabled to %s and advertising %s", b.config.FeedbackConvIDStr, b.config.FeedbackTeamAdvert)
	}
	if b.config.KVStoreTeam != "" {
		b.debug("keybase kvstore enabled in @%s", b.config.KVStoreTeam)
	}

	return nil
}
