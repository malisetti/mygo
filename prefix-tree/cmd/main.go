package main

import (
	"fmt"

	prefixtree "prefix-tree/trie"
)

func main() {
	trie := prefixtree.NewTrie[rune]()

	for _, fruit := range fruits {
		trie.Insert([]rune(fruit))
	}

	fmt.Printf("digraph trie {\n")
	prefixtree.DumpDot('_', trie, func(r rune) string {
		return string(r)
	})
	fmt.Printf("}\n")
}
