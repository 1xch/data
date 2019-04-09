package data

import (
	"fmt"
	"math/rand"
	"os"
	"path/filepath"
	"time"
)

var (
	currentDir                              string
	jsonLoc, yamlLoc                        string
	rs                                      []string
	si, ssi, bi, ii, ii64, ui, ui64, fi, vi Item
	testItems                               []Item
	base                                    *Vector
)

func testVector() *Vector {
	c := New("TEST")
	c.Set(
		testItems...,
	)
	return c
}

func init() {
	currentDir, _ = os.Getwd()
	rs = []string{"json", currentDir, "vector"}
	jsonLoc = filepath.Join(currentDir, fmt.Sprintf("%s.%s", "vector", "json"))
	yamlLoc = filepath.Join(currentDir, fmt.Sprintf("%s.%s", "vector", "yaml"))
	si = NewStringItem("a.string", "string")
	ssi = NewStringsItem("a.list", "a", "b", "c")
	bi = NewBoolItem("a.bool", false)
	ii = NewIntItem("a.int", 9)
	ii64 = NewInt64Item("a.int64", 999)
	ui = NewUintItem("a.uint", 1)
	ui64 = NewUint64Item("a.uint64", 111)
	fi = NewFloat64Item("a.float", 9.9)
	v := New("multi")
	v.Set(
		NewStringItem("vector.1", "ONE"),
		NewStringItem("vector.2", "TWO"),
	)
	vi = NewVectorItem(v.Tag(), v)
	always1 := NewStringItem("always.1", "the same")
	always2 := NewStringItem("always.2", "")
	always2.Provide("the same")
	ri := NewStringsItem("store.retrieval.string", rs...)
	testItems = []Item{
		ri, si, ssi, bi, ii, ii64, ui, ui64, fi, vi, always1, always2,
	}
	base = testVector().Clone()
	rand.Seed(time.Now().UnixNano())
}
