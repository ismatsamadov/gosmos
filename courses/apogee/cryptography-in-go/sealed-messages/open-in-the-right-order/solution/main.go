package main

import "fmt"

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
