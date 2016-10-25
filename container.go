package data

import (
	"bytes"
	"encoding/json"
	"fmt"
)

type Container struct {
	*Trie
}

func New(tag string, o ...Option) *Container {
	t := NewTrie(o...)
	if tag != "" {
		t.set(
			NewItem("container.tag", tag),
			NewItem("container.id", V4Quick()),
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
		i.SetString(t)
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
		nl = append(nl, v.Clone(""))
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
		ret[v.Undotted()] = v.ToString()
	}
	return ret
}

func (c *Container) MarshalJSON() ([]byte, error) {
	l := c.List()
	return json.Marshal(&l)
}

func (c *Container) UnmarshalJSON(b []byte) error {
	var bi []*BasicItem
	err := json.Unmarshal(b, &bi)
	if err != nil {
		return err
	}
	var i []Item
	for _, v := range bi {
		i = append(i, NewItem(v.Key, v.Value))
	}
	c.Set(i...)
	return nil
}

func (c *Container) MarshalYAML() (interface{}, error) {
	return c.List(), nil
}

func (c *Container) UnmarshalYAML(u func(interface{}) error) error {
	var bi []*BasicItem
	err := u(&bi)
	if err != nil {
		return err
	}
	var i []Item
	for _, v := range bi {
		i = append(i, NewItem(v.Key, v.Value))
	}
	c.Set(i...)
	return nil
}

func (c *Container) ToString(k string) string {
	if i := c.Get(k); i != nil {
		return i.ToString()
	}
	return ""
}

func (c *Container) ToKVString() string {
	return fmt.Sprintf("%s:%s", c.Tag(), c.Keys())
}

func (c *Container) ToKV() (string, interface{}) {
	return c.Tag(), c.Keys()
}

func (c *Container) ToBool(k string) bool {
	if i := c.Get(k); i != nil {
		return i.ToBool()
	}
	return false
}

func (c *Container) ToInt(k string) int {
	if i := c.Get(k); i != nil {
		return i.ToInt()
	}
	return 0
}

func (c *Container) ToFloat(k string) float64 {
	if i := c.Get(k); i != nil {
		return i.ToFloat()
	}
	return 0
}

func (c *Container) ToList(k string) []string {
	if i := c.Get(k); i != nil {
		return i.ToList()
	}
	return []string{}
}

func (c *Container) ToMap(k string) map[string]string {
	if i := c.Get(k); i != nil {
		return i.ToMap()
	}
	return nil
}
