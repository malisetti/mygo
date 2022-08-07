package trie

import (
	"fmt"
	"sort"

	"golang.org/x/exp/constraints"
)

type Trie[V constraints.Ordered] struct {
	root *node[V]
}

type node[V constraints.Ordered] struct {
	children map[V]*node[V]
}

func NewTrie[V constraints.Ordered]() *Trie[V] {
	return &Trie[V]{
		root: &node[V]{},
	}
}

func (trie *Trie[V]) Insert(xs []V) {
	root := trie.root
	for _, v := range xs {
		if root.children == nil {
			root.children = make(map[V]*node[V])
		}
		if root.children[v] == nil {
			// insert node and advance
			root.children[v] = &node[V]{
				children: make(map[V]*node[V]),
			}
		}

		root = root.children[v]
	}
}

func (trie *Trie[V]) Check(xs []V) (isWord, isSubStr bool) {
	root := trie.root
	for _, v := range xs {
		if root.children[v] == nil {
			return
		} else {
			root = root.children[v]
		}
	}

	for range root.children {
		return false, true
	}

	return true, false
}

func (trie *Trie[V]) Completions(xs []V) ([]V, error) {
	root := trie.root
	for _, v := range xs {
		root = root.children[v]
		if root == nil {
			return nil, fmt.Errorf("word not found")
		}
	}
	words := words(root, xs)
	return words, nil
}

func words[V constraints.Ordered](root *node[V], xs []V) []V {
	var words []V
	var search func(*node[V], []V)
	search = func(node *node[V], str []V) {
		if len(root.children) == 0 && len(str) > 0 {
			words = append(words, str...)
		} else {
			for r, child := range node.children {
				search(child, append(str, r))
			}
		}
	}
	search(root, xs)
	return words
}

type StringRep[V constraints.Ordered] func(V) string

func DumpDot[V constraints.Ordered](rootc V, trie *Trie[V], stringRep func(V) string) {
	var dump func(V, *node[V])
	dump = func(from V, node *node[V]) {
		var keys []V
		for to := range node.children {
			keys = append(keys, to)
		}
		if len(keys) == 0 {
			fmt.Printf("	\"%s\" -> \"%s\";\n", stringRep(from), "*")
			return
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		for _, to := range keys {
			fmt.Printf("    \"%s\" -> \"%s\";\n", stringRep(from), stringRep(to))
			dump(to, node.children[to])
		}
	}
	dump(rootc, trie.root)
}
