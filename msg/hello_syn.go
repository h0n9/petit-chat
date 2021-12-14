package msg

import (
	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/types"
)

type BodyHelloSyn struct {
	Persona *types.Persona `json:"persona"`
}

func (body *BodyHelloSyn) Check(box *Box, from *From) error {
	// if from.PeerID == box.myID {
	// 	return code.SelfMsg
	// }
	if !box.auth.IsPublic() && !box.auth.CanRead(from.ClientAddr) {
		return code.NonReadPermission
	}
	return nil
}

func (body *BodyHelloSyn) Execute(box *Box, hash types.Hash) error {
	err := box.join(body.Persona)
	if err != nil {
		return err
	}

	// encrypt b.secretKey with msh.Persona.PubKey.GetKey()
	encryptedSecretKey, err := body.Persona.PubKey.Encrypt(box.secretKey.GetKey())
	if err != nil {
		return err
	}

	msg, err := NewMsg(box.myID, box.myPersona.Address, hash, &BodyHelloAck{
		Personae:           box.personae,
		Auth:               box.auth,
		EncryptedSecretKey: encryptedSecretKey,
	})
	if err != nil {
		return err
	}

	err = box.Publish(msg, TypeHelloAck, false)
	if err != nil {
		return err
	}

	return nil
}
