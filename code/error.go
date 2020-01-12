package code

import "fmt"

var (
	AlreadySubscribingTopic = fmt.Errorf("already subscribing topic")
	NonSubscribingTopic     = fmt.Errorf("non subscribing topic")
	AlreadyExistingTopic    = fmt.Errorf("already existing topic")
	NonExistingTopic        = fmt.Errorf("non existing topic")
	AlreadyAppendedMsg      = fmt.Errorf("already appended msg")

	ImproperPubSub = fmt.Errorf("improper pubsub")
	ImproperSub    = fmt.Errorf("improper sub")
	ImproperNodeID = fmt.Errorf("improper node id")
)
