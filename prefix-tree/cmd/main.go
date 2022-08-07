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
	var stringRep prefixtree.StringRep[rune] = func(r rune) string {
		return string(r)
	}
	prefixtree.DumpDot('_', trie, stringRep)
	fmt.Printf("}\n")
}
