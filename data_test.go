package data

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

var (
	currentDir             string
	jsonLoc, yamlLoc       string
	rs                     string
	si, bi, ii, fi, li, mi Item
	testItems              []Item
	base                   *Container
)

func init() {
	currentDir, _ = os.Getwd()
	rs = strings.Join([]string{"json", currentDir, "container"}, ",")
	jsonLoc = filepath.Join(currentDir, fmt.Sprintf("%s.%s", "container", "json"))
	yamlLoc = filepath.Join(currentDir, fmt.Sprintf("%s.%s", "container", "yaml"))
	si = NewItem("a.string", "string")
	bi = NewItem("a.bool", "true")
	ii = NewItem("a.int", "9")
	fi = NewItem("a.float", "9.9")
	li = NewItem("a.list", "a,b,c,d,e,f,g")
	mi = NewItem("a.map", "a:one,b:two,c:3,d:d")
	always1 := NewItem("always.1", "the same")
	always2 := NewItem("always.2", "the same")
	testItems = []Item{
		NewItem("store.retrieval.string", rs),
		si, bi, ii, fi, li, mi, always1, always2,
	}
	base = testContainer().Clone()
}

func testContainer() *Container {
	c := New("TEST")
	c.Set(
		testItems...,
	)
	return c
}
