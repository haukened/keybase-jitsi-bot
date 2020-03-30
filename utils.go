package main

import (
	"encoding/json"
	"fmt"

	"github.com/teris-io/shortid"
	"github.com/ugorji/go/codec"
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

func encodeStructToJSONString(v interface{}) (string, error) {
	jh := codecHandle()
	var bytes []byte
	err := codec.NewEncoderBytes(&bytes, jh).Encode(v)
	if err != nil {
		return "", err
	}
	result := string(bytes)
	return result, nil
}

func decodeJSONStringToStruct(v interface{}, src string) error {
	bytes := []byte(src)
	jh := codecHandle()
	return codec.NewDecoderBytes(bytes, jh).Decode(v)
}

func codecHandle() *codec.JsonHandle {
	var jh codec.JsonHandle
	return &jh
}

func (b *bot) logError(err error) string {
	// generate the error id
	eid := shortid.MustGenerate()
	// send the error to the log
	b.debug("`%s` - %s", eid, err)
	// then return the error id for use
	return eid
}
