package emojiart

import "testing"
import "fmt"

func TestTreeJSON(t *testing.T) {
	output, err := treeJSON()
	str := string(output[:])
	fmt.Println(str)
	fmt.Println(err)
}

func TestBuildTree(t *testing.T) {
	root := BuildTree("apple", 0, len(emojiDictAvg["apple"])-1, 0)
	fmt.Print("this is the root ", root)
}
