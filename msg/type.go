package msg

import "github.com/h0n9/petit-chat/code"

type Type uint32

const (
	TypeNone Type = iota + 1 // Msg 0 means something wrong
	TypeHelloSyn
	TypeHelloAck
	TypeMeta
	TypeUpdate
	TypeRaw
	TypeBye // End of Subscription
)

var TypeStrMap = map[Type]string{
	TypeNone:     "TypeNone",
	TypeHelloSyn: "TypeHelloSyn",
	TypeHelloAck: "TypeHelloAck",
	TypeMeta:     "TypeMeta",
	TypeUpdate:   "TypeUpdate",
	TypeRaw:      "TypeRaw",
	TypeBye:      "TypeBye",
}

func (t Type) Base() Base {
	switch t {
	case TypeHelloSyn:
		return &HelloSyn{}
	case TypeHelloAck:
		return &HelloAck{}
	case TypeMeta:
		return &Meta{}
	case TypeUpdate:
		return &Update{}
	case TypeRaw:
		return &Raw{}
	case TypeBye:
		return &Bye{}
	default:
		return nil
	}
}

func (t Type) String() string {
	err := t.Check()
	if err != nil {
		return ""
	}
	return TypeStrMap[t]
}

func (t Type) Check() error {
	_, exist := TypeStrMap[t]
	if !exist {
		return code.UnknownMsgType
	}
	return nil
}
