package client

//
// import (
// 	"github.com/h0n9/petit-chat/code"
// 	"github.com/h0n9/petit-chat/msg"
// )
//
// type MsgHandler func(box *msg.Box, msg *msg.Msg) (bool, error)
//
// func DefaultMsgHandler(box *msg.Box, msg *msg.Msg) (bool, error) {
// 	eos := msg.IsEOS() && (msg.GetPeerID() == box.vault.GetID())
//
// 	// msg handling flow:
// 	//   check -> append -> execute -> (received)
//
// 	// check if msg is proper and can be supported on protocol
// 	// improper msgs are dropped here
// 	err := msg.check(box)
// 	if err != nil {
// 		return eos, err
// 	}
//
// 	// check msg.Body
// 	err = msg.Check(box)
// 	if err != nil && err != code.SelfMsg {
// 		return eos, err
// 	}
//
// 	// execute msg.Body
// 	err = msg.Execute(box)
// 	if err != nil {
// 		return eos, err
// 	}
//
// 	// append msg
// 	// readUntilIndex, err := box.append(msg)
// 	// if err != nil {
// 	// 	return eos, err
// 	// }
//
// 	// canRead := box.msgSubCh != nil
// 	// if canRead {
// 	// 	box.msgSubCh <- msg
// 	// 	box.state.SetReadUntilIndex(readUntilIndex)
// 	// }
// 	// if msg.GetType() <= TypeMeta {
// 	// 	return eos, nil
// 	// }
// 	// if msg.GetClientAddr() == box.vault.GetAddr() {
// 	// 	return eos, nil
// 	// }
// 	// meta := types.NewMeta(true, canRead, false)
// 	// msgMeta := NewMsgMeta(box, types.EmptyHash, msg.GetHash(), meta)
// 	// err = box.Publish(msgMeta, true)
// 	// if err != nil {
// 	// 	return eos, err
// 	// }
//
// 	return eos, nil
// }
//
