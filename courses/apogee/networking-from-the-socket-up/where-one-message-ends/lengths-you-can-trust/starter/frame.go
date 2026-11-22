package main

import "errors"

var (
	ErrBadLength = errors.New("bad length") // wrapped: "bad length 65535"
	ErrBadCRC    = errors.New("bad crc")
)
