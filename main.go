package main

import (
	"fmt"

	"github.com/iknite/molletree/history"
)

func main() {
	tree := history.NewTree()
	commitment := tree.Add("Hello world!")
	fmt.Println(commitment)
}
