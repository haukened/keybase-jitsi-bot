package main

import (
	"fmt"
	"log"
	"net/url"
	"strings"

	"samhofi.us/x/keybase/types/chat1"
)

// handlePayment controls how the bot reacts to wallet payments in chat
func (b *bot) handlePayment(m chat1.MsgSummary) {
	// there can be multiple payments on each message, iterate them
	for _, payment := range m.Content.Text.Payments {
		if strings.Contains(payment.PaymentText, b.k.Username) {
			// if the payment is successful put log the payment for wallet closure
			if payment.Result.ResultTyp__ == 0 && payment.Result.Error__ == nil {
				var replyInfo = botReply{convID: m.ConvID, msgID: m.Id}
				b.payments[*payment.Result.Sent__] = replyInfo
			} else {
				// if the payment fails, be sad
				b.k.ReactByConvID(m.ConvID, m.Id, ":cry:")
			}
		}
	}
}

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

func (b *bot) setKValue(convid chat1.ConvIDStr, msgID chat1.MessageID, args []string) {
	if args[0] != "set" {
		return
	}
	switch len(args) {
	case 3:
		if args[1] == "url" {
			// first validate the URL
			u, err := url.ParseRequestURI(args[2])
			if err != nil {
				b.k.ReplyByConvID(convid, msgID, "ERROR - `%s`", err)
				return
			}
			// then make sure its HTTPS
			if u.Scheme != "https" {
				b.k.ReplyByConvID(convid, msgID, "ERROR - HTTPS Required")
				return
			}
			// then get the current options
			var opts ConvOptions
			err = b.KVStoreGetStruct(convid, &opts)
			if err != nil {
				eid := b.logError(err)
				b.k.ReactByConvID(convid, msgID, "Error %s", eid)
				return
			}
			// then update the struct using only the scheme and hostname:port
			if u.Port() != "" {
				opts.CustomURL = fmt.Sprintf("%s://%s:%s/", u.Scheme, u.Hostname(), u.Port())
			} else {
				opts.CustomURL = fmt.Sprintf("%s://%s/", u.Scheme, u.Hostname())
			}
			// then write that back to kvstore, with revision
			err = b.KVStorePutStruct(convid, opts)
			if err != nil {
				eid := b.logError(err)
				b.k.ReactByConvID(convid, msgID, "ERROR %s", eid)
				return
			}
			b.k.ReactByConvID(convid, msgID, "OK!")
			return
		}
	default:
		return
	}
}

func (b *bot) listKValue(convid chat1.ConvIDStr, msgID chat1.MessageID, args []string) {
	if args[0] != "list" {
		return
	}
}
