package main

// ConvOptions stores team specific options like custom servers
type ConvOptions struct {
	ConvID               string `json:"converation_id,omitempty"`
	NotificationsEnabled bool   `json:"notifications_enabled,omitempty"`
	CustomURL            string `json:"custom_url,omitempty"`
}
