package llrbmerkle

import (
	"github.com/petar/GoLLRB/llrb"
)

//New create new new binary merkle tree
func New() *Tree {
	return &Tree{
		LlrbTree: llrb.New(),
	}
}
