package util

import (
	"fmt"
	"strings"
)

func PrettyMapOrder[Key comparable, Value any](keyStrs []string, datas ...map[Key]Value) {
	if len(datas) == 0 {
		return
	}
	keys := MapKeys(datas[0])

	maxLens := make(map[string]int)
	for _, keyStr := range keyStrs {
		maxLens[keyStr] = len(keyStr)
	}

	colsStrs := []map[string]string{}
	for _, data := range datas {
		colStrs := map[string]string{}

		for _, key := range keys {
			str := fmt.Sprintf("%v", data[key])
			colStrs[fmt.Sprintf("%v", key)] = str
			maxLen := maxLens[fmt.Sprintf("%v", key)]
			if maxLen < len(str) {
				maxLens[fmt.Sprintf("%v", key)] = len(str)
			}
		}

		colsStrs = append(colsStrs, colStrs)
	}

	var formats []string

	for _, key := range keyStrs {
		formats = append(formats, fmt.Sprintf("%%-%ds", maxLens[key]))
	}

	format := strings.Join(formats, "    ") + "\n"

	fmt.Printf(format, All2Any(keyStrs)...)
	for _, colStrs := range colsStrs {

		ss := []string{}
		for _, key := range keyStrs {
			ss = append(ss, colStrs[key])
		}

		fmt.Printf(format, All2Any(ss)...)
	}
}

func PrettyMap[Key comparable, Value any](datas ...map[Key]Value) {
	if len(datas) == 0 {
		return
	}

	keyStrs := All2Strings(MapKeys(datas[0]))

	PrettyMapOrder(keyStrs, datas...)
}

func All2Strings[v any](vs []v) []string {
	var result []string
	for _, v := range vs {
		result = append(result, fmt.Sprintf("%v", v))
	}
	return result
}

func All2Any[v any](vs []v) []any {
	var r []any
	for _, v := range vs {
		r = append(r, v)
	}
	return r
}

func MapKeys[Key comparable, Value any](m map[Key]Value) []Key {
	keys := make([]Key, len(m))
	i := 0
	for k := range m {
		keys[i] = k
		i++
	}
	return keys
}
