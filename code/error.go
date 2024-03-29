package code

import "fmt"

var (
	ImproperVault    = fmt.Errorf("improper vault")
	ImproperPubKey   = fmt.Errorf("improper pubkey")
	ImproperAddress  = fmt.Errorf("improper address")
	ImproperSigBytes = fmt.Errorf("improper sigBytes")
	FailedToVerify   = fmt.Errorf("failed to verify")

	AlreadySubscribingTopic = fmt.Errorf("already subscribing topic")
	NonSubscribingTopic     = fmt.Errorf("non subscribing topic")
	AlreadyExistingTopic    = fmt.Errorf("already existing topic")
	NonExistingTopic        = fmt.Errorf("non existing topic")

	NonExistingMsg = fmt.Errorf("non existing msg")
	SelfMsg        = fmt.Errorf("self msg")

	AlreadyAppendedCapsule = fmt.Errorf("already appended capsule")

	ImproperPubSub = fmt.Errorf("improper pubsub")
	ImproperSub    = fmt.Errorf("improper sub")
	ImproperNodeID = fmt.Errorf("improper node id")

	UnknownMsgType      = fmt.Errorf("unknown MsgType")
	ImproperParentHash  = fmt.Errorf("improper ParentHash")
	NonExistingParent   = fmt.Errorf("non existing Parent")
	AlreadyHavingParent = fmt.Errorf("already having Parent")
	TooBigMsgData       = fmt.Errorf("too big msg data")
	TooBigMsgMetadata   = fmt.Errorf("too big msg metadata")

	ExistingPersonaInBox       = fmt.Errorf("exisiting persona in box")
	NonExistingPersonaInBox    = fmt.Errorf("non exisiting persona in box")
	ImproperPersonaNickname    = fmt.Errorf("improper persona nickname")
	ImproperPersonaMetadata    = fmt.Errorf("improper persona metadata")
	ExistingPersonaNickname    = fmt.Errorf("existing nickname")
	NonExistingPersonaNickname = fmt.Errorf("non existing nickname")

	NonExistingPermission     = fmt.Errorf("non existing permission")
	AlreadyExistingPermission = fmt.Errorf("already existing permission")
	NonMinimumPermission      = fmt.Errorf("non minimum permission")
	NonReadPermission         = fmt.Errorf("non read permission")
	NonWritePermission        = fmt.Errorf("non write permission")
	NonExecutePermission      = fmt.Errorf("non execute permission")
)
