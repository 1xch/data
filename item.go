package data

import (
	"fmt"
	"strconv"
	"strings"
)

type Item struct {
	Key, Value string
}

type StringItem interface {
	ToString() string
	SetString(string)
}

type BoolItem interface {
	ToBool() bool
	SetBool(bool)
}

type IntItem interface {
	ToInt() int
	SetInt(int)
}

type FloatItem interface {
	ToFloat() float64
	SetFloat(float64)
}

type ListItem interface {
	ToList() []string
	SetList(...string)
}

type MapItem interface {
	ToMap() map[string]string
	SetMap(map[string]string)
}

func NewItem(k, v string) *Item {
	return &Item{k, v}
}

func (i *Item) ToString() string {
	return i.Value
}

func (i *Item) SetString(s string) {
	i.Value = s
}

func (i *Item) ToBool() bool {
	if vl, err := strconv.ParseBool(i.Value); err == nil {
		return vl
	}
	return false
}

func (i *Item) SetBool(b bool) {
	i.Value = strconv.FormatBool(b)
}

func (i *Item) ToInt() int {
	if vl, err := strconv.Atoi(i.Value); err == nil {
		return vl
	}
	return 0
}

func (i *Item) SetInt(in int) {
	i.Value = strconv.Itoa(in)
}

func (i *Item) ToFloat() float64 {
	if vl, err := strconv.ParseFloat(i.Value, 64); err == nil {
		return vl
	}
	return 0.0
}

func (i *Item) SetFloat(in float64) {
	i.Value = strconv.FormatFloat(in, 'f', 1, 64)
}

func (i *Item) ToList() []string {
	vl := i.Value
	spl := strings.Split(vl, ",")
	return spl
}

func (i *Item) SetList(l ...string) {
	list := strings.Join(l, ",")
	i.Value = list
}

func (i *Item) ToMap() map[string]string {
	ret := make(map[string]string)
	list := i.ToList()
	for _, v := range list {
		spl := strings.Split(v, ":")
		if len(spl) == 2 {
			ret[spl[0]] = spl[1]
		} else {
			ret[spl[0]] = "is not mappable"
		}
	}
	return ret
}

func (i *Item) SetMap(m map[string]string) {
	var set []string
	for k, v := range m {
		set = append(set, fmt.Sprintf("%s:%s", k, v))
	}
	i.Value = strings.Join(set, ",")
}
