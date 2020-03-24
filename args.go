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
	flags.StringVar(&cliConfig.logConvIDStr, "log-convid", "", "sets the keybase chat1.ConvIDStr to log to for feedback")
	if err := flags.Parse(args[1:]); err != nil {
		return err
	}

	// then override the environment vars if there were cli args
	if flags.NFlag() > 0 {
		if cliConfig.Debug == true {
			b.config.Debug = true
		}
		if cliConfig.logConvIDStr != "" {
			b.config.logConvIDStr = cliConfig.logConvIDStr
		}
	}

	// then print the running options
	b.debug("Debug Enabled")
	if b.config.logConvIDStr != "" {
		b.debug("Logging to conversation %s", b.config.logConvIDStr)
	}

	return nil
}
