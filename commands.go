package main

import (
	"fmt"
	"log"

	"samhofi.us/x/keybase/types/chat1"
)

func (b *bot) setupMeeting(convid chat1.ConvIDStr, sender string, words []string, membersType string) {
	b.debug("command recieved in conversation %s", convid)
	meeting, err := newJitsiMeeting()
	if err != nil {
		log.Println(err)
		message := fmt.Sprintf("@%s - I'm sorry, i'm not sure what happened... I was unable to set up a new meeting.\nI've written the appropriate logs and notified my humans.", sender)
		b.k.SendMessageByConvID(convid, message)
		return
	}
	message := fmt.Sprintf("@%s here's your meeting: %s", sender, meeting.getURL())
	b.k.SendMessageByConvID(convid, message)
}
