package types

type Meta uint8

const (
	metaNone     Meta = 0 // 0000 0000
	metaReceived Meta = 1 // 0000 0001
	metaRead     Meta = 2 // 0000 0010
	metaTyping   Meta = 4 // 0000 0100
)

func NewMeta(received, read, typing bool) Meta {
	var meta Meta
	if received {
		meta |= metaReceived
	}
	if read {
		meta |= metaRead
	}
	if typing {
		meta |= metaTyping
	}
	return meta
}

func (meta Meta) Received() bool {
	return meta&metaReceived == metaReceived
}

func (meta Meta) Read() bool {
	return meta&metaRead == metaRead
}

func (meta Meta) Typing() bool {
	return meta&metaTyping == metaTyping
}
