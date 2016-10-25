package data

import (
	"encoding/json"
	"fmt"
	"strconv"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Keyer interface {
	Key() string
	Undotted() string
	Change(string)
}

type Valuer interface {
	StringItem
	BoolItem
	IntItem
	FloatItem
	ListItem
	MapItem
}

type StringItem interface {
	ToString() string
	SetString(string)
}

//type ByteItem interface {
//	ToByte() []byte
//	SetByte([]byte)
//}

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

type KVer interface {
	ToKVString() string
	ToKV() (string, interface{})
}

type Transmitter interface {
	BasicItem() *BasicItem
	JsonTransmitter
	YamlTransmitter
}

type JsonTransmitter interface {
	json.Marshaler
	json.Unmarshaler
}

type YamlTransmitter interface {
	yaml.Marshaler
	yaml.Unmarshaler
}

type Item interface {
	Keyer
	Valuer
	KVer
	Transmitter
	Clone(string) Item
}

type item struct {
	key, Value string
}

func NewItem(k, v string) Item {
	return &item{k, v}
}

func (i *item) Key() string {
	return i.key
}

func (i *item) Undotted() string {
	k := strings.Split(i.key, ".")
	var j []string
	for _, kv := range k {
		j = append(j, strings.Title(kv))
	}
	return strings.Join(j, "")
}

func (i *item) Change(to string) {
	i.key = to
}

func (i *item) ToKVString() string {
	return fmt.Sprintf("%s:%s", i.key, i.Value)
}

func (i *item) ToKV() (string, interface{}) {
	return i.key, i.Value
}

func (i *item) ToString() string {
	return i.Value
}

func (i *item) SetString(s string) {
	i.Value = s
}

func (i *item) ToBool() bool {
	if vl, err := strconv.ParseBool(i.Value); err == nil {
		return vl
	}
	return false
}

func (i *item) SetBool(b bool) {
	i.Value = strconv.FormatBool(b)
}

func (i *item) ToInt() int {
	if vl, err := strconv.Atoi(i.Value); err == nil {
		return vl
	}
	return 0
}

func (i *item) SetInt(in int) {
	i.Value = strconv.Itoa(in)
}

func (i *item) ToFloat() float64 {
	if vl, err := strconv.ParseFloat(i.Value, 64); err == nil {
		return vl
	}
	return 0.0
}

func (i *item) SetFloat(in float64) {
	i.Value = strconv.FormatFloat(in, 'f', 1, 64)
}

func (i *item) ToList() []string {
	vl := i.Value
	spl := strings.Split(vl, ",")
	return spl
}

func (i *item) SetList(l ...string) {
	list := strings.Join(l, ",")
	i.Value = list
}

func (i *item) ToMap() map[string]string {
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

func (i *item) SetMap(m map[string]string) {
	var set []string
	for k, v := range m {
		set = append(set, fmt.Sprintf("%s:%s", k, v))
	}
	i.Value = strings.Join(set, ",")
}

type BasicItem struct {
	Key, Value string
}

func (i *item) BasicItem() *BasicItem {
	return &BasicItem{i.key, i.Value}
}

func (i *item) MarshalJSON() ([]byte, error) {
	return json.Marshal(i.BasicItem())
}

func (i *item) UnmarshalJSON(b []byte) error {
	var bi *BasicItem
	err := json.Unmarshal(b, &bi)
	if err != nil {
		return err
	}
	i.key = bi.Key
	i.Value = bi.Value
	return nil
}

func (i *item) MarshalYAML() (interface{}, error) {
	return i.BasicItem(), nil
}

func (i *item) UnmarshalYAML(u func(interface{}) error) error {
	var bi *BasicItem
	err := u(&bi)
	if err != nil {
		return err
	}
	i.key = bi.Key
	i.Value = bi.Value
	return nil
}

func (i *item) Clone(k string) Item {
	var key, value string
	if k != "" {
		key = k
	} else {
		key = i.Key()
	}
	value = i.Value
	return NewItem(key, value)
}
