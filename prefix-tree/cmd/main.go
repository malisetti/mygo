package main

import (
	"fmt"

	prefixtree "prefix-tree/trie"
)

func main() {
	trie := prefixtree.NewTrie()

	for _, fruit := range fruits {
		trie.Insert(fruit)
	}

	fmt.Printf("digraph trie {\n")
	prefixtree.DumpDot('_', trie)
	fmt.Printf("}\n")
}
