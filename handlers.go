package main

import (
	"fmt"
	"log"
	"strings"

	"samhofi.us/x/keybase"
	"samhofi.us/x/keybase/types/chat1"
	"samhofi.us/x/keybase/types/stellar1"
)

// RegisterHandlers is called by main to map these handler funcs to events
func (b *bot) registerHandlers() {
	chat := b.chatHandler
	conv := b.convHandler
	wallet := b.walletHandler
	err := b.errHandler

	b.handlers = keybase.Handlers{
		ChatHandler:         &chat,
		ConversationHandler: &conv,
		WalletHandler:       &wallet,
		ErrorHandler:        &err,
	}
}

// chatHandler should handle all messages coming from the chat
func (b *bot) chatHandler(m chat1.MsgSummary) {
	// only handle text, we don't really care about attachments
	if m.Content.TypeName != "text" {
		return
	}
	// if this chat message is a payment, add it to the bot payments
	if m.Content.Text.Payments != nil {
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
	// if the message is @myusername just perform the default function
	if strings.HasPrefix(m.Content.Text.Body, fmt.Sprintf("@%s", b.k.Username)) {
		words := strings.Fields(m.Content.Text.Body)
		if len(words) > 1 {
			switch words[1] {
			case "meet":
				b.setupMeeting(m.ConvID, m.Sender.Username, words, m.Channel.MembersType)
			case "feedback":
				b.sendFeedback(m.ConvID, m.Id, m.Sender.Username, words)
			case "hello":
				fallthrough
			case "help":
				b.sendWelcome(m.ConvID)
			}
		}
	}
	// its a command for me, iterate through extended commands
	if strings.HasPrefix(m.Content.Text.Body, "!") {
		// break up the message into words
		words := strings.Fields(m.Content.Text.Body)
		// strip the ! from the first word, and lowercase to derive the command
		thisCommand := strings.ToLower(strings.Replace(words[0], "!", "", 1))
		maybeSubCommand := strings.ToLower(words[1])
		// decide if this is askind for extended commands
		switch thisCommand {
		case "jitsi":
			switch maybeSubCommand {
			case "feedback":
				b.sendFeedback(m.ConvID, m.Id, m.Sender.Username, words)
			case "help":
				b.sendWelcome(m.ConvID)
			default:
				b.setupMeeting(m.ConvID, m.Sender.Username, words, m.Channel.MembersType)
			}
		default:
			return
		}
	}
}

// handle conversations (this fires when a new conversation is initiated)
// i.e. when someone opens a conversation to you but hasn't sent a message yet
func (b *bot) convHandler(m chat1.ConvSummary) {
	switch m.Channel.MembersType {
	case "team":
		b.debug("Added to new team: @%s (%s) Sending welcome message", m.Channel.Name, m.Id)
	case "impteamnative":
		b.debug("New conversation found %s (%s) Sending welcome message", m.Channel.Name, m.Id)
	default:
		b.debug("New convID found %s, sending welcome message.", m.Id)
	}
	b.sendWelcome(m.Id)
}

// this handles wallet events, like when someone send you money in chat
func (b *bot) walletHandler(m stellar1.PaymentDetailsLocal) {
	// if the payment is successful
	if m.Summary.StatusSimplified == 3 {
		// get the reply info and see if it exists
		replyInfo := b.payments[m.Summary.Id]
		if replyInfo.convID != "" {
			b.k.ReplyByConvID(replyInfo.convID, replyInfo.msgID, "Thank you so much!  I'll use this to offset my hosting costs!")
		}
	}
}

// this handles all errors returned from the keybase binary
func (b *bot) errHandler(m error) {
	log.Println("---[ error ]---")
	log.Println(p(m))
}
