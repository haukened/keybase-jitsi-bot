package main

import (
	"fmt"
	"log"

	"samhofi.us/x/keybase/types/chat1"
)

func (b *bot) setupMeeting(convid chat1.ConvIDStr, msgid chat1.MessageID, words []string, membersType string) {
	b.debug("command recieved in conversation %s", convid)
	meeting, err := newJitsiMeeting()
	if err != nil {
		log.Println(err)
		b.k.SendMessageByConvID(convid, "I'm sorry, i'm not sure what happened... I was unable to set up a new meeting.\nI've written the appropriate logs and notified my humans.")
		return
	}
	message := fmt.Sprintf("Here's your meeting:\n>URL: %s", meeting.getURL())
	b.k.SendMessageByConvID(convid, message)
}
