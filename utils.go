package main

import (
	"encoding/json"
	"fmt"

	"github.com/teris-io/shortid"
	"samhofi.us/x/keybase/types/chat1"
)

// this JSON pretty prints errors and debug
func p(b interface{}) string {
	s, _ := json.MarshalIndent(b, "", "  ")
	return string(s)
}

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

func (b *bot) logError(err error) string {
	// generate the error id
	eid := shortid.MustGenerate()
	// send the error to the log
	b.log("`%s` - %s", eid, err)
	// then return the error id for use
	return eid
}
