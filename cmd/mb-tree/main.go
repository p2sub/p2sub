package main

import (
	"math/big"
)

// MerkleBinaryTree binary hash tree
type MerkleBinaryTree struct {
	hash     *big.Int
	previous *MerkleBinaryTree
	left     *MerkleBinaryTree
	right    *MerkleBinaryTree
}

func newRoot() *MerkleBinaryTree {
	maxHash := hexToInt("ffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffffff")
	rootHash := new(big.Int)
	rootHash.Rsh(maxHash, 1)
	return &MerkleBinaryTree{hash: rootHash,
		previous: nil,
		left:     nil,
		right:    nil}
}

func newLeaf(hexString string, previousLeaf *MerkleBinaryTree) *MerkleBinaryTree {
	return &MerkleBinaryTree{hash: hexToInt(hexString),
		previous: nil,
		left:     nil,
		right:    nil}
}

func hexToInt(hexString string) *big.Int {
	tmp := new(big.Int)
	tmp.SetString(hexString, 16)
	return tmp
}

func isGt(x *big.Int, y *big.Int) bool {
	return x.Cmp(y) == 1
}

func isLt(x *big.Int, y *big.Int) bool {
	return x.Cmp(y) == -1
}

func isEq(x *big.Int, y *big.Int) bool {
	return x.Cmp(y) == 0
}

func (m *MerkleBinaryTree) insert(n *MerkleBinaryTree) {
	inserted := false
	w := m
	for inserted == false {
		if isGt(n.hash, m.hash) {
			if w.right == nil {
				w.right = n
				break
			}
			w = m.right
		} else if isLt(n.hash, m.hash) {
			w = m.left
		} else {
			inserted = true
			break
		}
		if w == nil {
			w = c
			inserted = true
		}
	}
}

func main() {

}
