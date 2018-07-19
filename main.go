package main

import (
	"fmt"

	"github.com/iknite/bygone-tree/history"
)

func main() {
	tree := history.NewTree()
	commitment := tree.Add("Hello world!")
	fmt.Println(commitment)
}
