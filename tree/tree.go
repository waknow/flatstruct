package tree

import (
	"fmt"
	"math/rand"
	"strings"
)

type Treer interface {
	ID() string
	Parent() Treer
	Children() []Treer
	String() string

	SetParent(Treer)
	SetChildren([]Treer)
}

func PrintTree(t Treer) {
	printTree([]bool{}, t)
}

func RandomPath(root Treer) string {
	if root == nil {
		return ""
	}
	var parts = []string{root.ID()}
	node := root
	for {
		if len(node.Children()) == 0 {
			break
		}

		child := node.Children()[rand.Intn(len(node.Children()))]
		parts = append(parts, child.ID())
		node = child
	}

	return strings.Join(parts, ".")
}

func ReversePath(path string) string {
	parts := strings.Split(path, ".")
	for i, j := 0, len(parts)-1; i < j; i, j = i+1, j-1 {
		parts[i], parts[j] = parts[j], parts[i]
	}
	return strings.Join(parts, ".")
}

func FindNode(root Treer, path string) Treer {
	var node Treer
	ps := []Treer{root}
	for _, part := range strings.Split(path, ".") {
		mached := false
		for _, p := range ps {
			if p.ID() == part {
				mached = true
				node = p
				ps = p.Children()
				break
			}
		}
		if !mached {
			return nil
		}
	}

	return node
}

func RebuildTreeByNode(node Treer) {
	if node == nil || node.Parent() == nil {
		return
	}

	parent := node.Parent()

	siblings := siblings(node)
	// fmt.Println(node.ID(), parent, "siblings:", siblings)
	parent.SetChildren(nil)
	node.SetChildren(append([]Treer{parent}, siblings...))
	node.SetParent(nil)

	RebuildTreeByNode(parent)
	parent.SetParent(node)
}

func siblings(node Treer) []Treer {
	if node == nil || node.Parent() == nil {
		return nil
	}

	var siblings []Treer
	for _, child := range node.Parent().Children() {
		if child != node {
			siblings = append(siblings, child)
		}
	}

	return siblings
}

func printTree(prefixes []bool, t Treer) {
	fmt.Print(getPrefix(prefixes), t, "\n")
	for idx, child := range t.Children() {
		printTree(append(prefixes, idx != len(t.Children())-1), child)
	}
}

func getPrefix(prefixes []bool) string {
	l := len(prefixes)

	if l == 0 {
		return ""
	}

	last := prefixes[l-1]
	prefixes = prefixes[:l-1]

	var s string
	for _, prefix := range prefixes {
		if prefix {
			s += "│  "
		} else {
			s += "   "
		}
	}

	if last {
		s += "├─ "
	} else {
		s += "└─ "
	}

	return s
}
