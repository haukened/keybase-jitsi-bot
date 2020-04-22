package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"

	"samhofi.us/x/keybase/v2/types/chat1"
)

// this JSON pretty prints errors and debug
func p(b interface{}) string {
	s, _ := json.MarshalIndent(b, "", "  ")
	return string(s)
}

// getFeedbackExtendedDescription returns the team name that feedback will be posted to, if configured
func getFeedbackExtendedDescription(bc botConfig) *chat1.UserBotExtendedDescription {
	if bc.FeedbackTeamAdvert != "" {
		return &chat1.UserBotExtendedDescription{
			Title:       "!jitsi feedback",
			DesktopBody: fmt.Sprintf("Please note: Your feedback will be public!\nYour feedback will be posted to %s", bc.FeedbackTeamAdvert),
			MobileBody:  fmt.Sprintf("Please note: Your feedback will be public!\nYour feedback will be posted to %s", bc.FeedbackTeamAdvert),
		}
	}
	return &chat1.UserBotExtendedDescription{
		Title:       fmt.Sprintf("!jitsi feedback"),
		DesktopBody: "Please note: Your feedback will be public!",
		MobileBody:  "Please note: Your feedback will be public!",
	}
}

// hasCommandPrefix determines if the command matches either command or name variant
func hasCommandPrefix(s string, baseCommand string, botName string, subCommands string) bool {
	// if this is actually a command
	if strings.HasPrefix(s, "!") || strings.HasPrefix(s, "@") {
		// generate the two possible command variants
		botCommand := fmt.Sprintf("%s %s", baseCommand, subCommands)
		nameCommand := fmt.Sprintf("%s %s", botName, subCommands)
		// then remove the ! or @ from the string
		s = strings.Replace(s, "!", "", 1)
		s = strings.Replace(s, "@", "", 1)
		// then check if either command variant is a match to the subCommands sent
		if strings.HasPrefix(s, botCommand) || strings.HasPrefix(s, nameCommand) {
			return true
		}
	}
	return false
}

// isRootCommand determines if the command is the root command or name with no arguments
func isRootCommand(s string, baseCommand string, botName string) bool {
	// the space after is important because keybase autocompletes ! and @ with a space after
	botCommand := fmt.Sprintf("!%s ", baseCommand)
	nameCommand := fmt.Sprintf("@%s ", botName)
	if s == botCommand || s == nameCommand {
		return true
	}
	return false
}

// this converts structs to slices of (Name, Value) pairs
func structToSlice(v interface{}) []reflectStruct {
	x := reflect.ValueOf(v)
	values := make([]reflectStruct, x.NumField())
	for i := 0; i < x.NumField(); i++ {
		values[i].Value = x.Field(i).Interface()
		values[i].Name = x.Type().Field(i).Name
	}
	return values
}
