package tree

import (
	"fmt"
	"math/rand"
	"strings"
)

type Nodes []*Node

func (n Nodes) String() string {
	var names []string
	for _, node := range n {
		names = append(names, node.Name)
	}
	return strings.Join(names, ", ")
}

type Node struct {
	Name string

	parent   *Node
	children []*Node
}

func NewNode(t Treer) Treer {
	if t == nil {
		return nil
	}

	v, ok := t.(*Node)
	if !ok {
		return nil
	}

	return &Node{
		Name: v.Name,
	}
}

func (n *Node) ID() string {
	if n == nil {
		return ""
	}
	return n.Name
}

func (n *Node) Parent() Treer {
	if n == nil {
		return nil
	}

	if n.parent == nil {
		return nil
	}

	return n.parent
}

func (n *Node) Children() []Treer {
	if n == nil {
		return nil
	}
	children := make([]Treer, len(n.children))
	for idx, child := range n.children {
		children[idx] = child
	}
	return children
}

func (n *Node) SetParent(parent Treer) {
	if n == nil {
		return
	}

	if parent == nil {
		n.parent = nil
		return
	}

	n.parent = parent.(*Node)
}

func (n *Node) SetChildren(children []Treer) {
	if n == nil {
		return
	}
	n.children = nil
	for _, child := range children {
		n.children = append(n.children, child.(*Node))
	}
}

func (n *Node) String() string {
	if n == nil {
		return "<nil>"
	}
	return n.Name
}

func RandomTree(depth int, maxChildren int) Treer {
	name := newName("Node")

	if depth <= 0 {
		return nil
	}

	root := &Node{Name: name()}
	if depth == 1 {
		return root
	}

	var parents []*Node = []*Node{root}
	for i := 1; i < depth; i++ {

		// fmt.Print(i, " parents: ")
		// for _, node := range parents {
		// 	fmt.Print(node.Name, " ")
		// }
		// fmt.Println()

		var children []*Node
		for _, node := range parents {
			// fmt.Println("depth", i, node.Name)
			addChildren(&node, maxChildren, name)
			for _, child := range node.children {
				children = append(children, child)
			}
		}

		// fmt.Print(i, " children: ")
		// for _, node := range children {
		// 	fmt.Print(node.Name, " ")
		// }
		// fmt.Println()

		parents = children
	}

	return root
}

func addChildren(root **Node, maxChildren int, name func() string) {
	childrenNum := rand.Intn(maxChildren) + 1
	// fmt.Println((*root).Name, childrenNum)

	for i := 0; i < childrenNum; i++ {
		child := &Node{Name: name()}
		child.parent = *root
		(*root).children = append((*root).children, child)
	}
}

func newName(prefix string) func() string {
	var count int
	return func() string {
		name := fmt.Sprintf("%s-%d", prefix, count)
		count++

		return name
	}
}
