package data

import "encoding/json"

type Container struct {
	Tag   string
	Items map[string]*Item
}

func NewContainer(tag string) *Container {
	return &Container{
		Tag:   tag,
		Items: make(map[string]*Item),
	}
}

func (c *Container) Tagged() string {
	return c.Tag
}

func (c *Container) Set(i *Item) {
	c.Items[i.Key] = i
}

func (c *Container) SetItem(is ...*Item) {
	for _, v := range is {
		c.Set(v)
	}
}

func (c *Container) Get(k string) *Item {
	if i, exists := c.Items[k]; exists {
		return i
	}
	return nil
}

func (c *Container) Clone(except ...string) *Container {
	n := NewContainer(c.Tagged())
	l := c.List()
	n.SetItem(l...)
	return n
}

func match(ss []string, s string) bool {
	for _, v := range ss {
		if v == s {
			return true
		}
	}
	return false
}

func (c *Container) List(except ...string) []*Item {
	var ret []*Item
	for k, v := range c.Items {
		if !match(except, k) {
			ret = append(ret, v)
		}
	}
	return ret
}

func (c *Container) ToString(k string) string {
	if i, exists := c.Items[k]; exists {
		return i.ToString()
	}
	return ""
}

func (c *Container) ToBool(k string) bool {
	if i, exists := c.Items[k]; exists {
		return i.ToBool()
	}
	return false
}

func (c *Container) ToInt(k string) int {
	if i, exists := c.Items[k]; exists {
		return i.ToInt()
	}
	return 0
}

func (c *Container) ToFloat(k string) float64 {
	if i, exists := c.Items[k]; exists {
		return i.ToFloat()
	}
	return 0
}

func (c *Container) ToList(k string) []string {
	if i, exists := c.Items[k]; exists {
		return i.ToList()
	}
	return []string{}
}

func (c *Container) ToMap(k string) map[string]string {
	if i, exists := c.Items[k]; exists {
		return i.ToMap()
	}
	return nil
}

func (c *Container) String() string {
	j, err := c.MarshalJSON()
	if err != nil {
		return err.Error()
	}
	return string(j)
}

func (c *Container) MarshalJSON() ([]byte, error) {
	l := c.List()
	return json.Marshal(&l)
}

func (c *Container) UnmarshalJSON(b []byte) error {
	var i []*Item
	err := json.Unmarshal(b, &i)
	if err != nil {
		return err
	}
	for _, v := range i {
		c.Set(v)
	}
	return nil
}

func (c *Container) MarshalYAML() (interface{}, error) {
	return c.List(), nil
}

func (c *Container) UnmarshalYAML(u func(interface{}) error) error {
	var i []*Item
	err := u(&i)
	if err != nil {
		return err
	}
	for _, v := range i {
		c.Set(v)
	}
	return nil
}
