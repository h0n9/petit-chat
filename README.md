# petit-chat

The program enables peer-to-peer chat which doesn't require a centralized
server relaying messages from peer to peer. It is implemented based on
[libp2p](https://github.com/libp2p/go-libp2p) for network layers and written
in [Go](https://golang.org). Basic ideas start from "Nobody wants his/her
private messages get relayed and stored by untrustful third-party service
providers even though they advertise the services' concrete security protocol".
In other words, any kind of data transmitted between clients remain on their
side, not on the side of middle men who could take off masks and reveal real
face of 'Big Brother'.

