package main

import "fmt"

// Accept checks one SEAL line meant for this satellite and returns the command sealed in it.
func (s *Satellite) Accept(line string) (Command, error) {
	sealed, err := parseSealed(line)
	if err != nil {
		return Command{}, err
	}
	if sealed.Seq <= s.Last {
		return Command{}, fmt.Errorf("%w: %d, last %d", ErrReplay, sealed.Seq, s.Last)
	}
	s.Last = sealed.Seq
	return openCommand(s.Keys[sealed.KeyID].Seal, sealed)
}
