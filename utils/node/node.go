package node

import (
	"fmt"
	"io"

	"github.com/Conquest-Reforged/ReforgedLauncher/utils/files"
)

var empty = &Node{}

type Node struct {
	value interface{}
}

func FromJson(r io.Reader) *Node {
	var value interface{}
	e := files.ParseJson(r, &value)
	if e != nil {
		return empty
	}
	return &Node{value: value}
}

func FromJsonAt(path string) *Node {
	var value interface{}
	e := files.ParseJsonAt(path, &value)
	if e != nil {
		return empty
	}
	return &Node{value: value}
}

func (n *Node) ToJson(w io.Writer) error {
	return files.WriteJson(w, n.value)
}

func (n *Node) ToJsonAt(path string) error {
	return files.WriteJsonAt(path, n.value)
}

func (n *Node) String() string {
	switch t := n.value.(type) {
	case string:
		return t
	default:
		return fmt.Sprint(n.value)
	}
}

func (n *Node) Bool() bool {
	switch t := n.value.(type) {
	case bool:
		return t
	default:
		return false
	}
}

func (n *Node) Int() int {
	switch t := n.value.(type) {
	case int:
		return t
	default:
		return 0
	}
}

func (n *Node) Float() float64 {
	switch t := n.value.(type) {
	case float64:
		return t
	default:
		return 0
	}
}

func (n *Node) Empty() bool {
	return n.value == nil
}

func (n *Node) Get(key interface{}) *Node {
	if n.value != nil {
		switch k := key.(type) {
		case string:
			return n.objGet(k)
		case int:
			return n.arrGet(k)
		}
	}
	return empty
}

func (n *Node) Set(key interface{}, value interface{}) bool {
	if n.value != nil {
		switch k := key.(type) {
		case string:
			return n.objSet(k, value)
		case int:
			return n.arrSet(k, value)
		}
	}
	return false
}

func (n *Node) objGet(key string) *Node {
	switch v := n.value.(type) {
	case map[string]interface{}:
		if i, ok := v[key]; ok && i != nil {
			return &Node{value: i}
		}
	}
	return empty
}

func (n *Node) arrGet(index int) *Node {
	if index >= 0 {
		switch v := n.value.(type) {
		case []interface{}:
			if index < len(v) {
				i := v[index]
				if i != nil {
					return &Node{value: i}
				}
			}
		}
	}
	return empty
}

func (n *Node) objSet(key string, value interface{}) bool {
	switch v := n.value.(type) {
	case map[string]interface{}:
		v[key] = value
		return true
	}
	return false
}

func (n *Node) arrSet(index int, value interface{}) bool {
	if index >= 0 {
		switch v := n.value.(type) {
		case []interface{}:
			if index < len(v) {
				v[index] = value
				return true
			}
		}
	}
	return false
}
