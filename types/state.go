package types

import (
	"time"

	"github.com/h0n9/petit-chat/code"
	"github.com/h0n9/petit-chat/crypto"
)

type State struct {
	personae        Personae
	auth            *Auth
	latestTimestamp time.Time
	readUntilIndex  Index
}

func NewState(public bool) *State {
	return &State{
		personae:        make(Personae),
		auth:            NewAuth(public, make(map[crypto.Addr]Perm)),
		latestTimestamp: time.Now(),
		readUntilIndex:  0,
	}
}

func (s *State) GetPersonae() Personae {
	return s.personae
}

func (s *State) SetPersonae(personae Personae) {
	s.personae = personae
}

func (s *State) GetPersona(addr crypto.Addr) *Persona {
	return s.personae[addr]
}

func (s *State) SetPersona(addr crypto.Addr, persona *Persona) {
	s.personae[addr] = persona
}

func (s *State) DeletePersona(addr crypto.Addr) {
	delete(s.personae, addr)
}

func (s *State) GetAuth() *Auth {
	return s.auth
}

func (s *State) SetAuth(auth *Auth) {
	s.auth = auth
}

func (s *State) GetLatestTimestamp() time.Time {
	return s.latestTimestamp
}

func (s *State) SetLatestTimestamp(latestTimestamp time.Time) {
	s.latestTimestamp = latestTimestamp
}

func (s *State) GetReadUntilIndex() Index {
	return s.readUntilIndex
}

func (s *State) SetReadUntilIndex(readUntilIndex Index) {
	s.readUntilIndex = readUntilIndex
}

func (s *State) Join(persona *Persona) error {
	if s.GetPersona(persona.Address) != nil {
		return nil // ignore even if existing
	}
	err := persona.Check()
	if err != nil {
		return err
	}
	s.SetPersona(persona.Address, persona)
	return nil
}

func (s *State) Leave(persona *Persona) error {
	if s.GetPersona(persona.Address) == nil {
		return code.NonExistingPersonaInBox
	}
	s.DeletePersona(persona.Address)
	return nil
}

func (s *State) Grant(persona *Persona, r, w, x bool) error {
	return s.auth.Grant(persona, r, w, x)
}
