package tree

import (
	"fmt"
	"math/rand"
	"strings"
)

type Noder interface {
	ID() string
	Parent() Noder
	Children() []Noder
	String() string

	SetParent(Noder)
	SetChildren([]Noder)
}

func Clone(t Noder, f func(t Noder) Noder) Noder {
	if t == nil {
		return nil
	}

	nt := f(t)
	nt.SetParent(nil)
	nt.SetChildren(nil)

	for _, child := range t.Children() {
		nt.SetChildren(append(nt.Children(), Clone(child, f)))
	}

	return nt
}

func Print(t Noder) {
	print([]bool{}, t)
}

func Prints(ts []Noder, titles ...string) {
	if titles == nil {
		titles = []string{}
	}

	sss := [][]string{}
	for i, t := range ts {
		ss := sprint([]bool{}, t)
		if len(titles) > 0 {
			if i < len(titles) {
				ss = append(ss, "", titles[i], "")
			} else {
				ss = append(ss, "", "", "")
			}
		}
		sss = append(sss, ss)
	}

	var formats []string
	var maxLine int
	for _, ss := range sss {
		max := 0
		if len(ss) > maxLine {
			maxLine = len(ss)
		}
		for _, s := range ss {
			if len(s) > max {
				max = len(s)
			}
		}
		formats = append(formats, fmt.Sprintf("%%-%ds", max))
	}

	format := strings.Join(formats, " ")
	for i := 0; i < maxLine; i++ {
		var args []interface{}
		for _, ss := range sss {
			if i < len(ss) {
				args = append(args, ss[i])
			} else {
				args = append(args, "")
			}
		}
		// fmt.Println("format", format)
		// fmt.Println("args", args)
		fmt.Printf(format+"\n", args...)
	}
}

func RandomPath(root Noder) string {
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

func FindNode(root Noder, path string) Noder {
	var node Noder
	ps := []Noder{root}
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

func Rebuild(node Noder) {
	if node == nil || node.Parent() == nil {
		return
	}

	parent := node.Parent()

	children := rebuildChildren(node)
	parent.SetChildren(nil)
	node.SetChildren(children)
	node.SetParent(nil)

	Rebuild(parent)
	parent.SetParent(node)
}

func rebuildChildren(node Noder) []Noder {
	if node == nil || node.Parent() == nil {
		return nil
	}

	parent := node.Parent()

	if len(parent.Children()) == 0 {
		return []Noder{parent}
	}

	var children []Noder
	for _, child := range parent.Children() {
		if child == node {
			children = append(children, parent)
		} else {
			children = append(children, child)
		}
	}

	return children
}

func print(prefixes []bool, t Noder) {
	fmt.Print(getPrefix(prefixes), t, "\n")
	for idx, child := range t.Children() {
		print(append(prefixes, idx != len(t.Children())-1), child)
	}
}

func sprint(prefixes []bool, t Noder) []string {
	var ss []string
	ss = append(ss, getPrefix(prefixes)+t.String())
	for idx, child := range t.Children() {
		ss = append(ss, sprint(append(prefixes, idx != len(t.Children())-1), child)...)
	}
	return ss
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
