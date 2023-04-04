package main

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
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
	fmt.Println("Hello, World!")

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

	var root = &value{
		name:  "root",
		kind:  kindOf(p),
		value: p,
	}
	travelInterface(&root, p)
	printTree(0, root)

	path := "root.Addrs"
	values := findValues(root, path)
	fmt.Println("find values", path, values)

	for _, value := range values {
		names, values, err := findRelateValues(value)
		if err != nil {
			panic(err)
		}
		fmt.Println("names", names)
		for _, value := range values {
			fmt.Println("values", value)
		}
	}
}

func findRelateValues(v *value) ([]string, [][]interface{}, error) {
	if v.kind == kindObject {
		return nil, nil, fmt.Errorf("cannot find relate values for object")
	}

	names, primaryValues, err := flatValue(v)
	if err != nil {
		return nil, nil, err
	}
	fmt.Println("names", names)
	fmt.Println("primarys", primaryValues)

	var relateValues []interface{}
	p := v.parent
	for p != nil {
		switch p.kind {
		case kindObject:
			for _, child := range p.children {
				if child == v {
					continue
				}

				switch child.kind {
				case kindValue:
					relateValues = append(relateValues, child.value)
					name, err := getName(child, names)
					if err != nil {
						return nil, nil, err
					}
					names = append(names, name)
				case kindObject, kindArray:
					bs, err := json.Marshal(child.value)
					if err != nil {
						return nil, nil, err
					}
					relateValues = append(relateValues, string(bs))
					name, err := getName(child, names)
					if err != nil {
						return nil, nil, err
					}
					names = append(names, name)
				default:
					return nil, nil, fmt.Errorf("unknown kind %s", child.kind)
				}
			}
		case kindArray:
			for _, child := range p.children {
				if child == v {
					continue
				}

				_, values, err := flatValue(child)
				if err != nil {
					return nil, nil, err
				}
				primaryValues = append(primaryValues, values...)
			}
		default:
			return nil, nil, fmt.Errorf("cannot find relate values for %s", p.kind)
		}

		p = p.parent
	}

	var values [][]interface{}

	for _, primaryValue := range primaryValues {
		var value []interface{}
		value = append(value, primaryValue...)
		value = append(value, relateValues...)
		values = append(values, value)
	}

	return names, values, nil
}

func flatValue(v *value) ([]string, [][]interface{}, error) {
	var names []string
	var valuesList [][]interface{}
	switch v.kind {
	case kindValue:
		names = []string{v.name}
		valuesList = [][]interface{}{{v.value}}
	case kindArray:
		for idx, child := range v.children {
			var values []interface{}
			for _, grandChild := range child.children {
				if idx == 0 {
					name, err := getName(grandChild, names)
					if err != nil {
						return nil, nil, err
					}
					names = append(names, name)
				}
				values = append(values, grandChild.value)
			}
			valuesList = append(valuesList, values)
		}
	default:
		return nil, nil, fmt.Errorf("unknown kind %s", v.kind)
	}

	return names, valuesList, nil
}

func getName(v *value, existedNames []string) (string, error) {
	if !hasString(v.name, existedNames) {
		return v.name, nil
	}

	p := v.parent
	var name string
	for p != nil {
		name = p.name + "." + name
		if !hasString(name, existedNames) {
			return name, nil
		}
		p = p.parent
	}

	return "", fmt.Errorf("cannot find a unique name for %s", v.name)
}

func hasString(s string, ss []string) bool {
	for _, v := range ss {
		if s == v {
			return true
		}
	}
	return false
}

func findValues(root *value, path string) []*value {
	p := root

	subs := strings.Split(path, ".")
	for idx, sub := range subs {
		if p.name == sub {
			if idx == len(subs)-1 {
				return []*value{p}
			}
			for _, child := range p.children {
				if child.name == subs[idx+1] {
					p = child
					break
				}
			}
		}
	}

	return nil
}

func printTree(level int, root *value) {
	if root.kind == kindValue {
		fmt.Printf("%s%s(%s): %v\n", strings.Repeat("\t", level), root.name, root.kind, root.value)
		return
	}
	fmt.Printf("%s%s(%s)\n", strings.Repeat("\t", level), root.name, root.kind)
	for _, child := range root.children {
		printTree(level+1, child)
	}
}

const (
	kindValue  = "value"
	kindArray  = "array"
	kindObject = "object"
)

type value struct {
	name  string
	kind  string
	value interface{}

	parent   *value
	children []*value
}

func (v value) String() string {
	return fmt.Sprintf("%s(%s): %v", v.name, v.kind, v.value)
}

func travelInterface(root **value, i interface{}) {
	switch (*root).kind {
	case kindValue:
	case kindArray:
		(*root).children = unpackArray(i)
	case kindObject:
		(*root).children = unpackObject(i)
	}

	for _, child := range (*root).children {
		child.parent = *root

		if child.kind != kindValue {
			travelInterface(&child, child.value)
		}
	}
}

func unpackObject(i interface{}) []*value {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	if t.Kind() != reflect.Struct {
		panic(fmt.Sprintf("expected struct, got %s", t.Kind()))
	}

	var values []*value
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		values = append(values, &value{
			name:  field.Name,
			kind:  kindOf(v.Field(i).Interface()),
			value: v.Field(i).Interface(),
		})
	}

	return values
}

func unpackArray(i interface{}) []*value {
	t := reflect.TypeOf(i)
	v := reflect.ValueOf(i)

	if t.Kind() != reflect.Slice && t.Kind() != reflect.Array {
		panic(fmt.Sprintf("expected slice, got %s", t.Kind()))
	}

	var values []*value
	for i := 0; i < v.Len(); i++ {
		values = append(values, &value{
			name:  fmt.Sprintf("%d", i),
			kind:  kindOf(v.Index(i).Interface()),
			value: v.Index(i).Interface(),
		})
	}

	return values
}

func kindOf(v interface{}) string {
	t := reflect.TypeOf(v)

	switch t.Kind() {
	case reflect.Array, reflect.Slice:
		return kindArray
	case reflect.Struct:
		return kindObject
	case reflect.String,
		reflect.Float32, reflect.Float64,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return kindValue
	default:
		panic(fmt.Sprintf("unknown type: %s", t.Kind()))
	}
}