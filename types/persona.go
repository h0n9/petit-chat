package types

import (
	"encoding/json"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
)

const (
	MaxPersonaNickname = 10  // 10 words, TODO: get from config
	MinPersonaNickname = 1   // 1 words, TODO: get from config
	MaxPersonaMetadata = 800 // 800 words, TODO: get from config
	MinPersonaMetadata = 0   // 0 words, TODO: get from config
)

type Persona struct {
	Nickname string        `json:"nickname"`
	Metadata []byte        `json:"metadata"`
	PubKey   crypto.PubKey `json:"pubkey"`
	Address  crypto.Addr   `json:"address"`
}

type Personae map[crypto.Addr]*Persona

func NewPersona(nickname string, metadata []byte, pubKey crypto.PubKey) (Persona, error) {
	p := Persona{
		Nickname: nickname,
		Metadata: metadata,
		PubKey:   pubKey,
		Address:  pubKey.Address(),
	}
	err := p.Check()
	if err != nil {
		return Persona{}, err
	}
	return p, nil
}

func (p *Persona) GetNickname() string {
	return p.Nickname
}

func (p *Persona) SetNickname(nickname string) error {
	err := checkNickname(nickname)
	if err != nil {
		return err
	}
	p.Nickname = nickname
	return nil
}

func (p *Persona) SetMetadata(md []byte) error {
	err := checkMetadata(md)
	if err != nil {
		return err
	}
	p.Metadata = md
	return nil
}

func (p *Persona) SetPubKey(pk crypto.PubKey) error {
	err := pk.Check()
	if err != nil {
		return err
	}
	p.PubKey = pk
	p.Address = pk.Address()
	return nil
}

func (p *Persona) Encapsulate() ([]byte, error) {
	return json.Marshal(p)
}

func (p *Persona) Decapsulate(data []byte) error {
	return json.Unmarshal(data, p)
}

func (p *Persona) Check() error {
	err := checkNickname(p.Nickname)
	if err != nil {
		return err
	}
	err = checkMetadata(p.Metadata)
	if err != nil {
		return err
	}
	err = checkPubKeyAddress(p.PubKey, p.Address)
	if err != nil {
		return err
	}
	return nil
}

func checkNickname(nn string) error {
	if !checkRange(len(nn), MinPersonaNickname, MaxPersonaNickname) {
		return code.ImproperPersonaNickname
	}
	// TODO: add nickname regex constraint (maybe ?)
	return nil
}

func checkMetadata(md []byte) error {
	if !checkRange(len(md), MinPersonaMetadata, MaxPersonaMetadata) {
		return code.ImproperPersonaMetadata
	}
	return nil
}

func checkPubKeyAddress(pk crypto.PubKey, addr crypto.Addr) error {
	err := pk.Check()
	if err != nil {
		return err
	}
	if !addr.IsDrivenFrom(pk) {
		return code.ImproperAddress
	}
	return nil
}

func checkRange(l, min, max int) bool {
	return (min <= l && l <= max)
}
