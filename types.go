package main

import "samhofi.us/x/keybase/v2/types/chat1"

// hold reply information when needed
type botReply struct {
	convID chat1.ConvIDStr
	msgID  chat1.MessageID
}

// ConvOptions stores team specific options like custom servers
type ConvOptions struct {
	ConvID string `json:"converation_id,omitempty"`
	//NotificationsEnabled bool   `json:"notifications_enabled,omitempty"`
	CustomURL string `json:"custom_url,omitempty"`
}

// reflectStruct holds information about reflected structs!
type reflectStruct struct {
	Name  string
	Value interface{}
}
