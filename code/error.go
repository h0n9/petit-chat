package code

import "fmt"

var (
	ImproperPubKey  = fmt.Errorf("improper pubkey")
	ImproperAddress = fmt.Errorf("improper address")

	AlreadySubscribingTopic = fmt.Errorf("already subscribing topic")
	NonSubscribingTopic     = fmt.Errorf("non subscribing topic")
	AlreadyExistingTopic    = fmt.Errorf("already existing topic")
	NonExistingTopic        = fmt.Errorf("non existing topic")
	AlreadyAppendedMsg      = fmt.Errorf("already appended msg")

	ImproperPubSub = fmt.Errorf("improper pubsub")
	ImproperSub    = fmt.Errorf("improper sub")
	ImproperNodeID = fmt.Errorf("improper node id")

	UnknownMsgType         = fmt.Errorf("unknown MsgType")
	ImproperParentMsgHash  = fmt.Errorf("improper ParentMsgHash")
	NonExistingParentMsg   = fmt.Errorf("non existing ParentMsg")
	AlreadyHavingParentMsg = fmt.Errorf("already having ParentMsg")

	ExistingPersonaInBox       = fmt.Errorf("exisiting persona in box")
	NonExistingPersonaInBox    = fmt.Errorf("non exisiting persona in box")
	ImproperPersonaNickname    = fmt.Errorf("improper persona nickname")
	ImproperPersonaMetadata    = fmt.Errorf("improper persona metadata")
	ExistingPersonaNickname    = fmt.Errorf("existing nickname")
	NonExistingPersonaNickname = fmt.Errorf("non existing nickname")
)
