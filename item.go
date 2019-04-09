package data

import (
	"encoding/json"
	"strings"

	yaml "gopkg.in/yaml.v2"
)

// A string key management interface.
type Keyer interface {
	Key() string
	KeyUndotted() string
	NewKey(string)
}

// An interface for managing any type of value.
type Valuer interface {
	Value() []byte
	Provided() interface{}
	Provide(interface{})
}

// An interface for managing transmission of values between formats
type Transmitter interface {
	JsonTransmitter
	YamlTransmitter
}

// A json transmitter interface.
type JsonTransmitter interface {
	json.Marshaler
	json.Unmarshaler
}

// A yaml transmitter interface.
type YamlTransmitter interface {
	yaml.Marshaler
	yaml.Unmarshaler
}

// An interface for encapsulating item cloning.
type Cloner interface {
	Clone() Item
}

// An interface for storing and transmitting single items composed of Keyer,
// Valuer, Transmitter, and Cloner interfaces.
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

// Returns the item's string key.
func (i *item) Key() string {
	return i.key
}

// Returns the item's string key with any dots removed:
// key.key.key becomes KeyKeyKey.
func (i *item) KeyUndotted() string {
	k := strings.Split(i.key, ".")
	var j []string
	for _, kv := range k {
		j = append(j, strings.Title(kv))
	}
	return strings.Join(j, "")
}

// Sets a new key for the item.
func (i *item) NewKey(k string) {
	i.key = k
}

// Returns an empty item with the provided key.
func KeyedItem(k string) Item {
	return &item{k, nil, nil}
}

// Returns this item's value as a []byte.
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

// Returns what this item was provided with as an interface{}
func (i *item) Provided() interface{} {
	return i.provided
}

// Provides this item with the given interface{}
func (i *item) Provide(p interface{}) {
	i.provided = p
	i.value = nil
}

// An intermediary unmarshaling type.
type Mtem struct {
	Key   string      `json:"key"`
	Value interface{} `json:"value"`
}

func fromVector(i *item) Item {
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

	vec := New("")
	if err := json.Unmarshal(v, &vec); err == nil {
		return &vectorItem{
			&item{
				key:      i.key,
				provided: vec,
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
	case int:
		return &intItem{i}
	case int64:
		return &int64Item{i}
	case uint:
		return &uintItem{i}
	case uint64:
		return &uint64Item{i}
	case float64:
		return &float64Item{i}
	case []interface{}:
		return fromVector(i)
	}

	return i
}

// json.Marshaler for this item.
func (i *item) MarshalJSON() ([]byte, error) {
	return json.Marshal(&Mtem{i.key, i.provided})
}

// json.Unmarshaler for this item.
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

// yaml.Marshaler for this item.
func (i *item) MarshalYAML() (interface{}, error) {
	return &Mtem{i.key, i.provided}, nil
}

// yaml.Unmarshaler for this item.
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

// Clones this item as a separate copy.
func (i *item) Clone() Item {
	ni := *i
	ni.value = nil
	return &ni
}

// An interface for a specific string type Item.
type StringItem interface {
	Item
	ToString() string
	SetString(string)
}

type stringItem struct {
	Item
}

// Creates a new StringItem from the provided key and value.
func NewStringItem(key, v string) StringItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &stringItem{i}
}

// Returns the string value of this StringItem.
func (i *stringItem) ToString() string {
	var ret string
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return err.Error()
	}
	return ret
}

// Sets a string value for this StringItem.
func (i *stringItem) SetString(s string) {
	i.Provide(s)
}

// Satisfies the Cloner interface for this StringItem.
func (i *stringItem) Clone() Item {
	ii := i.Item.Clone()
	return &stringItem{ii}
}

// An interface for a specific []string type Item.
type StringsItem interface {
	Item
	ToStrings() []string
	SetStrings(...string)
}

type stringsItem struct {
	Item
}

// Creates a new StringsItem from the provided key and string values.
func NewStringsItem(key string, v ...string) StringsItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &stringsItem{i}
}

//
func (i *stringsItem) ToStrings() []string {
	var ret []string
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return []string{err.Error()}
	}
	return ret
}

//
func (i *stringsItem) SetStrings(l ...string) {
	i.Provide(l)
}

//
func (i *stringsItem) Clone() Item {
	ii := i.Item.Clone()
	return &stringsItem{ii}
}

// An interface for a specific bool type Item.
type BoolItem interface {
	Item
	ToBool() bool
	SetBool(bool)
}

type boolItem struct {
	Item
}

// Creates a new BoolItem from the provided string key and boolean value.
func NewBoolItem(key string, v bool) BoolItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &boolItem{i}
}

