package model

import (
	"encoding/binary"
	"fmt"
)

type Puzzle struct {
	Nonce      Nonce
	Zeros      ZerosTarget
	PuzzleSize BytesSize
	Puzzle     Hash
}

type Solution struct {
	Nonce    Nonce
	Solution SolutionBytes
}

type (
	Hash          = []byte
	ZerosTarget   = byte
	SolutionBytes = [8]byte
	BytesSize     [2]byte
	Nonce         = [8]byte
)

var BytesOrder = binary.BigEndian

const (
	maxBytesSize uint16 = 64
)

func NewNonce(nonce uint64) Nonce {
	var bb Nonce
	BytesOrder.PutUint64(bb[:], nonce)
	return bb
}

func (s BytesSize) Num() uint16 {
	value := BytesOrder.Uint16(s[:])

	if value > maxBytesSize {
		value = maxBytesSize
	}

	return value
}

func (p Puzzle) String() string {
	return fmt.Sprintf(
		"%x:%x:%x:%x",
		p.Nonce, p.Zeros, p.PuzzleSize, p.Puzzle)
}

func (s Solution) String() string {
	return fmt.Sprintf(
		"%x:%x",
		s.Nonce, s.Solution)
}
