package main

import (
	"fmt"

	"github.com/iknite/molletree/balloon/hyper"
	"github.com/iknite/molletree/encoding/encstring"
)

func main() {
	tree := hyper.NewTree()
	fmt.Println(tree.Add(encstring.ToBytes("Hello world!")))
}
