package main

import (
	"crypto/hkdf"
	"crypto/sha256"
)

// fleetKeys derives every listed satellite's sign and seal keys from master, by info such as "seal|cirrine-1".
func fleetKeys(master Key, satIDs []string) (map[string][]byte, error) {
	keys := make(map[string][]byte, 2*len(satIDs))
	prk, err := hkdf.Extract(sha256.New, master.Secret, salt)
	if err != nil {
		return nil, err
	}
	for _, sat := range satIDs {
		for _, purpose := range []string{"sign", "seal"} {
			info := purpose + "|" + sat
			key, err := hkdf.Expand(sha256.New, prk, info, 32)
			if err != nil {
				return nil, err
			}
			keys[info] = key
		}
	}
	return keys, nil
}
