package main

import (
	"fmt"
	"log"
	"strings"

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

func (b *bot) sendFeedback(convid chat1.ConvIDStr, mesgID chat1.MessageID, sender string, words []string) {
	b.debug("feedback recieved in %s", convid)
	if b.config.FeedbackConvIDStr != "" {
		feedback := strings.Join(words[2:], " ")
		fcID := chat1.ConvIDStr(b.config.FeedbackConvIDStr)
		if _, err := b.k.SendMessageByConvID(fcID, "Feedback from @%s:\n```%s```", sender, feedback); err != nil {
			b.k.ReplyByConvID(convid, mesgID, "I'm sorry, I was unable to send your feedback because my benevolent overlords have not set a destination for feedback. :sad:")
			log.Printf("Unable to send feedback: %s", err)
		} else {
			b.k.ReplyByConvID(convid, mesgID, "Thanks! Your feedback has been sent to my human overlords!")
			b.debug("feedback sent")
		}
	} else {
		b.debug("feedback not enabled. set --feedback-convid or BOT_FEEDBACK_CONVID")
	}
}
