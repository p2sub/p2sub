package merkle

//New create a new merkle tree
func New() *Tree {
	return &Tree{Nodes: make([][]Node, 0), Root: nil}
}
