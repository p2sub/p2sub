package llrbmerkle

import (
	"math/big"

	"github.com/p2sub/p2sub/merkle"
	"github.com/petar/GoLLRB/llrb"
)

//Leaf leaf
type Leaf struct {
	Value big.Int
}

//Less is curren item less than?
func (l Leaf) Less(than llrb.Item) (r bool) {
	if b, ok := than.(Leaf); ok && l.Value.CmpAbs(&b.Value) == -1 {
		r = true
	}
	return
}

//Tree sorted and verifiable structure
type Tree struct {
	LlrbTree   *llrb.LLRB
	MerkleTree *merkle.Tree
}

//InsertHash to binary merkle tree
func (t *Tree) InsertHash(digest []byte) {
	//Insert digest to llrb tree
	t.LlrbTree.ReplaceOrInsert(LeafFromDigest(digest))
}

//CalculateMerkle root from sorted binary tree
func (t *Tree) CalculateMerkle() {
	t.MerkleTree = merkle.New()
	digests := make(chan []byte, t.LlrbTree.Len())
	go func(digests chan []byte) {
		t.LlrbTree.AscendGreaterOrEqual(t.LlrbTree.Min(), func(i llrb.Item) bool {
			if l, ok := i.(Leaf); ok {
				digests <- l.Value.Bytes()
			}
			return true
		})
		defer close(digests)
	}(digests)
	for item := range digests {
		t.MerkleTree.AppendHash(item)
	}
	defer t.MerkleTree.Calculate()
}

//LeafFromDigest create leafe from digest
func LeafFromDigest(digest []byte) Leaf {
	return Leaf{Value: *(new(big.Int).SetBytes(digest))}
}
