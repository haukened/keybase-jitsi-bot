package main

import (
	"log"
	"strings"

	"samhofi.us/x/keybase/v2"
	"samhofi.us/x/keybase/v2/types/chat1"
	"samhofi.us/x/keybase/v2/types/stellar1"
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
		b.handlePayment(m)
		return
	}

	// if its not a payment evaluate if this is a command at all
	if strings.HasPrefix(m.Content.Text.Body, "!") || strings.HasPrefix(m.Content.Text.Body, "@") {
		// first return if its not a command for me
		if !strings.Contains(m.Content.Text.Body, b.cmd()) && !strings.Contains(m.Content.Text.Body, b.k.Username) {
			return
		}
		// then check if this is the root command
		if isRootCommand(m.Content.Text.Body, b.cmd(), b.k.Username) {
			b.checkPermissionAndExecute("reader", m, b.handleMeeting)
			return
		}

		// then check help and welcome (non-permissions)
		// help
		if hasCommandPrefix(m.Content.Text.Body, b.cmd(), b.k.Username, "help") {
			b.handleWelcome(m.ConvID)
			return
		}
		// hello
		if hasCommandPrefix(m.Content.Text.Body, b.cmd(), b.k.Username, "hello") {
			b.handleWelcome(m.ConvID)
			return
		}

		// then check sub-command variants (permissions)
		// meet
		if hasCommandPrefix(m.Content.Text.Body, b.cmd(), b.k.Username, "meet") {
			b.checkPermissionAndExecute("reader", m, b.handleMeeting)
			return
		}
		// feedback
		if hasCommandPrefix(m.Content.Text.Body, b.cmd(), b.k.Username, "feedback") {
			b.checkPermissionAndExecute("reader", m, b.handleFeedback)
			return
		}
		// config commands
		if hasCommandPrefix(m.Content.Text.Body, b.cmd(), b.k.Username, "config") {
			b.checkPermissionAndExecute("admin", m, b.handleConfigCommand)
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
	b.handleWelcome(m.Id)
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
