package data

import (
	"bytes"
	"encoding/json"
	"strconv"
	"sync"
)

// A sync.Mutex bound struct that wraps a Trie holding package level Item.
type Vector struct {
	l  *sync.RWMutex
	o  []Option
	bl []string
	*Trie
}

//
func New(tag string, o ...Option) *Vector {
	t := NewTrie(o...)
	v := &Vector{
		nil, o, make([]string, 0), t,
	}
	v.mutexSet()
	v.Set(NewStringItem("vector.tag", tag))
	return v
}

func (v *Vector) mutexSet() {
	if v.l == nil {
		v.l = &sync.RWMutex{}
	}
}

func (v *Vector) trieSet() {
	if v.Trie == nil {
		v.Trie = NewTrie(v.o...)
	}
}

// unmarshaling a raw Vector can be painful in certain situations,
// this ensures the process is less so
func (v *Vector) ensureNotEmpty() {
	v.mutexSet()
	v.trieSet()
}

//
func (v *Vector) Tag() string {
	return v.ToString("vector.tag")
}

//
func (v *Vector) Retag(t string) {
	i := v.Get("vector.tag")
	if i != nil {
		i.Provide(t)
	}
}

//
func (v *Vector) Keys() []string {
	var ret []string
	w := func(p Prefix, i Item) error {
		ret = append(ret, string(p))
		return nil
	}
	v.walk(nil, w)
	return ret
}

//
func (v *Vector) Get(k string) Item {
	v.l.RLock()
	key := Prefix(k)
	i := v.get(key)
	v.l.RUnlock()
	return i
}

//
func (v *Vector) Match(k string) []Item {
	v.l.RLock()
	defer v.l.RUnlock()
	var ret []Item
	bk := []byte(k)
	w := func(p Prefix, i Item) error {
		if bytes.Contains(p, bk) {
			ret = append(ret, i)
		}
		return nil
	}
	v.walk(nil, w)
	return ret
}

func (v *Vector) Blacklist(keys ...string) {
	v.bl = append(v.bl, keys...)
}

func inList(k string, bl []string) bool {
	for _, v := range bl {
		if k == v {
			return true
		}
	}
	return false
}

func notBlacklisted(bl []string, i []Item) []Item {
	var ret []Item
	for _, ii := range i {
		if !inList(ii.Key(), bl) {
			ret = append(ret, ii)
		}
	}
	return ret
}

//
func (v *Vector) Set(i ...Item) {
	nbi := notBlacklisted(v.bl, i)
	v.l.Lock()
	v.set(nbi...)
	v.l.Unlock()
}

//
func (v *Vector) Merge(vs ...*Vector) {
	for _, vv := range vs {
		l := vv.List("vector.tag", "vector.id")
		v.Set(l...)
	}
}

//
func (v *Vector) Clone(except ...string) *Vector {
	except = append(except, "vector.tag")
	n := New(v.Tag())
	l := v.List(except...)
	var nl []Item
	for _, i := range l {
		nl = append(nl, i.Clone())
	}
	n.Set(nl...)
	return n
}

//
func (v *Vector) CloneAs(tag string, except ...string) *Vector {
	nc := v.Clone(except...)
	nc.Retag(tag)
	return nc
}

func match(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

// Returns a list of Item, EXCEPT those matching the provided key strings.
func (v *Vector) List(except ...string) []Item {
	v.l.RLock()
	defer v.l.RUnlock()
	var ret []Item
	w := func(p Prefix, i Item) error {
		if !match(except, i.Key()) {
			ret = append(ret, i)
		}
		return nil
	}
	v.walk(nil, w)
	return ret
}

// Clears the Vector of all Item.
func (v *Vector) Clear() {
	v.reset()
}

// Clears the Vector of all Item, except those matching the internal "vector"
// key e.g. "vector.tag", "vector.id", etc et al.
func (v *Vector) Reset() {
	ci := v.Match("vector")
	v.reset()
	v.Set(ci...)
}

// Returns the Vector data as a map[string]interface{} suitable for use with
// text.Template or html.Template. Keys are undotted form(e.g. key.key becomes
// KeyKey).
func (v *Vector) TemplateData() map[string]interface{} {
	ret := make(map[string]interface{})
	l := v.List()
	for _, i := range l {
		ret[i.KeyUndotted()] = i.Provided()
	}
	return ret
}

// json.Marshaler
func (v *Vector) MarshalJSON() ([]byte, error) {
	l := v.List()
	return json.Marshal(&l)
}

// json.Unmarshaler
func (v *Vector) UnmarshalJSON(b []byte) error {
	v.ensureNotEmpty()
	var i []*Mtem
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	var ii []Item
	for _, v := range i {
		ii = append(ii, fromMtem(v))
	}
	v.Set(ii...)
	return nil
}

// yaml.Marshaler
func (v *Vector) MarshalYAML() (interface{}, error) {
	return v.List(), nil
}

// yaml.Unmarshaler
func (v *Vector) UnmarshalYAML(u func(interface{}) error) error {
	v.ensureNotEmpty()
	var i []*Mtem
	err := u(&i)
	if err != nil {
		return err
	}
	var ii []Item
	for _, v := range i {
		ii = append(ii, fromMtem(v))
	}
	v.Set(ii...)
	return nil
}

// Return a string from key matching a stored StringItem.
func (v *Vector) ToString(k string) string {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(StringItem); ok {
			return ii.ToString()
		}
	}
	return ""
}

