package main

import (
	"fmt"
	"net/url"
	"strings"

	"samhofi.us/x/keybase/types/chat1"
)

/*
**** this function is a special case on parameters as it must be called from 2 handlers which
**** get their information from separate types.  As a result we're only passing the conversation id.
**** because of this we can't wrap handleWelcome with permissions, not that you'd want to.
 */

// handleWelcome sends the welcome message to new conversations
func (b *bot) handleWelcome(id chat1.ConvIDStr) {
	b.k.SendMessageByConvID(id, "Hello there!! I'm the Jitsi meeting bot, made by @haukened\nI can start Jitsi meetings right here in this chat!\nI can be activated in 2 ways:\n    1. `@jitsibot`\n    2.`!jitsi`\nYou can provide feedback to my humans using:\n    1. `@jitsibot feedback <type anything>`\n    2. `!jitsibot feedback <type anything>`\nYou can also join @jitsi_meet to talk about features, enhancements, or talk to live humans! Everyone is welcome!\nI also accept donations to offset hosting costs, just send some XLM to my wallet if you feel like it by typing `+5XLM@jitsibot`\nIf you ever need to see this message again, ask me for help or say hello to me!")
}

/*
**** all other commands here-below should only accept a single argument of type chat1.MsgSummary
**** in order to be compliant with the permissions wrapper.  Anything not should be explicitly notated.
 */

// handlePayment controls how the bot reacts to wallet payments in chat
func (b *bot) handlePayment(m chat1.MsgSummary) {
	// there can be multiple payments on each message, iterate them
	for _, payment := range m.Content.Text.Payments {
		if strings.Contains(payment.PaymentText, b.k.Username) {
			// if the payment is successful put log the payment for wallet closure
			if payment.Result.ResultTyp__ == 0 && payment.Result.Error__ == nil {
				var replyInfo = botReply{convID: m.ConvID, msgID: m.Id}
				b.payments[*payment.Result.Sent__] = replyInfo
				b.log("payment recieved %s", payment.PaymentText)
			} else {
				// if the payment fails, be sad
				b.k.ReactByConvID(m.ConvID, m.Id, ":cry:")
			}
		}
	}
}

// handleMeeting starts a new jitsi meeting
func (b *bot) handleMeeting(m chat1.MsgSummary) {
	b.debug("command recieved in conversation %s", m.ConvID)
	// currently we aren't sending dial-in information, so don't get it just generate the name
	// use the simple method
	meeting, err := newJitsiMeetingSimple()
	if err != nil {
		eid := b.logError(err)
		message := fmt.Sprintf("@%s - I'm sorry, i'm not sure what happened... I was unable to set up a new meeting.\nI've written the appropriate logs and notified my humans. Please reference Error ID %s", m.Sender.Username, eid)
		b.k.SendMessageByConvID(m.ConvID, message)
		return
	}
	b.k.SendMessageByConvID(m.ConvID, "@%s here's your meeting: %s", m.Sender.Username, meeting.getURL())
}

// handleFeedback sends feedback to a keybase chat, if configured
func (b *bot) handleFeedback(m chat1.MsgSummary) {
	b.log("feedback recieved in %s", m.ConvID)
	if b.config.FeedbackConvIDStr != "" {
		args := strings.Fields(m.Content.Text.Body)
		feedback := strings.Join(args[2:], " ")
		fcID := chat1.ConvIDStr(b.config.FeedbackConvIDStr)
		if _, err := b.k.SendMessageByConvID(fcID, "Feedback from @%s:\n```%s```", m.Sender.Username, feedback); err != nil {
			eid := b.logError(err)
			b.k.ReactByConvID(m.ConvID, m.Id, "Error ID %s", eid)
		} else {
			b.k.ReplyByConvID(m.ConvID, m.Id, "Thanks! Your feedback has been sent to my human overlords!")
		}
	} else {
		b.k.ReplyByConvID(m.ConvID, m.Id, "I'm sorry, I was unable to send your feedback because my benevolent overlords have not set a destination for feedback. :sob:")
		b.log("user tried to send feedback, but feedback is not enabled. set --feedback-convid or BOT_FEEDBACK_CONVID")
	}
}

// handleSetCommand processes all settings SET calls
func (b *bot) handleSetCommand(m chat1.MsgSummary) {
	b.debug("%s called set command in %s", m.Sender.Username, m.ConvID)
	// first normalize the text and extract the arguments
	args := strings.Fields(strings.ToLower(m.Content.Text.Body))
	if args[1] != "set" {
		return
	}
	switch len(args) {
	case 4:
		if args[2] == "url" {
			// first validate the URL
			u, err := url.ParseRequestURI(args[3])
			if err != nil {
				eid := b.logError(err)
				b.k.ReactByConvID(m.ConvID, m.Id, "Error ID %s", eid)
				return
			}
			// then make sure its HTTPS
			if u.Scheme != "https" {
				b.k.ReactByConvID(m.ConvID, m.Id, "ERROR: HTTPS Required")
				return
			}
			// then get the current options
			var opts ConvOptions
			err = b.KVStoreGetStruct(m.ConvID, &opts)
			if err != nil {
				eid := b.logError(err)
				b.k.ReactByConvID(m.ConvID, m.Id, "Error ID %s", eid)
				return
			}
			// then update the struct using only the scheme and hostname:port
			if u.Port() != "" {
				opts.CustomURL = fmt.Sprintf("%s://%s:%s/", u.Scheme, u.Hostname(), u.Port())
			} else {
				opts.CustomURL = fmt.Sprintf("%s://%s/", u.Scheme, u.Hostname())
			}
			// then write that back to kvstore, with revision
			err = b.KVStorePutStruct(m.ConvID, opts)
			if err != nil {
				eid := b.logError(err)
				b.k.ReactByConvID(m.ConvID, m.Id, "Error ID %s", eid)
				return
			}
			b.k.ReactByConvID(m.ConvID, m.Id, "OK!")
			return
		}
	default:
		return
	}
}

// handleListCommand lists settings for the conversation
func (b *bot) handleListCommand(m chat1.MsgSummary) {
	// first normalize the text and extract the arguments
	args := strings.Fields(strings.ToLower(m.Content.Text.Body))
	if args[0] != "list" {
		return
	}
}
