package msg

import "github.com/h0n9/petit-chat/code"

type Type uint32

const (
	TypeNone Type = iota + 1 // Msg 0 means something wrong
	TypeRaw
	TypeHelloSyn
	TypeHelloAck
	TypeBye // End of Subscription
	TypeUpdate
	TypeMeta
)

var TypeStrMap = map[Type]string{
	TypeNone:     "TypeNone",
	TypeRaw:      "TypeRaw",
	TypeHelloSyn: "TypeHelloSyn",
	TypeHelloAck: "TypeHelloAck",
	TypeBye:      "TypeBye",
	TypeUpdate:   "TypeUpdate",
	TypeMeta:     "TypeMeta",
}

func (t Type) Base() Base {
	switch t {
	case TypeRaw:
		return &Raw{}
	case TypeHelloSyn:
		return &HelloSyn{}
	case TypeHelloAck:
		return &HelloAck{}
	case TypeBye:
		return &Bye{}
	case TypeUpdate:
		return &Update{}
	case TypeMeta:
		return &Meta{}
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