// Set a StringItem with the provided key and value.
func (v *Vector) SetString(k, vi string) {
	ni := NewStringItem(k, vi)
	v.Set(ni)
}

// Return an array of strings from a matching key.
// Storing a StringsItem is relatively faster, but will attempt to return strings
// from a StringItem.
func (v *Vector) ToStrings(k string) []string {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(StringsItem); ok {
			return ii.ToStrings()
		}
		//
	}
	return []string{}
}

// Set a StringsItem with the provided key and string values.
func (v *Vector) SetStrings(k string, vi ...string) {
	ni := NewStringsItem(k, vi...)
	v.Set(ni)
}

// Return a boolean from a matching key.
// Storing a BoolItem is relatively faster, but will attempt to return a bool
// from a StringItem.
func (v *Vector) ToBool(k string) bool {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(BoolItem); ok {
			return ii.ToBool()
		}
		if iis, ok := i.(StringItem); ok {
			v := iis.ToString()
			if b, err := strconv.ParseBool(v); err == nil {
				return b
			}
		}
	}
	return false
}

// Set a BoolItem with the provided key and boolean value.
func (v *Vector) SetBool(k string, vi bool) {
	ni := NewBoolItem(k, vi)
	v.Set(ni)
}

// Return an integer from a matching key.
// Storing an IntItem is relatively faster, but will attempt to return an int
// from a StringItem.
func (v *Vector) ToInt(k string) int {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(IntItem); ok {
			return ii.ToInt()
		}
		if iii, ok := i.(StringItem); ok {
			v := iii.ToString()
			if ri, err := strconv.ParseInt(v, 10, 64); err == nil {
				return int(ri)
			}
		}
	}
	return 0
}

// Set an IntItem with the provided key and integer value.
func (v *Vector) SetInt(k string, vi int) {
	ni := NewIntItem(k, vi)
	v.Set(ni)
}

// Return an int64 from a matching key.
// Storing an Int64Item is relatively faster, but will attempt to return an int64
// from a StringItem.
func (v *Vector) ToInt64(k string) int64 {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(Int64Item); ok {
			return ii.ToInt64()
		}
		if iii, ok := i.(StringItem); ok {
			v := iii.ToString()
			if r, err := strconv.ParseInt(v, 10, 64); err == nil {
				return r
			}
		}
	}
	return 0
}

// Set an Int64Item with the provided key and int64 value.
func (v *Vector) SetInt64(k string, vi int64) {
	ni := NewInt64Item(k, vi)
	v.Set(ni)
}

// Return an uint from a key matching a stored UintItem.
// Storing a UintItem is relatively faster, but will attempt to return a uint
// from a StringItem.
func (v *Vector) ToUint(k string) uint {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(UintItem); ok {
			return ii.ToUint()
		}
		if iii, ok := i.(StringItem); ok {
			v := iii.ToString()
			if r, err := strconv.ParseUint(v, 10, 64); err == nil {
				return uint(r)
			}
		}
	}
	return 0
}

// Set an UintItem with the provided key and uint value.
func (v *Vector) SetUint(k string, vi uint) {
	ni := NewUintItem(k, vi)
	v.Set(ni)
}

// Return an uint64 from a key matching a stored Uint64Item.
// Storing a Uint64Item is relatively faster, but will attempt to return a uint64
// from a StringItem.
func (v *Vector) ToUint64(k string) uint64 {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(Uint64Item); ok {
			return ii.ToUint64()
		}
		if iii, ok := i.(StringItem); ok {
			v := iii.ToString()
			if r, err := strconv.ParseUint(v, 10, 64); err == nil {
				return r
			}
		}
	}
	return 0
}

// Set an Uint64Item with the provided key and uint64 value.
func (v *Vector) SetUint64(k string, vi uint64) {
	ni := NewUint64Item(k, vi)
	v.Set(ni)
}

// Return a float64 from a key matching a stored Float64Item.
// Storing a float64Item is relatively faster, but will attempt to return a float64
// from a StringItem.
func (v *Vector) ToFloat64(k string) float64 {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(Float64Item); ok {
			return ii.ToFloat64()
		}
		if iii, ok := i.(StringItem); ok {
			v := iii.ToString()
			if rf, err := strconv.ParseFloat(v, 64); err == nil {
				return rf
			}
		}
	}
	return 0
}

/// Set a Float64Item with the provided key and float64 value.
func (v *Vector) SetFloat64(k string, vi float64) {
	ni := NewFloat64Item(k, vi)
	v.Set(ni)
}

// Return a *Vector from a key matching a stored VectorItem.
func (v *Vector) ToVector(k string) *Vector {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(VectorItem); ok {
			return ii.ToVector()
		}
	}
	return nil
}

// Set a VectorItem with the provided key and *Vector value.
func (v *Vector) SetVector(k string, vi *Vector) {
	ni := NewVectorItem(k, vi)
	v.Set(ni)
}
