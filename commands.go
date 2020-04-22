package main

import (
	"bytes"
	"fmt"
	"net/url"
	"strings"
	"text/tabwriter"

	"samhofi.us/x/keybase/v2/types/chat1"
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
	// check and see if this conversation has a custom URL
	opts := ConvOptions{}
	err := b.KVStoreGetStruct(m.ConvID, &opts)
	if err != nil {
		b.debug("unable to get conversation options")
		eid := b.logError(err)
		b.k.ReactByConvID(m.ConvID, m.Id, "Error ID %s", eid)
		return
	}
	// currently we aren't sending dial-in information, so don't get it just generate the name
	// use the simple method
	meeting, err := newJitsiMeetingSimple()
	if err != nil {
		eid := b.logError(err)
		message := fmt.Sprintf("@%s - I'm sorry, i'm not sure what happened... I was unable to set up a new meeting.\nI've written the appropriate logs and notified my humans. Please reference Error ID %s", m.Sender.Username, eid)
		b.k.SendMessageByConvID(m.ConvID, message)
		return
	}
	// then set the Custom server URL, if it exists
	if opts.ConvID == string(m.ConvID) && opts.CustomURL != "" {
		meeting.CustomServer = opts.CustomURL
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

// handleConfigCommand dispatches config calls
func (b *bot) handleConfigCommand(m chat1.MsgSummary) {
	args := strings.Fields(strings.ToLower(m.Content.Text.Body))
	if args[1] != "config" {
		return
	}
	if len(args) >= 3 {
		switch args[2] {
		case "set":
			b.handleConfigSet(m)
			return
		case "list":
			b.handleConfigList(m)
			return
		case "help":
			b.handleConfigHelp(m)
		}
	}
}

// handleConfigSet processes all settings SET calls
// this should be called from b.handleConfigCommand()
func (b *bot) handleConfigSet(m chat1.MsgSummary) {
	// first normalize the text and extract the arguments
	args := strings.Fields(strings.ToLower(m.Content.Text.Body))
	if args[2] != "set" {
		return
	}
	b.debug("config set called by @%s in %s", m.Sender.Username, m.ConvID)
	switch len(args) {
	case 5:
		if args[3] == "customurl" {
			// first validate the URL
			u, err := url.ParseRequestURI(args[4])
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
				opts.CustomURL = fmt.Sprintf("%s://%s:%s", u.Scheme, u.Hostname(), u.Port())
			} else {
				opts.CustomURL = fmt.Sprintf("%s://%s", u.Scheme, u.Hostname())
			}
			// ensure that the struct has convid filled out (if its new it won't)
			if opts.ConvID == "" {
				opts.ConvID = string(m.ConvID)
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

// handleConfigList lists settings for the conversation
// this should be called from b.handleConfigCommand()
func (b *bot) handleConfigList(m chat1.MsgSummary) {
	// first normalize the text and extract the arguments
	args := strings.Fields(strings.ToLower(m.Content.Text.Body))
	if args[2] != "list" {
		return
	}
	// get the ConvOptions
	var opts ConvOptions
	err := b.KVStoreGetStruct(m.ConvID, &opts)
	if err != nil {
		eid := b.logError(err)
		b.k.ReactByConvID(m.ConvID, m.Id, "Error ID %s", eid)
		return
	}
	// then reflect the struct to a list
	configOpts := structToSlice(opts)
	// Then iterate those through a tabWriter
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 0, 3, ' ', 0)
	fmt.Fprintln(w, "Config Options for this channel:\n```")
	for _, opt := range configOpts {
		if opt.Name == "ConvID" {
			// dont print out the conversation id, thats not an option.
			continue
		}
		fmt.Fprintf(w, "%s\t%v\t\n", opt.Name, opt.Value)
	}
	fmt.Fprintln(w, "```")
	w.Flush()
	b.k.ReplyByConvID(m.ConvID, m.Id, buf.String())
}

// handleConfigHelp shows config help
// this should be called from b.handleConfigCommand()
func (b *bot) handleConfigHelp(m chat1.MsgSummary) {
	// first normalize the text and extract the arguments
	args := strings.Fields(strings.ToLower(m.Content.Text.Body))
	if args[2] != "help" {
		return
	}
	b.debug("config help called by @%s in %s", m.Sender.Username, m.ConvID)
}
