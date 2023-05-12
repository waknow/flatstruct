package flat

import (
	"encoding/json"
	"fmt"
	"reflect"
	"sort"
	"strings"
	"struct-flat/tree"
	"struct-flat/util"
)

func Flat(i interface{}, primaryPath string) ([]map[string]interface{}, []string, error) {
	var root = &value{
		name:  "root",
		kind:  kindOf(i),
		value: i,
	}
	travelInterface(&root, i)

	primaryNode := tree.FindNode(root, primaryPath)
	if primaryNode == nil {
		return nil, nil, fmt.Errorf("cannot find node %s", primaryPath)
	}

	primaryValues := primaryNode.Children()
	if len(primaryValues) == 0 {
		return nil, nil, nil
	}

	tree.Rebuild(primaryNode)

	primaryRows, err := toRows(primaryValues)
	if err != nil {
		return nil, nil, err
	}

	rows, err := toRow("", primaryNode)
	if err != nil {
		return nil, nil, err
	}

	primaryKeys := sort.StringSlice(util.MapKeys(primaryRows[0]))
	rowKeys := sort.StringSlice(util.MapKeys(rows))

	keys := append(primaryKeys, rowKeys...)

	var result []map[string]interface{}
	for _, primaryRow := range primaryRows {
		m := map[string]interface{}{}
		for _, k := range keys {
			if v, ok := primaryRow[k]; ok {
				m[k] = v
				continue
			}
			if v, ok := rows[k]; ok {
				m[k] = v
				continue
			}
		}
		result = append(result, m)
	}

	return result, keys, nil
}

func toRows(nodes []tree.Noder) ([]map[string]interface{}, error) {
	values := []map[string]interface{}{}
	for _, node := range nodes {
		v := node.(*value)
		if v.kind != kindObject {
			return nil, fmt.Errorf("node %s kind is not object", v.name)
		}
		m := map[string]interface{}{}
		for _, child := range node.Children() {
			v := child.(*value)
			if v.kind != kindValue {
				return nil, fmt.Errorf("node %s kind is not value", v.name)
			}
			m[strings.ToLower(v.name)] = v.value
		}
		values = append(values, m)
	}

	return values, nil
}

func toRow(path string, node tree.Noder) (map[string]interface{}, error) {
	path = joinPath(path, node.ID())
	var result = map[string]interface{}{}
	// fmt.Println("toRow", node.ID(), len(node.Children()))
	for _, child := range node.Children() {
		value := child.(*value)
		// fmt.Println("toRow", node.ID(), "->", value.name, value.kind, len(value.children))

		//skip original root node
		if value.kind == kindObject && len(value.children) == 0 {
			continue
		}

		switch value.kind {
		case kindObject:
			m, err := toRow(path, child)
			if err != nil {
				return nil, err
			}
			for k, v := range m {
				result[k] = v
			}
		case kindArray:
			bs, err := json.Marshal(value.value)
			if err != nil {
				return nil, fmt.Errorf("marshal array %s error: %s", value.name, err)
			}
			result[joinPath(path, value.name)] = string(bs)
		case kindValue:
			result[joinPath(path, value.name)] = value.value
		default:
			return nil, fmt.Errorf("unknown kind %s", value.kind)
		}

		// fmt.Println("toRow", result)
	}

	return result, nil
}

func joinPath(path string, s string) string {
	s = strings.ToLower(s)
	if path == "" {
		return s
	}
	if s == "" {
		return path
	}

	return path + "." + s
}

type value struct {
	name  string
	kind  string
	value interface{}

	parent   *value
	children []*value
}

func NewValue(t tree.Noder) tree.Noder {
	if t == nil {
		return nil
	}

	v, ok := t.(*value)
	if !ok {
		return nil
	}

	return &value{
		name:  v.name,
		kind:  v.kind,
		value: v.value,
	}
}

func (v value) ID() string {
	return v.name
}

func (v value) Parent() tree.Noder {
	if v.parent == nil {
		return nil
	}
	return v.parent
}

func (v value) Children() []tree.Noder {
	var children []tree.Noder
	for _, child := range v.children {
		children = append(children, child)
	}
	return children
}

func (v value) String() string {
	if len(v.children) == 0 {
		return fmt.Sprintf("%s(%s): %v", v.name, v.kind, v.value)
	}
	return fmt.Sprintf("%s(%s)", v.name, v.kind)
}

func (v *value) SetParent(parent tree.Noder) {
	if v == nil {
		return
	}
	if parent == nil {
		v.parent = nil
		return
	}
	v.parent = parent.(*value)
}

func (v *value) SetChildren(children []tree.Noder) {
	if v == nil {
		return
	}
	v.children = nil
	for _, child := range children {
		v.children = append(v.children, child.(*value))
	}
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

const (
	kindValue  = "value"
	kindArray  = "array"
	kindObject = "object"
)

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
