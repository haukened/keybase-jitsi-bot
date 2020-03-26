package main

import (
	"fmt"
	"log"
	"strings"

	"samhofi.us/x/keybase/types/chat1"
)

func (b *bot) setupMeeting(convid chat1.ConvIDStr, sender string, args []string, membersType string) {
	b.debug("command recieved in conversation %s", convid)
	meeting, err := newJitsiMeetingSimple()
	if err != nil {
		log.Println(err)
		message := fmt.Sprintf("@%s - I'm sorry, i'm not sure what happened... I was unable to set up a new meeting.\nI've written the appropriate logs and notified my humans.", sender)
		b.k.SendMessageByConvID(convid, message)
		return
	}
	message := fmt.Sprintf("@%s here's your meeting: %s", sender, meeting.getURL())
	b.k.SendMessageByConvID(convid, message)
}

func (b *bot) sendFeedback(convid chat1.ConvIDStr, mesgID chat1.MessageID, sender string, args []string) {
	b.debug("feedback recieved in %s", convid)
	if b.config.FeedbackConvIDStr != "" {
		feedback := strings.Join(args, " ")
		fcID := chat1.ConvIDStr(b.config.FeedbackConvIDStr)
		if _, err := b.k.SendMessageByConvID(fcID, "Feedback from @%s:\n```%s```", sender, feedback); err != nil {
			b.k.ReplyByConvID(convid, mesgID, "I'm sorry, I was unable to send your feedback because my benevolent overlords have not set a destination for feedback. :sad:")
			log.Printf("Unable to send feedback: %s", err)
		} else {
			b.k.ReplyByConvID(convid, mesgID, "Thanks! Your feedback has been sent to my human overlords!")
		}
	} else {
		b.debug("feedback not enabled. set --feedback-convid or BOT_FEEDBACK_CONVID")
	}
}

func (b *bot) sendWelcome(convid chat1.ConvIDStr) {
	b.k.SendMessageByConvID(convid, "Hello there!! I'm the Jitsi meeting bot, made by @haukened\nI can start Jitsi meetings right here in this chat!\nI can be activated in 2 ways:\n    1. `@jitsibot`\n    2.`!jitsi`\nYou can provide feedback to my humans using:\n    1. `@jitsibot feedback <type anything>`\n    2. `!jitsibot feedback <type anything>`\nYou can also join @jitsi_meet to talk about features, enhancements, or talk to live humans! Everyone is welcome!\nI also accept donations to offset hosting costs, just send some XLM to my wallet if you feel like it by typing `+5XLM@jitsibot`\nIf you ever need to see this message again, ask me for help or say hello to me!")
}
