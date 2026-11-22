package main

import (
	"fmt"
	"time"
)

// SatKeys are the keys a satellite holds for one master version.
type SatKeys struct {
	Sign, Seal []byte
	Retire     time.Time // refused from this instant; zero: not scheduled
	Revoked    bool
}

// Satellite stands in for a satellite's flight software: it opens the SEAL lines meant for it.
type Satellite struct {
	ID   string
	Keys map[string]SatKeys // by key ID
	Last uint64             // the last sequence number it accepted
	Now  func() time.Time
}

// keyFor returns the keys of master version id, or why the satellite won't use them now.
func (s *Satellite) keyFor(id string) (SatKeys, error) {
	k, ok := s.Keys[id]
	switch {
	case !ok:
		return SatKeys{}, fmt.Errorf("%w: %q", ErrUnknownKey, id)
	case k.Revoked:
		return SatKeys{}, fmt.Errorf("%w: %q", ErrRevoked, id)
	case !k.Retire.IsZero() && !s.Now().Before(k.Retire):
		return SatKeys{}, fmt.Errorf("%w: %q", ErrRetired, id)
	}
	return k, nil
}

// Accept checks one SEAL line meant for this satellite and returns the command sealed in it.
func (s *Satellite) Accept(line string) (Command, error) {
	sealed, err := parseSealed(line)
	if err != nil {
		return Command{}, err
	}
	if sealed.SatID != s.ID {
		return Command{}, fmt.Errorf("%w: %q", ErrWrongSat, sealed.SatID)
	}
	keys, err := s.keyFor(sealed.KeyID)
	if err != nil {
		return Command{}, err
	}
	if sealed.Seq <= s.Last {
		return Command{}, fmt.Errorf("%w: %d, last %d", ErrReplay, sealed.Seq, s.Last)
	}
	c, err := openCommand(keys.Seal, sealed)
	if err != nil {
		return Command{}, err
	}
	s.Last = sealed.Seq
	return c, nil
}

// Receive is what the flight software does with a SEAL line since it learned LOAD_KEY: it checks the
// line as Accept does, then installs the keys a LOAD_KEY carries. A refused line changes nothing.
func (s *Satellite) Receive(line string) (Command, error) {
	last := s.Last
	c, err := s.Accept(line)
	if err != nil || c.Op != "LOAD_KEY" {
		return c, err
	}
	if err := s.load(c.Args); err != nil {
		s.Last = last
		return Command{}, err
	}
	return c, nil
}

// load installs the key version a LOAD_KEY carries: its ID, then its sign and seal keys in hex. A
// version the satellite holds as revoked stays revoked.
func (s *Satellite) load(args []string) error {
	if len(args) != 3 {
		return fmt.Errorf("%w: LOAD_KEY takes 3 arguments", ErrMalformed)
	}
	if held, ok := s.Keys[args[0]]; ok && held.Revoked {
		return fmt.Errorf("%w: %q", ErrRevoked, args[0])
	}
	sign, signOK := hexBytes(args[1], 32, 32)
	seal, sealOK := hexBytes(args[2], 32, 32)
	if !signOK || !sealOK {
		return fmt.Errorf("%w: LOAD_KEY keys are 64 hex digits each", ErrMalformed)
	}
	s.Keys[args[0]] = SatKeys{Sign: sign, Seal: seal}
	return nil
}
