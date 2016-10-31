package data

import (
	"encoding/json"
	"math"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

type Keyer interface {
	Key() string
	KeyUndotted() string
	NewKey(string)
}

type Valuer interface {
	Value() []byte
	Provided() interface{}
	Provide(interface{})
}

type Transmitter interface {
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

type Cloner interface {
	Clone() Item
}

type Item interface {
	Keyer
	Valuer
	Transmitter
	Cloner
}

type item struct {
	key      string
	provided interface{}
	value    []byte
}

func (i *item) Key() string {
	return i.key
}

func (i *item) KeyUndotted() string {
	k := strings.Split(i.key, ".")
	var j []string
	for _, kv := range k {
		j = append(j, strings.Title(kv))
	}
	return strings.Join(j, "")
}

func (i *item) NewKey(k string) {
	i.key = k
}

func KeyedItem(k string) Item {
	return &item{k, nil, nil}
}

func (i *item) Value() []byte {
	if i.value == nil {
		b, err := json.Marshal(i.provided)
		if err != nil {
			r := []byte(err.Error())
			i.value = r
			return r
		}
		i.value = b
	}
	return i.value
}

func (i *item) Provided() interface{} {
	return i.provided
}

func (i *item) Provide(p interface{}) {
	i.provided = p
	i.value = nil
}

type Mtem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func fromFloat(i *item) Item {
	v, ok := i.provided.(float64)
	if ok {
		if math.Mod(v, 1) == 0 {
			return &intItem{i}
		}
	}
	return &floatItem{i}
}

func fromArrayedInterface(i *item) Item {
	v := i.Value()

	var s []string
	if err := json.Unmarshal(v, &s); err == nil {
		return &stringsItem{
			&item{
				key:      i.key,
				provided: s,
			},
		}
	}

	c := New("")
	if err := json.Unmarshal(v, &c); err == nil {
		return &multiItem{
			&item{
				key:      i.key,
				provided: c,
			},
		}

	}

	return i
}

func fromMtem(m *Mtem) Item {
	i := &item{
		key:      m.Key,
		provided: m.Value,
	}

	switch m.Value.(type) {
	case string:
		return &stringItem{i}
	case bool:
		return &boolItem{i}
	case float64:
		return fromFloat(i)
	case []interface{}:
		return fromArrayedInterface(i)
	}

	return i
}

func (i *item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&Mtem{i.key, i.provided})
}

func (i *item) UnmarshalJSON(b []byte) error {
	var m Mtem
	err := json.Unmarshal(b, &m)
	if err != nil {
		return err
	}
	i.key = m.Key
	i.provided = m.Value
	return nil
}

func (i *item) MarshalYAML() (interface{}, error) {
	return &Mtem{i.key, i.provided}, nil
}

func (i *item) UnmarshalYAML(u func(interface{}) error) error {
	var m Mtem
	err := u(&m)
	if err != nil {
		return err
	}
	i.key = m.Key
	i.provided = m.Value
	return nil
}

func (i *item) Clone() Item {
	ni := *i
	ni.value = nil
	return &ni
}

type StringItem interface {
	Item
	ToString() string
	SetString(string)
}

type stringItem struct {
	Item
}

func NewStringItem(key, v string) StringItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &stringItem{i}
}

func (i *stringItem) ToString() string {
	var ret string
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return err.Error()
	}
	return ret
}

func (i *stringItem) SetString(s string) {
	i.Provide(s)
}

func (i *stringItem) Clone() Item {
	ii := i.Item.Clone()
	return &stringItem{ii}
}

type StringsItem interface {
	Item
	ToStrings() []string
	SetStrings(...string)
}

type stringsItem struct {
	Item
}

func NewStringsItem(key string, v ...string) StringsItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &stringsItem{i}
}

func (i *stringsItem) ToStrings() []string {
	var ret []string
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return []string{err.Error()}
	}
	return ret
}

func (i *stringsItem) SetStrings(l ...string) {
	i.Provide(l)
}

func (i *stringsItem) Clone() Item {
	ii := i.Item.Clone()
	return &stringsItem{ii}
}

type BoolItem interface {
	Item
	ToBool() bool
	SetBool(bool)
}

type boolItem struct {
	Item
}

func NewBoolItem(key string, v bool) BoolItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &boolItem{i}
}

func (i *boolItem) ToBool() bool {
	var ret bool
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return false
	}
	return ret
}

func (i *boolItem) SetBool(v bool) {
	i.Provide(v)
}

func (i *boolItem) Clone() Item {
	ii := i.Item.Clone()
	return &boolItem{ii}
}

type IntItem interface {
	Item
	ToInt() int
	SetInt(int)
}

type intItem struct {
	Item
}

func NewIntItem(key string, v int) IntItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &intItem{i}
}

func (i *intItem) ToInt() int {
	var ret int
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return 0
	}
	return ret
}

func (i *intItem) SetInt(v int) {
	i.Provide(v)
}

func (i *intItem) Clone() Item {
	ii := i.Item.Clone()
	return &intItem{ii}
}

type FloatItem interface {
	Item
	ToFloat() float64
	SetFloat(float64)
}

type floatItem struct {
	Item
}

func NewFloatItem(key string, v float64) FloatItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &floatItem{i}
}

func (i *floatItem) ToFloat() float64 {
	var ret float64
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return 0
	}
	return ret
}

func (i *floatItem) SetFloat(v float64) {
	i.Provide(v)
}

func (i *floatItem) Clone() Item {
	ii := i.Item.Clone()
	return &floatItem{ii}
}

type MultiItem interface {
	Item
	ToMulti() *Container
	SetMulti(*Container)
}

type multiItem struct {
	Item
}

func NewMultiItem(key string, v *Container) Item {
	i := KeyedItem(key)
	i.Provide(v)
	return &multiItem{i}
}

func (i *multiItem) ToMulti() *Container {
	ret := New("")
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return nil
	}
	return ret
}

func (i *multiItem) SetMulti(v *Container) {
	i.Provide(v)
}

func (i *multiItem) Clone() Item {
	c := i.ToMulti()
	ii := i.Item.Clone()
	ii.Provide(c)
	return &multiItem{ii}
}
