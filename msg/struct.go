package msg

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/crypto"
	"github.com/h0n9/petit-chat/types"
)

type MsgStructText struct {
	Data     []byte `json:"data"`
	Encoding string `json:"encoding"`
}

func NewMsgStructText(data []byte, encoding string) *MsgStructText {
	return &MsgStructText{Data: data, Encoding: encoding}
}

func (mst *MsgStructText) Encapsulate() ([]byte, error) {
	return json.Marshal(mst)
}

func (mst *MsgStructText) Decapsulate(data []byte) error {
	return json.Unmarshal(data, mst)
}

type MsgStructImage struct {
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

func NewMsgStructImage(data []byte, extension string) *MsgStructImage {
	return &MsgStructImage{Data: data, Extension: extension}
}

func (msi *MsgStructImage) Encapsulate() ([]byte, error) {
	return json.Marshal(msi)
}

func (msi *MsgStructImage) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msi)
}

type MsgStructVideo struct {
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

func NewMsgStructVideo(data []byte, extension string) *MsgStructVideo {
	return &MsgStructVideo{Data: data, Extension: extension}
}

func (msv *MsgStructVideo) Encapsulate() ([]byte, error) {
	return json.Marshal(msv)
}

func (msv *MsgStructVideo) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msv)
}

type MsgStructAudio struct {
	Data      []byte `json:"data"`
	Extension string `json:"extension"`
}

func NewMsgStructAudio(data []byte, extension string) *MsgStructAudio {
	return &MsgStructAudio{Data: data, Extension: extension}
}

func (msa *MsgStructAudio) Encapsulate() ([]byte, error) {
	return json.Marshal(msa)
}

func (msa *MsgStructAudio) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msa)
}

type MsgStructRaw struct {
	Data []byte `json:"data"`
}

func NewMsgStructRaw(data []byte) *MsgStructRaw {
	return &MsgStructRaw{Data: data}
}

func (msr *MsgStructRaw) Encapsulate() ([]byte, error) {
	return json.Marshal(msr)
}

func (msr *MsgStructRaw) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msr)
}

type MsgStructHello struct {
	Persona   *types.Persona    `json:"persona"`
	SecretKey *crypto.SecretKey `json:"secret_key"`
}

func NewMsgStructHello(persona *types.Persona, secretKey *crypto.SecretKey) *MsgStructHello {
	return &MsgStructHello{Persona: persona, SecretKey: secretKey}
}

func (msh *MsgStructHello) Encapsulate() ([]byte, error) {
	return json.Marshal(msh)
}

func (msh *MsgStructHello) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msh)
}

type MsgStructBye struct {
	Persona *types.Persona `json:"persona"`
}

func NewMsgStructBye(persona *types.Persona) *MsgStructBye {
	return &MsgStructBye{Persona: persona}
}

func (msb *MsgStructBye) Encapsulate() ([]byte, error) {
	return json.Marshal(msb)
}

func (msb *MsgStructBye) Decapsulate(data []byte) error {
	return json.Unmarshal(data, msb)
}
