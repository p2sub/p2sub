package main

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"math/big"

	"github.com/petar/GoLLRB/llrb"
)

//Digest size for hash key
const (
	HashKeySize = 32
)

//HashKey digiest of data
type HashKey [HashKeySize]byte

//Merkle node in merkle tree
type MerkleNode struct {
	letf   *MerkleNode
	right  *MerkleNode
	parent *MerkleNode
	digest HashKey
}

func bytesToHashKey(b []byte) (r HashKey) {
	copy(r[:], b[:HashKeySize])
	return r
}

//BinaryMerkleTree sorted and verifiable structure
type BinaryMerkleTree struct {
	tree       *llrb.LLRB
	storage    map[HashKey][]byte
	merkleRoot *MerkleNode
}

//New create new new binary merkle tree
func New() *BinaryMerkleTree {
	return &BinaryMerkleTree{
		tree:    llrb.New(),
		storage: make(map[HashKey][]byte),
	}
}

func (t *BinaryMerkleTree) CalculateMerkleRoot() {
	t.tree.AscendGreaterOrEqual(t.tree.Min(), func(i llrb.Item) bool {
		l, _ := i.(Leaf)
		fmt.Println(l.Value.Text(16))
		return true
	})
}

//Insert data to binary merkle tree
func (t *BinaryMerkleTree) Insert(data []byte) {
	//Calculate digest
	h := sha256.New()
	h.Write(data)
	d := h.Sum(nil)
	k := bytesToHashKey(d)
	fmt.Println(hex.EncodeToString(d))
	//Insert digest to llrb tree
	t.tree.ReplaceOrInsert(LeafFromDigest(k))

	//Insert digest to storage
	t.storage[k] = data
}

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

//LeafFromDigest create leafe from digest
func LeafFromDigest(k HashKey) Leaf {
	n := new(big.Int)
	d := make([]byte, 32)
	copy(d, k[:])
	n.SetBytes(d)
	return Leaf{Value: *n}
}

/*
func newLeaf(val big.Int) Leaf {
	return Leaf{Value: val}
}

func newBigInt(val int) big.Int {
	n := new(big.Int)
	n.SetInt64(int64(val))
	return *n
}*/

func main() {
	/*
		tree := llrb.New()
		tree.ReplaceOrInsert(newLeaf(newBigInt(1)))
		tree.ReplaceOrInsert(newLeaf(newBigInt(2)))
		tree.ReplaceOrInsert(newLeaf(newBigInt(3)))
		tree.ReplaceOrInsert(newLeaf(newBigInt(4)))
		//tree.DeleteMin()
		//tree.Delete(llrb.Inf(1))
		fmt.Println("Item value: ", tree.Min(), tree.Max())
		tree.AscendGreaterOrEqual(tree.Min(), func(i llrb.Item) bool {
			k, _ := i.(Leaf)
			fmt.Println("Item value: ", k.Value.Text(10))
			return true
		})
	*/
	tree := New()
	tree.Insert([]byte("I'm"))
	tree.Insert([]byte("invincible!"))
	tree.Insert([]byte("I am"))
	tree.Insert([]byte("Ironman!"))
	tree.CalculateMerkleRoot()
}
