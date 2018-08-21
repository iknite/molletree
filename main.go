package main

import (
	"fmt"

	"github.com/iknite/molletree/balloon"
)

func main() {
	balloon := balloon.NewBalloon()
	fmt.Println(balloon.Add("Hello world!"))
}
