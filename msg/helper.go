package msg

import "github.com/h0n9/petit-chat/types"

type Helper interface {
	// accessors
	GetVault() *types.Vault
	GetState() *types.State
	GetStore() *CapsuleStore
	GetPeerID() types.ID

	// operators
	Publish(msg *Msg, encrypt bool) error
}
