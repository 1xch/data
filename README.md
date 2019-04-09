# data

General purpose custom data containment for use with the Go programming language.

## Item

An interface (and package default specific structure) for managing Key & Value 
for data in addition to providing facility for transmission and cloning. 

## Trie

A trie constructed from https://github.com/tchap/go-patricia with a variety of 
changes and additions(and more likely the  opposite of optimisations).

## Vector

A sync.Mutex bound struct wrapping an Item Trie with customized control.

## Store

A structure for managing Vectors across varied formats(e.g. json, yaml, etc.).
