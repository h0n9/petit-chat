package msg

import "github.com/h0n9/petit-chat/code"

type Type uint32

const (
	TypeNone Type = iota + 1 // Msg 0 means something wrong
	TypeRaw
	TypeHelloSyn
	TypeHelloAck
	TypeBye // End of Subscription
	TypeUpdateSyn
	TypeUpdateAck
)

var TypeMap = map[Type]string{
	TypeNone:      "TypeNone",
	TypeRaw:       "TypeRaw",
	TypeHelloSyn:  "TypeHelloSyn",
	TypeHelloAck:  "TypeHelloAck",
	TypeBye:       "TypeBye",
	TypeUpdateSyn: "TypeUpdateSyn",
	TypeUpdateAck: "TypeUpdateAck",
}

func (t Type) String() string {
	err := t.Check()
	if err != nil {
		return ""
	}
	return TypeMap[t]
}

func (t Type) Check() error {
	_, exist := TypeMap[t]
	if !exist {
		return code.UnknownMsgType
	}
	return nil
}
