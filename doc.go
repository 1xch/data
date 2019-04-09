// Package data contains generalized data management tactics & strategies.
// Currently this package provides four interrelated structures:
//
// - Trie
//   A trie constructed from https://github.com/tchap/go-patricia with a
//   variety of changes and additions. No optimisation is promised.
//
// - Item
//   A general interface for managing a key and a value. A key is a string and
//   a value may be anything. The package provides for a variety of common types,
//   providing example for any type you might have need to construct.
//
// - Vector
//   A sync.Mutex bound struct wrapping a Trie holding any number of package
//   level Item.
//
// - Store
//   An interface for managing the storage of Vector in and out of any variety
//   of formats. Package provides common stores to take a Vector to stdout(out
//   only), json, formatted json, and yaml. An example use might take a Vector
//   to json, sent elsewhere and modified, returned and used as a Vector, viewed
//   in a terminal, saved as yaml and returned Vector, etc et al. Store is meant
//   as a rough data interchange manager mediating Vector to any format you might
//   need or want.
package data
