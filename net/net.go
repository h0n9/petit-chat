package net

// Multi Address used on QUIC protocol is formed as follows:
// ex) /ip4/0.0.0.0/udp/61881/quic

const (
	TransportProtocol = "quic"
	ProtocolID        = "/petit-chat/1.0.0"
	RendezVous        = "t'as bien dormi ?"

	DefaultListenAddr = "/ip4/0.0.0.0/udp"
	MinListenPort     = 49152
	MaxListenPort     = 65535
)
