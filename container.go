package data

import (
	"bytes"
	"encoding/json"
)

type Container struct {
	*Trie
}

func New(tag string, o ...Option) *Container {
	t := NewTrie(o...)
	if tag != "" {
		t.set(
			NewStringItem("container.tag", tag),
			NewStringItem("container.id", V4Quick()),
		)
	}
	return &Container{
		Trie: t,
	}
}

func (c *Container) Tag() string {
	return c.ToString("container.tag")
}

func (c *Container) Retag(t string) {
	i := c.Get("container.tag")
	if i != nil {
		i.Provide(t)
	}
}

func (c *Container) Keys() []string {
	var ret []string
	v := func(p Prefix, i Item) error {
		ret = append(ret, string(p))
		return nil
	}
	c.walk(nil, v)
	return ret
}

func (c *Container) Get(k string) Item {
	key := Prefix(k)
	return c.get(key)
}

func (c *Container) Match(k string) []Item {
	var ret []Item
	bk := []byte(k)
	v := func(p Prefix, i Item) error {
		if bytes.Contains(p, bk) {
			ret = append(ret, i)
		}
		return nil
	}
	c.walk(nil, v)
	return ret
}

func (c *Container) Set(i ...Item) {
	c.set(i...)
}

func (c *Container) Merge(cs ...*Container) {
	for _, cc := range cs {
		l := cc.List("container.tag", "container.id")
		c.Set(l...)
	}
}

func (c *Container) Clone(except ...string) *Container {
	except = append(except, "container.tag", "container.id")
	n := New(c.Tag())
	l := c.List(except...)
	var nl []Item
	for _, v := range l {
		nl = append(nl, v.Clone())
	}
	n.Set(nl...)
	return n
}

func (c *Container) CloneAs(tag string, except ...string) *Container {
	nc := c.Clone(except...)
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

func (c *Container) List(except ...string) []Item {
	var ret []Item
	v := func(p Prefix, i Item) error {
		if !match(except, i.Key()) {
			ret = append(ret, i)
		}
		return nil
	}
	c.walk(nil, v)
	return ret
}

func (c *Container) Clear() {
	c.reset()
}

func (c *Container) Reset() {
	ci := c.Match("container")
	c.reset()
	c.Set(ci...)
}

func (c *Container) TemplateData() map[string]interface{} {
	ret := make(map[string]interface{})
	l := c.List()
	for _, v := range l {
		ret[v.KeyUndotted()] = v.Provided()
	}
	return ret
}

func (c *Container) MarshalJSON() ([]byte, error) {
	l := c.List()
	return json.Marshal(&l)
}

func (c *Container) UnmarshalJSON(b []byte) error {
	var i []*Mtem
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	var ii []Item
	for _, v := range i {
		ii = append(ii, fromMtem(v))
	}
	c.Set(ii...)
	return nil
}

func (c *Container) MarshalYAML() (interface{}, error) {
	return c.List(), nil
}

func (c *Container) UnmarshalYAML(u func(interface{}) error) error {
	var i []*Mtem
	err := u(&i)
	if err != nil {
		return err
	}
	var ii []Item
	for _, v := range i {
		ii = append(ii, fromMtem(v))
	}
	c.Set(ii...)
	return nil
}

func (c *Container) ToString(k string) string {
	if i := c.Get(k); i != nil {
		if ii, ok := i.(StringItem); ok {
			return ii.ToString()
		}
	}
	return ""
}

func (c *Container) ToStrings(k string) []string {
	if i := c.Get(k); i != nil {
		if ii, ok := i.(StringsItem); ok {
			return ii.ToStrings()
		}
	}
	return []string{}
}

func (c *Container) ToBool(k string) bool {
	if i := c.Get(k); i != nil {
		if ii, ok := i.(BoolItem); ok {
			return ii.ToBool()
		}
	}
	return false
}

func (c *Container) ToInt(k string) int {
	if i := c.Get(k); i != nil {
		if ii, ok := i.(IntItem); ok {
			return ii.ToInt()
		}
	}
	return 0
}

func (c *Container) ToFloat(k string) float64 {
	if i := c.Get(k); i != nil {
		if ii, ok := i.(FloatItem); ok {
			return ii.ToFloat()
		}
	}
	return 0
}

func (c *Container) ToMulti(k string) *Container {
	if i := c.Get(k); i != nil {
		if ii, ok := i.(MultiItem); ok {
			return ii.ToMulti()
		}
	}
	return nil
}
