package trie

import (
	"fmt"
	"sort"
)

type Trie struct {
	root *node
}

type node struct {
	children map[rune]*node
}

func NewTrie() *Trie {
	return &Trie{
		root: &node{},
	}
}

func (trie *Trie) Insert(str string) {
	root := trie.root
	for _, v := range str {
		if root.children == nil {
			root.children = make(map[rune]*node)
		}
		if root.children[v] == nil {
			// insert node and advance
			root.children[v] = &node{
				children: make(map[rune]*node),
			}
		}

		root = root.children[v]
	}
}

func (trie *Trie) Check(str string) (isWord, isSubStr bool) {
	root := trie.root
	for _, v := range str {
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

func (trie *Trie) Completions(str string) ([]string, error) {
	root := trie.root
	for _, v := range str {
		root = root.children[v]
		if root == nil {
			return nil, fmt.Errorf("%s word not found", str)
		}
	}
	words := words(root, str)
	return words, nil
}

func words(root *node, str string) []string {
	var words []string
	var search func(*node, []rune)
	search = func(node *node, str []rune) {
		if len(root.children) == 0 && len(str) > 0 {
			words = append(words, string(str))
		} else {
			for r, child := range node.children {
				search(child, append(str, r))
			}
		}
	}
	search(root, []rune(str))
	return words
}

func DumpDot(rootc rune, trie *Trie) {
	var dump func(rune, *node)
	dump = func(from rune, node *node) {
		var keys []rune
		for to := range node.children {
			keys = append(keys, to)
		}
		if len(keys) == 0 {
			fmt.Printf("	\"%c\" -> \"%c\";\n", from, '*')
			return
		}
		sort.Slice(keys, func(i, j int) bool {
			return keys[i] < keys[j]
		})
		for _, to := range keys {
			fmt.Printf("    \"%c\" -> \"%c\";\n", from, to)
			dump(to, node.children[to])
		}
	}
	dump(rootc, trie.root)
}
