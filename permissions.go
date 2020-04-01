package main

import (
	"strings"

	"samhofi.us/x/keybase/types/chat1"
)

// checkPermissionAndExecute will check the minimum required role for the permission and execute the handler function if allowed
func (b *bot) checkPermissionAndExecute(requiredRole string, m chat1.MsgSummary, f func(chat1.MsgSummary)) {
	// get the members of the conversation
	b.debug("Executing permissions check")
	// currently this doesn't work due to a keybase bug
	// the workaround is to check the general channel the old way
	//conversation, err := b.k.ListMembersOfConversation(m.ConvID)

	// **** <workaround>
	channel := chat1.ChatChannel{
		Name:        m.Channel.Name,
		MembersType: m.Channel.MembersType,
		TopicName:   "general",
	}
	conversation, err := b.k.ListMembersOfChannel(channel)
	/// **** </workaround>

	if err != nil {
		eid := b.logError(err)
		b.k.ReactByConvID(m.ConvID, m.Id, "Error ID %s", eid)
		return
	}
	// create a map of valid roles, according to @dxb struc
	memberTypes := make(map[string]struct{})
	memberTypes["owner"] = struct{}{}
	memberTypes["admin"] = struct{}{}
	memberTypes["writer"] = struct{}{}
	memberTypes["reader"] = struct{}{}

	// if the role is not in the map, its an invalid role
	if _, ok := memberTypes[strings.ToLower(requiredRole)]; !ok {
		// the role passed was not valid, so bail
		b.log("ERROR: %s is not a valid permissions level", requiredRole)
		return
	}

	// then descend permissions from top down
	for _, member := range conversation.Members.Owners {
		if strings.ToLower(member.Username) == strings.ToLower(m.Sender.Username) {
			f(m)
			return
		}
		b.debug("no")
	}
	// if the required role was owner, return and don't evaluate the rest
	if strings.ToLower(requiredRole) == "owner" {
		b.debug("user does not have required permission of: owner")
		return
	}
	// admins
	for _, member := range conversation.Members.Admins {
		if strings.ToLower(member.Username) == strings.ToLower(m.Sender.Username) {
			f(m)
			return
		}
	}
	if strings.ToLower(requiredRole) == "admin" {
		b.debug("user does not have required permission of: admin")
		return
	}
	// writers
	for _, member := range conversation.Members.Writers {
		if strings.ToLower(member.Username) == strings.ToLower(m.Sender.Username) {
			f(m)
			return
		}
	}
	if strings.ToLower(requiredRole) == "writer" {
		b.debug("user does not have required permission of: writer")
		return
	}
	// readers
	for _, member := range conversation.Members.Readers {
		if strings.ToLower(member.Username) == strings.ToLower(m.Sender.Username) {
			f(m)
			return
		}
	}
	// just return - restricted bots shouldn't be able to run commands
	b.debug("user does not have required permission of: reader")
	return
}