//
func (i *boolItem) ToBool() bool {
	var ret bool
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return false
	}
	return ret
}

//
func (i *boolItem) SetBool(v bool) {
	i.Provide(v)
}

//
func (i *boolItem) Clone() Item {
	ii := i.Item.Clone()
	return &boolItem{ii}
}

// An interface for a specific int type Item.
type IntItem interface {
	Item
	ToInt() int
	SetInt(int)
}

type intItem struct {
	Item
}

// Creates a new IntItem from the provided string key and int value.
func NewIntItem(key string, v int) IntItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &intItem{i}
}

//
func (i *intItem) ToInt() int {
	var ret int
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return 0
	}
	return ret
}

//
func (i *intItem) SetInt(v int) {
	i.Provide(v)
}

//
func (i *intItem) Clone() Item {
	ii := i.Item.Clone()
	return &intItem{ii}
}

// An interface for a specific int64 type Item.
type Int64Item interface {
	Item
	ToInt64() int64
	SetInt64(int)
}

type int64Item struct {
	Item
}

// Creates a new Int64Item from the provided string key and int64 value.
func NewInt64Item(key string, v int64) Int64Item {
	i := KeyedItem(key)
	i.Provide(v)
	return &int64Item{i}
}

//
func (i *int64Item) ToInt64() int64 {
	var ret int64
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return 0
	}
	return ret
}

//
func (i *int64Item) SetInt64(v int) {
	i.Provide(v)
}

//
func (i *int64Item) Clone() Item {
	ii := i.Item.Clone()
	return &int64Item{ii}
}

// An interface for a specific uint type Item.
type UintItem interface {
	Item
	ToUint() uint
	SetUint(uint)
}

type uintItem struct {
	Item
}

// Creates a new UintItem from the provided string key and uint value.
func NewUintItem(key string, v uint) UintItem {
	i := KeyedItem(key)
	i.Provide(v)
	return &uintItem{i}
}

//
func (i *uintItem) ToUint() uint {
	var ret uint
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return 0
	}
	return ret
}

//
func (i *uintItem) SetUint(v uint) {
	i.Provide(v)
}

//
func (i *uintItem) Clone() Item {
	ii := i.Item.Clone()
	return &uintItem{ii}
}

// An interface for a specific uint64 type Item.
type Uint64Item interface {
	Item
	ToUint64() uint64
	SetUint64(uint64)
}

type uint64Item struct {
	Item
}

// Creates a new Uint64Item from the provided string key and uint64 value.
func NewUint64Item(key string, v uint64) Uint64Item {
	i := KeyedItem(key)
	i.Provide(v)
	return &uint64Item{i}
}

//
func (i *uint64Item) ToUint64() uint64 {
	var ret uint64
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return 0
	}
	return ret
}

//
func (i *uint64Item) SetUint64(v uint64) {
	i.Provide(v)
}

//
func (i *uint64Item) Clone() Item {
	ii := i.Item.Clone()
	return &uint64Item{ii}
}

// An interface for a specific float64 type Item.
type Float64Item interface {
	Item
	ToFloat64() float64
	SetFloat(float64)
}

type float64Item struct {
	Item
}

// Creates a new Float64Item from the provided string key and float64 value.
func NewFloat64Item(key string, v float64) Float64Item {
	i := KeyedItem(key)
	i.Provide(v)
	return &float64Item{i}
}

//
func (i *float64Item) ToFloat64() float64 {
	var ret float64
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return 0
	}
	return ret
}

//
func (i *float64Item) SetFloat(v float64) {
	i.Provide(v)
}

//
func (i *float64Item) Clone() Item {
	ii := i.Item.Clone()
	return &float64Item{ii}
}

// An interface for a specific Vector type Item, i.e store multiple vectors
// within a single vector.
type VectorItem interface {
	Item
	ToVector() *Vector
	SetVector(*Vector)
}

type vectorItem struct {
	Item
}

// Creates a new VectorItem from the provided string key and *Vector value.
func NewVectorItem(key string, v *Vector) Item {
	i := KeyedItem(key)
	i.Provide(v)
	return &vectorItem{i}
}

//
func (i *vectorItem) ToVector() *Vector {
	ret := New("")
	err := json.Unmarshal(i.Value(), &ret)
	if err != nil {
		return nil
	}
	return ret
}

//
func (i *vectorItem) SetVector(v *Vector) {
	i.Provide(v)
}

//
func (i *vectorItem) Clone() Item {
	c := i.ToVector()
	ii := i.Item.Clone()
	ii.Provide(c)
	return &vectorItem{ii}
}
