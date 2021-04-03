package code

import "fmt"

var (
	ImproperPubKey = fmt.Errorf("improper pubkey")

	AlreadySubscribingTopic = fmt.Errorf("already subscribing topic")
	NonSubscribingTopic     = fmt.Errorf("non subscribing topic")
	AlreadyExistingTopic    = fmt.Errorf("already existing topic")
	NonExistingTopic        = fmt.Errorf("non existing topic")
	AlreadyAppendedMsg      = fmt.Errorf("already appended msg")

	AlreadyExistingNickname = fmt.Errorf("already existing nickname")
	NonExistingNickname     = fmt.Errorf("non existing nickname")

	ImproperPubSub = fmt.Errorf("improper pubsub")
	ImproperSub    = fmt.Errorf("improper sub")
	ImproperNodeID = fmt.Errorf("improper node id")

	UnknownMsgType         = fmt.Errorf("unknown MsgType")
	ImproperParentMsgHash  = fmt.Errorf("improper ParentMsgHash")
	NonExistingParentMsg   = fmt.Errorf("non existing ParentMsg")
	AlreadyHavingParentMsg = fmt.Errorf("already having ParentMsg")

	ImproperPersonaNickname = fmt.Errorf("improper persona nickname")
	ImproperPersonaMetadata = fmt.Errorf("improper persona metadata")
)
