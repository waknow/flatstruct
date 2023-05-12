package main

import (
	"fmt"

	"struct-flat/flat"
	"struct-flat/tree"
	"struct-flat/util"
)

type Address struct {
	Street string
	City   string
}

type Person struct {
	Name    string
	Age     int
	Aliases []string
	Addrs   []Address
}

func main() {
	p := Person{
		Name: "John",
		Age:  30,
		Aliases: []string{
			"Johnny",
			"Johny",
		},
		Addrs: []Address{
			{
				Street: "123 Main St",
				City:   "Anytown",
			},
			{
				Street: "456 Main St",
				City:   "Anytown",
			},
		},
	}

	result, keys, err := flat.Flat(p, "root.Addrs")
	if err != nil {
		fmt.Println("flat", err)
		return
	}

	if len(result) == 0 {
		fmt.Println("empty result")
		return
	}

	util.PrettyMapOrder(keys, result...)
}

func randTree() {
	var titles []string
	var trees []tree.Noder

	r := tree.RandomTree(4, 3)
	titles = append(titles, "random tree")
	trees = append(trees, tree.Clone(r, tree.NewNode))

	path := tree.RandomPath(r)
	fmt.Println("random path", path)
	nr := tree.FindNode(r, path)

	tree.Rebuild(nr)
	titles = append(titles, "rebuild by "+nr.ID())
	trees = append(trees, tree.Clone(nr, tree.NewNode))

	tree.Rebuild(r)
	titles = append(titles, "rebuild by "+r.ID())
	trees = append(trees, tree.Clone(r, tree.NewNode))

	tree.Prints(trees, titles...)
}
