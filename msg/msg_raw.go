package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/code"
)

// TODO: fix size constraint
const (
	MaxDataSize     = 5000
	MaxMetadataSize = 1000
)

type MsgStructRaw struct {
	Data     []byte `json:"data"`
	Metadata []byte `json:"metadata"`
}

func NewMsgStructRaw(data, metadata []byte) (*MsgStructRaw, error) {
	msr := MsgStructRaw{Data: data, Metadata: metadata}
	err := msr.check()
	if err != nil {
		return nil, err
	}
	return &msr, nil
}

func (msr *MsgStructRaw) GetData() []byte {
	return msr.Data
}

func (msr *MsgStructRaw) GetMetadata() []byte {
	return msr.Metadata
}

func (msr *MsgStructRaw) Encapsulate() ([]byte, error) {
	return json.Marshal(msr)
}

func (msr *MsgStructRaw) Decapsulate(data []byte) error {
	err := json.Unmarshal(data, msr)
	if err != nil {
		return err
	}
	err = msr.check()
	if err != nil {
		return err
	}
	return nil
}

func (msr *MsgStructRaw) check() error {
	if len(msr.Data) > MaxDataSize {
		return code.TooBigMsgData
	}
	if len(msr.Metadata) > MaxMetadataSize {
		return code.TooBigMsgMetadata
	}
	return nil
}
