package merkle

import (
	"encoding/hex"
	"fmt"

	"github.com/p2sub/p2sub/utilities"
)

//Node node in merkle tree
type Node struct {
	Left   *Node
	Right  *Node
	Digest []byte
}

//Tree in memory merkle tree
type Tree struct {
	Nodes [][]Node
	Root  *Node
}

//Append leaf to tree
func (t *Tree) Append(node Node) {
	t.updateOrInsertNode(0, -1, node)
}

//AppendHash to merkle tree
func (t *Tree) AppendHash(hash []byte) {
	t.updateOrInsertNode(0, -1, Node{Digest: hash})
}

//AppendData to merkle tree
func (t *Tree) AppendData(dat []byte) {
	t.updateOrInsertNode(0, -1, Node{Digest: utilities.FastSha256(dat)})
}

//Calculate merkle root
func (t *Tree) Calculate(level ...int) {
	//Get current level
	c := 0
	if len(level) == 1 {
		c = level[0]
	}
	//Get total element of curren level
	e := len(t.Nodes[c])
	//Found the root
	if e == 1 {
		t.Root = &t.Nodes[c][0]
		return
	}
	for i := 0; i < e; i += 2 {
		if i+1 >= e {
			//Duplicate node
			t.updateOrInsertNode(c+1, i/2, Node{
				Right:  &t.Nodes[c][i],
				Left:   nil,
				Digest: t.Nodes[c][i].Digest,
			})
		} else {
			//Calculate next level node
			t.updateOrInsertNode(c+1, i/2, Node{
				Right:  &t.Nodes[c][i],
				Left:   &t.Nodes[c][i+1],
				Digest: utilities.FastSha256(t.Nodes[c][i].Digest, t.Nodes[c][i+1].Digest),
			})
		}
	}
	t.Calculate(c + 1)
}

//PrintTree for fun @todo remove later
func (t *Tree) PrintTree() {
	for l := 0; l < len(t.Nodes); l++ {
		for i := 0; i < len(t.Nodes[l]); i++ {
			fmt.Printf("%s\t", hex.EncodeToString(t.Nodes[l][i].Digest)[:16])
		}
		fmt.Printf("\n")
	}
}

func (t *Tree) updateOrInsertNode(level int, index int, node Node) {
	if level >= len(t.Nodes) {
		t.Nodes = append(t.Nodes, make([]Node, 0))
	}
	if index >= len(t.Nodes[level]) || index < 0 {
		t.Nodes[level] = append(t.Nodes[level], node)
	} else {
		t.Nodes[level][index] = node
	}
}
