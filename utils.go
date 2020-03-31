package main

import (
	"encoding/json"
	"fmt"

	"samhofi.us/x/keybase/types/chat1"
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
