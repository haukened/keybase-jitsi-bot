package main

import (
	"bytes"
	"fmt"
	"log"
	"text/tabwriter"

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
	var buf bytes.Buffer
	w := tabwriter.NewWriter(&buf, 0, 4, 3, ' ', 0)
	fmt.Fprintln(w, "Here's your meeting:")
	fmt.Fprintf(w, "URL:\t%s\n", meeting.getURL())
	fmt.Fprintf(w, "PIN:\t%s\n", meeting.getPIN())
	fmt.Fprintln(w, "Dial In:\t")
	fmt.Fprintln(w, "```")
	for _, phone := range meeting.Phone {
		fmt.Fprintf(w, "    %s\t%s\t\n", phone.Country, phone.Number)
	}
	fmt.Fprintln(w, "```")
	w.Flush()
	b.k.SendMessageByConvID(convid, buf.String())
}
