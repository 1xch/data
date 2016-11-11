package data

import (
	"bytes"
	"encoding/json"
	"sync"
)

type Vector struct {
	l *sync.RWMutex
	*Trie
}

func New(tag string, o ...Option) *Vector {
	l := &sync.RWMutex{}
	t := NewTrie(o...)
	v := &Vector{
		l, t,
	}
	if tag != "" {
		v.Set(
			NewStringItem("vector.tag", tag),
			NewStringItem("vector.id", V4Quick()),
		)
	}
	return v
}

func (v *Vector) Tag() string {
	return v.ToString("vector.tag")
}

func (v *Vector) Retag(t string) {
	i := v.Get("vector.tag")
	if i != nil {
		i.Provide(t)
	}
}

func (v *Vector) Keys() []string {
	var ret []string
	w := func(p Prefix, i Item) error {
		ret = append(ret, string(p))
		return nil
	}
	v.walk(nil, w)
	return ret
}

func (v *Vector) Get(k string) Item {
	v.l.RLock()
	key := Prefix(k)
	i := v.get(key)
	v.l.RUnlock()
	return i
}

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

func (v *Vector) Set(i ...Item) {
	v.l.Lock()
	v.set(i...)
	v.l.Unlock()
}

func (v *Vector) Merge(vs ...*Vector) {
	for _, vv := range vs {
		l := vv.List("vector.tag", "vector.id")
		v.Set(l...)
	}
}

func (v *Vector) Clone(except ...string) *Vector {
	except = append(except, "vector.tag", "vector.id")
	n := New(v.Tag())
	l := v.List(except...)
	var nl []Item
	for _, i := range l {
		nl = append(nl, i.Clone())
	}
	n.Set(nl...)
	return n
}

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

func (v *Vector) Clear() {
	v.reset()
}

func (v *Vector) Reset() {
	ci := v.Match("vector")
	v.reset()
	v.Set(ci...)
}

func (v *Vector) TemplateData() map[string]interface{} {
	ret := make(map[string]interface{})
	l := v.List()
	for _, i := range l {
		ret[i.KeyUndotted()] = i.Provided()
	}
	return ret
}

func (v *Vector) MarshalJSON() ([]byte, error) {
	l := v.List()
	return json.Marshal(&l)
}

func (v *Vector) UnmarshalJSON(b []byte) error {
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

func (v *Vector) MarshalYAML() (interface{}, error) {
	return v.List(), nil
}

func (v *Vector) UnmarshalYAML(u func(interface{}) error) error {
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

func (v *Vector) ToString(k string) string {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(StringItem); ok {
			return ii.ToString()
		}
	}
	return ""
}

func (v *Vector) SetString(k, vi string) {
	ni := NewStringItem(k, vi)
	v.Set(ni)
}

func (v *Vector) ToStrings(k string) []string {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(StringsItem); ok {
			return ii.ToStrings()
		}
	}
	return []string{}
}

func (v *Vector) SetStrings(k string, vi ...string) {
	ni := NewStringsItem(k, vi...)
	v.Set(ni)
}

func (v *Vector) ToBool(k string) bool {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(BoolItem); ok {
			return ii.ToBool()
		}
	}
	return false
}

func (v *Vector) SetBool(k string, vi bool) {
	ni := NewBoolItem(k, vi)
	v.Set(ni)
}

func (v *Vector) ToInt(k string) int {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(IntItem); ok {
			return ii.ToInt()
		}
	}
	return 0
}

func (v *Vector) SetInt(k string, vi int) {
	ni := NewIntItem(k, vi)
	v.Set(ni)
}

func (v *Vector) ToFloat(k string) float64 {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(FloatItem); ok {
			return ii.ToFloat()
		}
	}
	return 0
}

func (v *Vector) SetFloat(k string, vi float64) {
	ni := NewFloatItem(k, vi)
	v.Set(ni)
}

func (v *Vector) ToVector(k string) *Vector {
	if i := v.Get(k); i != nil {
		if ii, ok := i.(VectorItem); ok {
			return ii.ToVector()
		}
	}
	return nil
}

func (v *Vector) SetVector(k string, vi *Vector) {
	ni := NewVectorItem(k, vi)
	v.Set(ni)
}
