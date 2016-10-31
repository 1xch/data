package data

import (
	"fmt"
	"os"
	"path/filepath"
)

var (
	currentDir             string
	jsonLoc, yamlLoc       string
	rs                     []string
	si, bi, ii, fi, li, mi Item
	testItems              []Item
	base                   *Container
)

func init() {
	currentDir, _ = os.Getwd()
	rs = []string{"json", currentDir, "container"}
	jsonLoc = filepath.Join(currentDir, fmt.Sprintf("%s.%s", "container", "json"))
	yamlLoc = filepath.Join(currentDir, fmt.Sprintf("%s.%s", "container", "yaml"))
	si = NewStringItem("a.string", "string")
	li = NewStringsItem("a.list", "a", "b", "c")
	bi = NewBoolItem("a.bool", false)
	ii = NewIntItem("a.int", 9)
	fi = NewFloatItem("a.float", 9.9)
	c := New("multi")
	c.Set(
		NewStringItem("multi.1", "ONE"),
		NewStringItem("multi.2", "TWO"),
	)
	mi = NewMultiItem(c.Tag(), c)
	always1 := NewStringItem("always.1", "the same")
	always2 := NewStringItem("always.2", "")
	always2.Provide("the same")
	ri := NewStringsItem("store.retrieval.string", rs...)
	testItems = []Item{
		ri, si, li, bi, ii, fi, mi, always1, always2,
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
